package postal

import (
	"errors"
	"fmt"
	"github.com/tsawler/toolbox"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	SMTP            = 1
	MailGun         = 2
	defaultTemplate = "default.gohtml"
)

var mapLock sync.Mutex

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

func checkTemplateDir() error {
	// Get default templates from github.
	var t toolbox.Tools
	_ = t.CreateDirIfNotExist(service.TemplateDir)

	_, err := os.Stat(fmt.Sprintf("%s/default.gohtml", service.TemplateDir))
	if os.IsNotExist(err) {
		fmt.Println("Getting default template from remote source...")
		resp, err := http.Get("https://raw.githubusercontent.com/tsawler/postal-templates/main/action.html")
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := os.WriteFile(fmt.Sprintf("%s/default.gohtml", service.TemplateDir), body, 0666); err != nil {
			return err
		}
		fmt.Println("Done!")

	}
	return nil
}

var service Service

// New returns a new mail dispatcher, with a job queue channel and error channel. Send email
// by calling the Send method on the returned *MailDispatcher.
func New(s Service) (*MailDispatcher, error) {
	service = s

	// Sanity check.
	if service.Method == 0 {
		service.Method = SMTP
	}

	if service.MaxMessages == 0 {
		service.MaxMessages = 100
	}

	if service.MaxWorkers == 0 {
		service.MaxWorkers = 2
	}

	if service.SMTPServer == "" && service.Method == SMTP {
		return nil, errors.New("invalid smtp server")
	}

	if service.SMTPPort == 0 && service.Method == SMTP {
		return nil, errors.New("invalid smtp port")
	}

	if service.APIKey == "" && service.Method == MailGun {
		return nil, errors.New("api key required")
	}

	if service.Method == MailGun && service.Domain == "" {
		return nil, errors.New("domain required when using mailgun")
	}

	if service.ErrorChan == nil {
		return nil, errors.New("ErrorChan must be specified")
	}

	if service.TemplateDir == "" {
		service.TemplateDir = "./templates/mail"
	} else if strings.HasSuffix(service.TemplateDir, "/") {
		service.TemplateDir = strings.TrimSuffix(service.TemplateDir, "/")
	}

	// Get the default template if it does not exist.
	err := checkTemplateDir()
	if err != nil {
		return nil, err
	}

	// Set up the template map for caching templates.
	service.templateMap = make(map[string]*template.Template)

	return &MailDispatcher{
		workerPool: make(chan chan MailProcessingJob),
		maxWorkers: s.MaxWorkers,
		JobQueue:   make(chan MailProcessingJob, s.MaxMessages),
		ErrorChan:  service.ErrorChan,
	}, nil
}
