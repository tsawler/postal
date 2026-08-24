package postal

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tsawler/toolbox"
)

const (
	SMTP            = 1
	MailGun         = 2
	defaultTemplate = "default.gohtml"
)

// Default timeouts. Each may be overridden per Service.
const (
	// DefaultConnectTimeout is how long to wait for an SMTP connection.
	DefaultConnectTimeout = 10 * time.Second

	// DefaultSendTimeout is how long to wait for a message to be transmitted
	// once connected. It also bounds the Mailgun API call.
	DefaultSendTimeout = 10 * time.Second

	// DefaultErrorChanGrace is how long a worker will wait to hand a result to
	// a Service.ErrorChan that nobody is currently reading, before abandoning
	// that result and returning to the pool.
	//
	// A worker is never blocked for longer than this. Previously the hand-off
	// was an unconditional channel send, so a result arriving after its caller
	// had stopped waiting parked that worker permanently; once every worker was
	// parked the dispatcher accepted mail forever without sending any of it.
	// Use SendAndWait to avoid the situation entirely.
	DefaultErrorChanGrace = 5 * time.Second
)

// Errors returned by the dispatcher.
var (
	// ErrDispatcherClosed is returned by SendAndWait once Close has been called.
	ErrDispatcherClosed = errors.New("postal: dispatcher is closed")

	// ErrNotRunning is returned by SendAndWait when Run has not been called.
	// Without workers, a queued message would never be consumed.
	ErrNotRunning = errors.New("postal: dispatcher is not running; call Run first")
)

// Service is the type used to create a MailDispatcher.
//
// A Service is copied into the MailDispatcher that New returns, and each
// dispatcher keeps its own configuration. Multiple dispatchers may be created
// and used concurrently within one process.
type Service struct {
	Method        int        // How to send the message: postal.SMTP or postal.MailGun.
	ServerURL     string     // The URL of the server mail is sent from.
	SMTPServer    string     // The SMTP server.
	SMTPPort      int        // The SMTP server's port.
	SMTPUser      string     // The username for the SMTP server.
	SMTPPassword  string     // The password for the SMTP server.
	ErrorChan     chan error // Optional channel for results (an error, or nil on success).
	MaxWorkers    int        // Maximum number of workers in the pool.
	MaxMessages   int        // How big the buffer should be for the JobQueue.
	Domain        string     // The domain used to send mail.
	APIKey        string     // The API key for mailgun.
	SendingFromEU bool       // If using mailgun and sending from EU, set to true.
	TemplateDir   string     // Where templates are stored.

	// ConnectTimeout bounds establishing the SMTP connection.
	// Zero means DefaultConnectTimeout.
	ConnectTimeout time.Duration

	// SendTimeout bounds transmitting the message once connected, and bounds
	// the Mailgun API call. Zero means DefaultSendTimeout.
	SendTimeout time.Duration

	// ErrorChanGrace bounds how long a worker waits to deliver a result to
	// ErrorChan when no reader is ready. Zero means DefaultErrorChanGrace.
	//
	// Callers needing results delivered reliably should use SendAndWait, which
	// carries a private result channel and never drops anything.
	ErrorChanGrace time.Duration
}

func (s Service) connectTimeout() time.Duration {
	if s.ConnectTimeout <= 0 {
		return DefaultConnectTimeout
	}
	return s.ConnectTimeout
}

func (s Service) sendTimeout() time.Duration {
	if s.SendTimeout <= 0 {
		return DefaultSendTimeout
	}
	return s.SendTimeout
}

func (s Service) errorChanGrace() time.Duration {
	if s.ErrorChanGrace <= 0 {
		return DefaultErrorChanGrace
	}
	return s.ErrorChanGrace
}

// MaxSendDuration reports the longest a single delivery attempt can take,
// excluding template rendering. Callers imposing their own deadline on
// SendAndWait should keep it comfortably above this value, or slow but
// successful sends will be abandoned and misreported as failures.
func (s Service) MaxSendDuration() time.Duration {
	return s.connectTimeout() + s.sendTimeout()
}

// checkTemplateDir ensures the template directory exists and holds a default
// template, fetching one if it is absent.
func checkTemplateDir(templateDir string) error {
	// Get default templates from github.
	var t toolbox.Tools
	_ = t.CreateDirIfNotExist(templateDir)

	_, err := os.Stat(fmt.Sprintf("%s/default.gohtml", templateDir))
	if os.IsNotExist(err) {
		fmt.Println("Getting default template from remote source...")
		resp, err := http.Get("https://raw.githubusercontent.com/tsawler/postal-templates/main/action.html")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := os.WriteFile(fmt.Sprintf("%s/default.gohtml", templateDir), body, 0666); err != nil {
			return err
		}
		fmt.Println("Done!")
	}

	return nil
}

// New returns a new mail dispatcher. Call Run to start its worker pool, then
// send mail with SendAndWait (synchronous, result returned to the caller) or
// Send (fire and forget, result delivered to Service.ErrorChan if set). Call
// Close when finished.
func New(s Service) (*MailDispatcher, error) {
	// Sanity check.
	if s.Method == 0 {
		s.Method = SMTP
	}

	if s.MaxMessages == 0 {
		s.MaxMessages = 100
	}

	if s.MaxWorkers == 0 {
		s.MaxWorkers = 2
	}

	if s.SMTPServer == "" && s.Method == SMTP {
		return nil, errors.New("invalid smtp server")
	}

	if s.SMTPPort == 0 && s.Method == SMTP {
		return nil, errors.New("invalid smtp port")
	}

	if s.APIKey == "" && s.Method == MailGun {
		return nil, errors.New("api key required")
	}

	if s.Method == MailGun && s.Domain == "" {
		return nil, errors.New("domain required when using mailgun")
	}

	// ErrorChan is deliberately optional: SendAndWait reports results directly
	// to its caller and needs no shared channel.

	if s.TemplateDir == "" {
		s.TemplateDir = "./templates/mail"
	} else if strings.HasSuffix(s.TemplateDir, "/") {
		s.TemplateDir = strings.TrimSuffix(s.TemplateDir, "/")
	}

	// Get the default template if it does not exist.
	if err := checkTemplateDir(s.TemplateDir); err != nil {
		return nil, err
	}

	return &MailDispatcher{
		service:     s,
		workerPool:  make(chan chan MailProcessingJob),
		maxWorkers:  s.MaxWorkers,
		JobQueue:    make(chan MailProcessingJob, s.MaxMessages),
		ErrorChan:   s.ErrorChan,
		templateMap: make(map[string]*template.Template),
		quit:        make(chan struct{}),
	}, nil
}
