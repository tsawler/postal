<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![Version](https://img.shields.io/badge/goversion-1.22.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tsawler/postal/master/LICENSE.md)
<a href="https://pkg.go.dev/github.com/tsawler/postal"><img src="https://img.shields.io/badge/godoc-reference-%23007d9c.svg"></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/tsawler/postal)](https://goreportcard.com/report/github.com/tsawler/postal)

# Postal

Postal is a simple package which allows you to send email messages, either via an SMTP server
or using [MailGun](https://www.mailgun.com/)'s api. Postal implements a simple worker pool to
make sending messages both easy and efficient.

## Installation

Install it the usual way:

~~~
go get -u github.com/tsawler/postal
~~~

## Usage

To create a worker pool to send email, first define a `postal.Service` object with settings appropriate for your environment.
This is the `postal.Service` type:

~~~go
// Service is the type used to create a MailDispatcher.
type Service struct {
    Method        int                           // How to send the message: postal.SMTP or postal.MailGun.
    ServerURL     string                        // The URL of the server mail is sent from.
    SMTPServer    string                        // The SMTP server.
    SMTPPort      int                           // The SMTP server's port.
    SMTPUser      string                        // The username for the SMTP server.
    SMTPPassword  string                        // The password for the SMTP server.
    ErrorChan     chan error                    // A channel to send errors (or nil) to.
    MaxWorkers    int                           // Maximum number of workers in the pool.
    MaxMessages   int                           // How big the buffer should be for the JobQueue.
    Domain        string                        // The domain used to send mail.
    APIKey        string                        // The API key for mailgun.
    SendingFromEU bool                          // If using mailgun and sending from EU, set to true.
    TemplateDir   string                        // Where templates are stored.
    templateMap   map[string]*template.Template // The map of preprocessed html templates.
}
~~~

Then, create a `postal.MailDispatcher` by calling `postal.New` with your `postal.Service` as the parameter:

~~~go
dispatcher, _ := postal.Service(myService)
~~~

Create a mail message by defining a `postal.MailData` object. This is the type for `postal.MailData`:

~~~go
// MailData holds all information for a given message.
type MailData struct {
	ToName       string         // The name of the recipient.
	ToAddress    string         // The email address of the recipient.
	FromName     string         // The name of the sender.
	FromAddress  string         // THe email address of the sender.
	AdditionalTo []string       // Additional TO recipients.
	Subject      string         // The subject of the email message.
	Content      template.HTML  // The content of the message, as HTML.
	Template     string         // The template to use. If not specified, will use a simple default template.
	CC           []string       // A slice of CC recipient emails.
	BCC          []string       // A slice of BCC recipient emails.
	Attachments  []string       // A slice of attachments, which must exist on disk (i.e. []string{"./files/myfile.pdf"}).
	Data         map[string]any // Data which is to be passed to the Go template.
	InlineImages []string       // A slice of images to be inlined in the email. PNG is preferred.
	ServerURL    string         // The URL of the server (for backlinks in message).
}
~~~

Finally, to send the message, call the `Send()` method on your dispatcher. Errors will be sent back on the `Service.ErrorChan`, and `nil` will
be sent back if the mail was queued successfully.

~~~go
dispatcher.Send(msg)

// Wait for a response.
err := <-service.ErrorChan
~~~

An example of sending via SMTP:

~~~go
package main

import (
	"fmt"
	"github.com/tsawler/postal"
)

func main() {
	// Define a postal.Service variable with values appropriate for your environment.
	service := postal.Service{
		Method:         postal.SMTP,       // Method to send. Options are postal.SMTP or postal.MailGun.
		ServerURL:     "http://localhost", // The URL of the server (for backlinks in mail).
		SMTPServer:    "localhost",        // The name of the SMTP server you are sending through.
		SMTPPort:      1025,               // The port the SMTP server listens on.
		SMTPUser:      "",                 // SMTP username.
		SMTPPassword:  "",                 // SMTP password.
		ErrorChan:     make(chan error),   // A channel to receive errors (or nil for success).
		MaxWorkers:    2,                  // Number of workers in our pool.
		MaxMessages:   10,                 // Buffer size for send channel.
		APIKey:        "",                 // API key is used for mailgun.
		SendingFromEU: false,              // If using mailgun, set to true if sending from EU.
		Domain:        "yourdomain.com",   // Only used for mailgun.
		TemplateDir:   "./templates",      // The directory where mail templates live.
	}
	
	// Get a dispatcher by calling postal.New(service).
	dispatcher, _ := postal.New(service)

	// Run the worker pool.
	dispatcher.Run()

	// Create a mail message.
	msg := postal.MailData{
		ToName:      "Me",
		ToAddress:   "me@here.com",
		FromName:    "Jack",
		FromAddress: "jack@there.com",
		Subject:     "Test subject",
		Content:     "Hello, world!",
		CC:          []string{"you@here.com", "him@here.com"},
		//Template:    "my-template.gohtml", // You can specify your own template, or leave this out and use the default.
	}

	// Send the message by calling dispatcher.Send.
	fmt.Println("Sending mail")
	dispatcher.Send(msg)

	// Wait for something back from ErrorChan.
	fmt.Println("Waiting for response")
	err := <-service.ErrorChan
	fmt.Println("Error", err)
}
~~~

An example sending through MailGun's api:

~~~go
package main

import (
	"fmt"
	"github.com/tsawler/postal"
)

func main() {
	
	// Define a postal.Service variable with values appropriate for your environment.
	service := postal.Service{
		Method:        postal.MailGun,     // Method to send. Options are postal.SMTP or postal.MailGun.
		ServerURL:    "http://localhost",  // The URL of your server, in case you want backlinks in your message.
		ErrorChan:     make(chan error),   // A channel to receive error messages (or nil for successful sends).
		MaxWorkers:    2,                  // The number of workers you want in the worker pool.
		MaxMessages:   10,                 // The size of the job queue (a buffered channel).
		SendingFromEU: false,              // Set to true if sending from European Union (mailgun only).
		TemplateDir:   "./templates",      // The directory where mail templates live.
		Domain:        "yourdomain.com",   // The domain you are sending from (mailgun only).
		APIKey:        "your-mailgun-api", // The mailgun api key (mailgun only).
	}
	
	// Call postal.New with your service to get a message dispatcher.
	dispatcher, _ := postal.New(service)

	// Run the worker pool associated with the dispatcher.
	dispatcher.Run()

	// Create a mail message, of postal.MailData type.
	msg := postal.MailData{
		ToName:      "Me",
		ToAddress:     "me@here.com",
		FromName:      "Jack",
		FromAddress:   "jack@there.com",
		Subject:       "Test subject",
		Content:       "Hello, world!",
		CC:            []string{"you@here.com", "him@here.com"},
	}

	// Send the message.
	fmt.Println("Sending mail")
	dispatcher.Send(msg)

	
	// Wait for a response.
	fmt.Println("Waiting for response")
	err := <-service.ErrorChan
	
	// If err is nil, then the message was sent.
	fmt.Println("Error", err)
}
~~~