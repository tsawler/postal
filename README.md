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
		ServerURL:    "http://localhost", // The URL of your server, in case you want backlinks in your message.
		SMTPServer:   "localhost",        // The address of the SMTP server you are sending through..
		SMTPPort:     1025,               // The port the SMTP server listens on.
		SMTPUser:     "",                 // SMTP username.
		SMTPPassword: "",                 // SMTP password.
		ErrorChan:    make(chan error),   // A channel to receive error messages (or nil for successful sends).
		MaxWorkers:   2,                  // The number of workers you want in the worker pool.
		MaxMessages:  10,                 // The size of the job queue (a buffered channel).
	}
	
	// Call postal.New with your service to get a message dispatcher.
	dispatcher, _ := postal.New(service)

	// Run the worker pool associated with the dispatcher.
	dispatcher.Run()

	// Create a mail message, of postal.MailData type.
	msg := postal.MailData{
		ToName:      "Me",
		ToAddress:   "me@here.com",
		FromName:    "Jack",
		FromAddress: "jack@there.com",
		Subject:     "Test subject",
		Content:     "Hello, world!",
		CC:          []string{"you@here.com", "him@here.com"},
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
		ServerURL:    "http://localhost", // The URL of your server, in case you want backlinks in your message.
		ErrorChan:     make(chan error),   // A channel to receive error messages (or nil for successful sends).
		MaxWorkers:    2,                  // The number of workers you want in the worker pool.
		MaxMessages:   10,                 // The size of the job queue (a buffered channel).
		SendingFromEU: false,              // Set to true if sending from European Union.
		Domain:        "yourdomain.com",   // The domain you are sending from.
		APIKey:        "your-mailgun-api", // The mailgun api key.
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