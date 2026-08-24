package postal

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"sync"
	"time"

	"github.com/aymerick/douceur/inliner"
	"github.com/k3a/html2text"
	"github.com/mailgun/mailgun-go/v4"
	mail "github.com/xhit/go-simple-mail/v2"
)

// MailProcessingJob is the unit of work performed by a worker: one message to
// deliver, plus where to report the outcome.
type MailProcessingJob struct {
	MailMessage MailData

	// resultChan, when non-nil, receives this job's single result and nothing
	// else. It is created by SendAndWait with a buffer of one, so the worker's
	// hand-off always completes immediately even if the caller has already
	// given up waiting. This is what correlates a result with its message;
	// Service.ErrorChan, being shared by every in-flight message, cannot.
	resultChan chan error
}

// newWorker takes a numeric id, a channel which accepts the chan
// MailProcessingJob type, and the dispatcher the worker belongs to.
func newWorker(id int, workerPool chan chan MailProcessingJob, md *MailDispatcher) worker {
	return worker{
		id:         id,
		jobQueue:   make(chan MailProcessingJob),
		workerPool: workerPool,
		md:         md,
	}
}

// worker holds info for a pool worker. It has the numeric id of the worker,
// the job queue, the worker pool chan, and the dispatcher it serves. A chan chan
// is used when the thing you want to send down a channel is another channel to
// send things back.
// See http://tleyden.github.io/blog/2013/11/23/understanding-chan-chans-in-go/
type worker struct {
	id         int
	jobQueue   chan MailProcessingJob      // Where we send jobs to process.
	workerPool chan chan MailProcessingJob // Our worker pool channel.
	md         *MailDispatcher             // The dispatcher this worker serves.
}

// start starts an individual worker.
func (w worker) start() {
	go func() {
		defer w.md.wg.Done()
		for {
			// Add jobQueue to the worker pool.
			select {
			case w.workerPool <- w.jobQueue:
			case <-w.md.quit:
				return
			}

			// Wait for a job to come back.
			select {
			case job := <-w.jobQueue:
				w.processJob(job)
			case <-w.md.quit:
				return
			}
		}
	}()
}

// MailDispatcher is the main interface to this package. Calling New returns an
// instance of this type. Send mail with SendAndWait or Send.
//
// All configuration and template state is held per dispatcher, so several
// dispatchers may be created and used concurrently in one process.
type MailDispatcher struct {
	service    Service                     // This dispatcher's configuration.
	workerPool chan chan MailProcessingJob // Our worker pool channel.
	maxWorkers int                         // The maximum number of workers in our pool.
	JobQueue   chan MailProcessingJob      // The channel we send work to.
	ErrorChan  chan error                  // Optional channel results are sent to.

	templateMap map[string]*template.Template // Cache of preprocessed html templates.
	mapLock     sync.Mutex                    // Guards templateMap.

	quit      chan struct{}  // Closed by Close to stop workers and the dispatch loop.
	closeOnce sync.Once      // Makes Close idempotent.
	wg        sync.WaitGroup // Tracks workers and the dispatch loop.
	running   bool           // Whether Run has been called.
	runLock   sync.Mutex     // Guards running.
}

// Send takes a msg in postal.MailData format, wraps it in a
// postal.MailProcessingJob and sends it to the job queue for delivery.
//
// Send is fire and forget: it returns as soon as the message is queued, and the
// result is delivered to Service.ErrorChan if one was configured. Because that
// channel is shared by every in-flight message, a result read from it cannot be
// attributed to any particular Send. Prefer SendAndWait when the caller needs
// to know whether this message was delivered.
func (md *MailDispatcher) Send(msg MailData) {
	md.JobQueue <- MailProcessingJob{MailMessage: msg}
}

// SendAndWait delivers one message and returns its result.
//
// The result belongs to this message alone: the job carries a private buffered
// channel, so a worker never blocks handing the result back and no other caller
// can consume it. If ctx expires first, the worker still completes and reports
// without blocking, and ctx.Err() is returned here.
//
// The context should allow at least Service.MaxSendDuration, or slow but
// successful sends will be abandoned and misreported as failures.
func (md *MailDispatcher) SendAndWait(ctx context.Context, msg MailData) error {
	md.runLock.Lock()
	running := md.running
	md.runLock.Unlock()
	if !running {
		return ErrNotRunning
	}

	// Buffered so the worker's hand-off below never blocks, whatever this
	// caller does. This is the property the old shared unbuffered channel
	// lacked, and the reason workers could be parked forever.
	result := make(chan error, 1)
	job := MailProcessingJob{MailMessage: msg, resultChan: result}

	select {
	case md.JobQueue <- job:
	case <-ctx.Done():
		return ctx.Err()
	case <-md.quit:
		return ErrDispatcherClosed
	}

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-md.quit:
		return ErrDispatcherClosed
	}
}

