package postal

import (
	"bytes"
	"context"
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

// MailDispatcher holds info for a dispatcher.
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

func (w worker) sendViaMailGun(m MailProcessingJob) {

	_, formattedMessage, err := w.buildMessage(m)
	if err != nil {
		service.ErrorChan <- err
		return
	}

	mg := mailgun.NewMailgun(service.Domain, service.APIKey)
	if service.SendingFromEU {
		mg.SetAPIBase("https://api.eu.mailgun.net/v3")
	}

	message := mg.NewMessage(m.MailMessage.FromAddress, m.MailMessage.Subject, formattedMessage, m.MailMessage.ToAddress)

	if len(m.MailMessage.AdditionalTo) > 0 {
		for _, x := range m.MailMessage.AdditionalTo {
			err := message.AddRecipient(x)
			if err != nil {
				service.ErrorChan <- err
				return
			}
		}
	}

	if len(m.MailMessage.CC) > 0 {
		for _, x := range m.MailMessage.CC {
			message.AddCC(x)
		}
	}

	if len(m.MailMessage.Attachments) > 0 {
		for _, x := range m.MailMessage.Attachments {
			message.AddAttachment(x)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err = mg.Send(ctx, message)
	service.ErrorChan <- err
}

func (w worker) sendViaSMTP(m MailProcessingJob) {
	plainText, formattedMessage, err := w.buildMessage(m)
	if err != nil {
		service.ErrorChan <- err
		return
	}

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
		log.Println(err)
		service.ErrorChan <- err
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.MailMessage.FromAddress).
		AddTo(m.MailMessage.ToAddress).
		SetSubject(m.MailMessage.Subject)

	if len(m.MailMessage.AdditionalTo) > 0 {
		for _, x := range m.MailMessage.AdditionalTo {
			email.AddTo(x)
		}
	}

	if len(m.MailMessage.CC) > 0 {
		for _, x := range m.MailMessage.CC {
			email.AddCc(x)
		}
	}

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

func (w worker) buildMessage(m MailProcessingJob) (string, string, error) {
	var templateToParse string
	if m.MailMessage.Template == "" {
		templateToParse = bootstrapTemplate
	} else {
		templateToParse = m.MailMessage.Template
	}

	t, err := template.New("msg").Parse(templateToParse)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	data := struct {
		Content       template.HTML
		From          string
		FromName      string
		PreferenceMap map[string]string
		ServerUrl     string
		IntMap        map[string]int
		StringMap     map[string]string
		FloatMap      map[string]float32
		RowSets       map[string]interface{}
	}{
		Content:   m.MailMessage.Content,
		FromName:  m.MailMessage.FromName,
		From:      m.MailMessage.FromAddress,
		ServerUrl: m.MailMessage.ServerURL,
		IntMap:    m.MailMessage.IntMap,
		StringMap: m.MailMessage.StringMap,
		FloatMap:  m.MailMessage.FloatMap,
		RowSets:   m.MailMessage.RowSets,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
		return "", "", err
	}

	result := tpl.String()

	plainText, err := html2text.FromString(result, html2text.Options{PrettyTables: true})
	if err != nil {
		plainText = ""
	}

	var formattedMessage string

	formattedMessage, err = inliner.Inline(result)
	if err != nil {
		log.Println(err)
		service.ErrorChan <- err
		return "", "", err
	}
	return plainText, formattedMessage, nil
}
