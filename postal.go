package postal

import "errors"

// Service is the type used to create a MailDispatcher.
type Service struct {
	ServerURL    string     // The URL of the server mail is sent from.
	SMTPServer   string     // The SMTP server.
	SMTPPort     int        // The SMTP server's port.
	SMTPUser     string     // The username for the SMTP server.
	SMTPPassword string     // The password for the SMTP server.
	ErrorChan    chan error // A channel to send errors (or nil) to.
	MaxWorkers   int        // Maximum number of workers in the pool.
	MaxMessages  int        // How big the buffer should be for the JobQueue.
}

var service Service

// New returns a new mail dispatcher, with a job queue channel and error channel. Send email
// by calling the Send method on the returned *MailDispatcher.
func New(s Service) (*MailDispatcher, error) {
	service = s

	// Sanity check.
	if service.MaxMessages == 0 {
		service.MaxMessages = 100
	}

	if service.MaxWorkers == 0 {
		service.MaxWorkers = 2
	}

	if service.SMTPServer == "" {
		return nil, errors.New("invalid smtp server")
	}

	if service.SMTPPort == 0 {
		return nil, errors.New("invalid smtp port")
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