// Run runs the workers in our worker pool. It is safe to call more than once;
// subsequent calls are no-ops.
func (md *MailDispatcher) Run() {
	md.runLock.Lock()
	defer md.runLock.Unlock()
	if md.running {
		return
	}
	md.running = true

	for i := 0; i < md.maxWorkers; i++ {
		md.wg.Add(1)
		worker := newWorker(i+1, md.workerPool, md)
		worker.start()
	}

	md.wg.Add(1)
	go md.dispatch()
}

// Close stops the dispatch loop and the worker pool, and waits for in-flight
// deliveries to finish. It is safe to call more than once. A dispatcher cannot
// be restarted after Close.
func (md *MailDispatcher) Close() {
	md.closeOnce.Do(func() {
		close(md.quit)
	})
	md.wg.Wait()
}

// dispatch waits for a MailProcessingJob to come in over the job queue and
// hands it to a free worker.
//
// Unlike the previous implementation this does not spawn a goroutine per job;
// it waits for a worker to become free, so a burst of mail applies
// backpressure through JobQueue instead of accumulating goroutines.
func (md *MailDispatcher) dispatch() {
	defer md.wg.Done()

	for {
		// Wait for a job to come in.
		var job MailProcessingJob
		select {
		case job = <-md.JobQueue:
		case <-md.quit:
			return
		}

		// Wait for a free worker, then hand off the unit of work.
		select {
		case workerJobQueue := <-md.workerPool:
			select {
			case workerJobQueue <- job:
			case <-md.quit:
				return
			}
		case <-md.quit:
			return
		}
	}
}

// report delivers a job's single result.
//
// Exactly one report happens per job. Delivery never blocks a worker
// indefinitely: a private result channel is buffered, and the shared ErrorChan
// is given only a bounded grace period before the result is dropped.
func (md *MailDispatcher) report(job MailProcessingJob, err error) {
	if job.resultChan != nil {
		// Buffered with capacity one and written exactly once, so this cannot
		// block regardless of whether the caller is still waiting.
		job.resultChan <- err
		return
	}

	if md.ErrorChan == nil {
		return
	}

	// A reader that is already waiting takes this immediately. Otherwise wait
	// only briefly: losing a result is far better than losing a worker.
	timer := time.NewTimer(md.service.errorChanGrace())
	defer timer.Stop()

	select {
	case md.ErrorChan <- err:
	case <-timer.C:
		if err != nil {
			log.Printf("postal: no reader on ErrorChan; dropping send result: %v", err)
		}
	case <-md.quit:
	}
}

// processJob delivers one message and reports the outcome exactly once.
func (w worker) processJob(m MailProcessingJob) {
	var err error

	switch w.md.service.Method {
	case SMTP:
		err = w.sendViaSMTP(m)
	case MailGun:
		err = w.sendViaMailGun(m)
	default:
		err = fmt.Errorf("postal: unknown send method %d", w.md.service.Method)
	}

	w.md.report(m, err)
}

// sendViaMailGun attempts to send an email using MailGun's api, returning the
// outcome. It never writes to a result channel itself; processJob reports once.
func (w worker) sendViaMailGun(m MailProcessingJob) error {
	service := w.md.service

	// Get the message body in both formats.
	plainTextMessage, formattedMessage, err := w.md.buildMessage(m)
	if err != nil {
		return err
	}

	// Create a mailgun client.
	mg := mailgun.NewMailgun(service.Domain, service.APIKey)
	if service.SendingFromEU {
		mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	}

	// Create the message in MailGun format.
	fromAddr := m.MailMessage.FromAddress
	if m.MailMessage.FromName != "" {
		fromAddr = fmt.Sprintf("%s <%s>", m.MailMessage.FromName, m.MailMessage.FromAddress)
	}
	message := mg.NewMessage(fromAddr, m.MailMessage.Subject, plainTextMessage, m.MailMessage.ToAddress)

	// Set HTML body only for HTML content type (default).
	contentType := m.MailMessage.ContentType
	if contentType == "" || contentType == "text/html" {
		message.SetHtml(formattedMessage)
	}

	// Set reply-to address.
	if m.MailMessage.ReplyTo != "" {
		message.SetReplyTo(m.MailMessage.ReplyTo)
	}

	// Add additional to recipients.
	if len(m.MailMessage.AdditionalTo) > 0 {
		for _, x := range m.MailMessage.AdditionalTo {
			if err := message.AddRecipient(x); err != nil {
				return err
			}
		}
	}

	// Add cc recipients.
	if len(m.MailMessage.CC) > 0 {
		for _, x := range m.MailMessage.CC {
			message.AddCC(x)
		}
	}

	// Add bcc recipients.
	if len(m.MailMessage.BCC) > 0 {
		for _, x := range m.MailMessage.BCC {
			message.AddBCC(x)
		}
	}

	// Add attachments.
	if len(m.MailMessage.Attachments) > 0 {
		for _, x := range m.MailMessage.Attachments {
			message.AddAttachment(x)
		}
	}

	// To add image to template, use this syntax:
	//     <img alt="alt text" src="cid:filename.png">
	if len(m.MailMessage.InlineImages) > 0 {
		for _, x := range m.MailMessage.InlineImages {
			message.AddInline(x)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), service.sendTimeout())
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	return err
}

