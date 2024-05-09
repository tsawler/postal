package postal

import "errors"

const (
	SMTP    = 1
	MailGun = 2
)

// Service is the type used to create a MailDispatcher.
type Service struct {
	Method        int        // How to send the message: postal.SMTP or postal.MailGun.
	ServerURL     string     // The URL of the server mail is sent from.
	SMTPServer    string     // The SMTP server.
	SMTPPort      int        // The SMTP server's port.
	SMTPUser      string     // The username for the SMTP server.
	SMTPPassword  string     // The password for the SMTP server.
	ErrorChan     chan error // A channel to send errors (or nil) to.
	MaxWorkers    int        // Maximum number of workers in the pool.
	MaxMessages   int        // How big the buffer should be for the JobQueue.
	Domain        string     // The domain used to send mail.
	APIKey        string     // The API key for mailgun.
	SendingFromEU bool       // If using mailgun and sending from EU, set to true.
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

	if service.Method == MailGun && service.Domain == "" {
		return nil, errors.New("domain required when using mailgun")
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

	if service.Domain == "" && service.Method == MailGun {
		return nil, errors.New("domain required")
	}

	if service.ErrorChan == nil {
		return nil, errors.New("ErrorChan must be specified")
	}

	return &MailDispatcher{
		workerPool: make(chan chan MailProcessingJob),
		maxWorkers: s.MaxWorkers,
		JobQueue:   make(chan MailProcessingJob, s.MaxMessages),
		ErrorChan:  service.ErrorChan,
	}, nil
}
