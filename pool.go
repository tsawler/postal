package postal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aymerick/douceur/inliner"
	"github.com/mailgun/mailgun-go/v4"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"jaytaylor.com/html2text"
	"log"
	"time"
)

// MailProcessingJob is the unit of work to be performed. We wrap this type
// around a Video, which has all the information we need about the input source
// and what we want the output to look like.
type MailProcessingJob struct {
	MailMessage MailData
}

// newWorker takes a numeric id and a channel which accepts the chan MailProcessingJob
// type, and returns a worker object.
func newWorker(id int, workerPool chan chan MailProcessingJob) worker {
	return worker{
		id:         id,
		jobQueue:   make(chan MailProcessingJob),
		workerPool: workerPool,
	}
}

// worker holds info for a pool worker. It has the numeric id of the worker,
// the job queue, and the worker pool chan. A chan chan is used when the thing you want to
// send down a channel is another channel to send things back.
// See http://tleyden.github.io/blog/2013/11/23/understanding-chan-chans-in-go/
type worker struct {
	id         int
	jobQueue   chan MailProcessingJob      // Where we send jobs to process.
	workerPool chan chan MailProcessingJob // Our worker pool channel.
}

// start starts an individual worker.
func (w worker) start() {
	go func() {
		for {
			// Add jobQueue to the worker pool.
			w.workerPool <- w.jobQueue

			// Wait for a job to come back.
			job := <-w.jobQueue

			// Process the video with a worker.
			w.processJob(job)
		}
	}()
}

// MailDispatcher is the main interface to this package. Calling new returns an instance onf this type.
// This type has one method, Send, which is used to send an email.
type MailDispatcher struct {
	workerPool chan chan MailProcessingJob // Our worker pool channel.
	maxWorkers int                         // The maximum number of workers in our pool.
	JobQueue   chan MailProcessingJob      // The channel we send work to.
	ErrorChan  chan error                  // The channel we send errors (or nil) to.
}

// Send takes a msg in postal.MailData format, wraps it in a postal.MailProcessingJob
// and send it to the job queue for delivery.
func (md *MailDispatcher) Send(msg MailData) {
	job := MailProcessingJob{
		MailMessage: msg,
	}
	md.JobQueue <- job
}

// Run runs the workers in our worker pool.
func (md *MailDispatcher) Run() {
	for i := 0; i < md.maxWorkers; i++ {
		worker := newWorker(i+1, md.workerPool)
		worker.start()
	}

	go md.dispatch()
}

// dispatch waits for a MailProcessingJob job to come in over the job queue to send to a worker.
func (md *MailDispatcher) dispatch() {
	for {
		// Wait for a job to come in.
		job := <-md.JobQueue

		go func() {
			workerJobQueue := <-md.workerPool // assign a channel from our worker pool to workerJobPool.
			workerJobQueue <- job             // Send the unit of work to our queue.
		}()
	}
}

// processJob processes the main queue job by trying to deliver an email message. The resulting error, which will be
// nil if delivery is successful, is sent to the error chan.
func (w worker) processJob(m MailProcessingJob) {
	switch service.Method {
	case SMTP:
		w.sendViaSMTP(m)
	case MailGun:
		w.sendViaMailGun(m)
	}
}

// sendViaMailGun attempts to send an email using MailGun's api.
func (w worker) sendViaMailGun(m MailProcessingJob) {
	// Get the message body in both formats.
	plainTextMessage, formattedMessage, err := w.buildMessage(m)
	if err != nil {
		service.ErrorChan <- err
		return
	}

	// Create a mailgun client.
	mg := mailgun.NewMailgun(service.Domain, service.APIKey)
	if service.SendingFromEU {
		mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	}

	// Create the message in MailGun format.
	message := mg.NewMessage(m.MailMessage.FromAddress, m.MailMessage.Subject, plainTextMessage, m.MailMessage.ToAddress)
	message.SetHtml(formattedMessage)

	// Add additional to recipients.
	if len(m.MailMessage.AdditionalTo) > 0 {
		for _, x := range m.MailMessage.AdditionalTo {
			err := message.AddRecipient(x)
			if err != nil {
				service.ErrorChan <- err
				return
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

	// Set a 10-second context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10-second timeout.
	_, _, err = mg.Send(ctx, message)
	service.ErrorChan <- err
}

// sendViaSMTP attempts to send a message using an SMTP server.
func (w worker) sendViaSMTP(m MailProcessingJob) {
	// Get the message body in both formats.
	plainText, formattedMessage, err := w.buildMessage(m)
	if err != nil {
		service.ErrorChan <- err
		return
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
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		service.ErrorChan <- err
		return
	}

	// Create the mail message.
	email := mail.NewMSG()
	email.SetFrom(m.MailMessage.FromAddress).
		AddTo(m.MailMessage.ToAddress).
		SetSubject(m.MailMessage.Subject)

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

	email.SetBody(mail.TextPlain, plainText)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	}
	service.ErrorChan <- err
}

// buildMessage takes a mail processing job and sends back the message in two formats:
// plaintext and HTML.
func (w worker) buildMessage(m MailProcessingJob) (string, string, error) {
	var templateToParse string
	if m.MailMessage.Template == "" {
		templateToParse = fmt.Sprintf("%s/%s", service.TemplateDir, defaultTemplate)
		m.MailMessage.Template = defaultTemplate
	} else {
		templateToParse = m.MailMessage.Template
	}

	// check to see if the template exists in the cache
	var tmpl *template.Template

	// Lock the template map.
	mapLock.Lock()
	val, ok := service.templateMap[m.MailMessage.Template]
	if ok {
		// In cache, so use that.
		tmpl = val
	} else {
		// Not in cache, so create and add to cache.
		t, err := template.New(m.MailMessage.Template).ParseFiles(templateToParse)
		if err != nil {
			mapLock.Unlock()
			return "", "", err
		}
		tmpl = t
		service.templateMap[m.MailMessage.Template] = tmpl
	}
	// Unlock the map.
	mapLock.Unlock()

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
	plainText, err := html2text.FromString(result, html2text.Options{PrettyTables: true})
	if err != nil {
		plainText = ""
	}

	// Create html version of message.
	formattedMessage, err := inliner.Inline(result)
	if err != nil {
		service.ErrorChan <- err
		return "", "", err
	}

	return plainText, formattedMessage, nil
}