// sendViaSMTP attempts to send a message using an SMTP server, returning the
// outcome. It never writes to a result channel itself; processJob reports once.
func (w worker) sendViaSMTP(m MailProcessingJob) error {
	service := w.md.service

	// Get the message body in both formats.
	plainText, formattedMessage, err := w.md.buildMessage(m)
	if err != nil {
		return err
	}

	// Create smtp client.
	server := mail.NewSMTPClient()
	server.Host = service.SMTPServer
	server.Port = service.SMTPPort
	server.Username = service.SMTPUser
	server.Password = service.SMTPPassword
	if service.SMTPServer == "localhost" {
		server.Authentication = mail.AuthPlain
	} else {
		server.Authentication = mail.AuthLogin
	}
	server.Encryption = mail.EncryptionTLS

	server.KeepAlive = false
	server.ConnectTimeout = service.connectTimeout()
	server.SendTimeout = service.sendTimeout()

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Create the mail message.
	fromAddr := m.MailMessage.FromAddress
	if m.MailMessage.FromName != "" {
		fromAddr = fmt.Sprintf("%s <%s>", m.MailMessage.FromName, m.MailMessage.FromAddress)
	}
	email := mail.NewMSG()
	email.SetFrom(fromAddr).
		AddTo(m.MailMessage.ToAddress).
		SetSubject(m.MailMessage.Subject)

	// Set reply-to address.
	if m.MailMessage.ReplyTo != "" {
		email.SetReplyTo(m.MailMessage.ReplyTo)
	}

	// Add additional to recipients.
	if len(m.MailMessage.AdditionalTo) > 0 {
		for _, x := range m.MailMessage.AdditionalTo {
			email.AddTo(x)
		}
	}

	// Add cc recipients.
	if len(m.MailMessage.CC) > 0 {
		for _, x := range m.MailMessage.CC {
			email.AddCc(x)
		}
	}

	// Add bcc recipients.
	if len(m.MailMessage.BCC) > 0 {
		for _, x := range m.MailMessage.BCC {
			email.AddBcc(x)
		}
	}

	// Add attachments.
	if len(m.MailMessage.Attachments) > 0 {
		for _, x := range m.MailMessage.Attachments {
			email.AddAttachment(x)
		}
	}

	// To add image to template, use this syntax:
	//     <img alt="alt text" src="cid:filename.png">
	if len(m.MailMessage.InlineImages) > 0 {
		for _, x := range m.MailMessage.InlineImages {
			email.AddInline(x)
		}
	}

	// Set body based on content type.
	contentType := m.MailMessage.ContentType
	if contentType == "" {
		contentType = "text/html"
	}

	switch contentType {
	case "text/plain":
		email.SetBody(mail.TextPlain, plainText)
	case "text/html":
		email.SetBody(mail.TextPlain, plainText)
		email.AddAlternative(mail.TextHTML, formattedMessage)
	default:
		// For other content types (e.g., application/xml), send as plain text body.
		email.SetBody(mail.TextPlain, formattedMessage)
	}

	if err := email.Send(smtpClient); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// buildMessage takes a mail processing job and sends back the message in two
// formats: plaintext and HTML.
//
// It only returns errors. An earlier version also wrote failures to the shared
// error channel before returning them, so a single message could produce two
// results and permanently desynchronize every later reader.
func (md *MailDispatcher) buildMessage(m MailProcessingJob) (string, string, error) {
	var templateToParse string
	if m.MailMessage.Template == "" {
		templateToParse = fmt.Sprintf("%s/%s", md.service.TemplateDir, defaultTemplate)
		m.MailMessage.Template = defaultTemplate
	} else {
		templateToParse = fmt.Sprintf("%s/%s", md.service.TemplateDir, m.MailMessage.Template)
	}

	// check to see if the template exists in the cache
	var tmpl *template.Template

	// Lock the template map.
	md.mapLock.Lock()
	val, ok := md.templateMap[m.MailMessage.Template]
	if ok {
		// In cache, so use that.
		tmpl = val
	} else {
		// Not in cache, so create and add to cache.
		t, err := template.New(m.MailMessage.Template).ParseFiles(templateToParse)
		if err != nil {
			md.mapLock.Unlock()
			return "", "", err
		}
		tmpl = t
		md.templateMap[m.MailMessage.Template] = tmpl
	}
	// Unlock the map.
	md.mapLock.Unlock()

	data := struct {
		Content   template.HTML
		From      string
		FromName  string
		ServerUrl string
		Data      map[string]any
	}{
		Content:   m.MailMessage.Content,
		FromName:  m.MailMessage.FromName,
		From:      m.MailMessage.FromAddress,
		ServerUrl: m.MailMessage.ServerURL,
		Data:      m.MailMessage.Data,
	}

	// Execute the template with data.
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", "", err
	}

	// Get the rendered template as a string.
	result := tpl.String()

	// Create plaintext version of message.
	plainText := html2text.HTML2Text(result)

	// Create html version of message.
	formattedMessage, err := inliner.Inline(result)
	if err != nil {
		return "", "", err
	}

	return plainText, formattedMessage, nil
}
