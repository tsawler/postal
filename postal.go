package postal

type Service struct {
	ServerURL    string
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	ErrorChan    chan error
	MaxWorkers   int
}

var service Service

func New(s Service) *MailDispatcher {
	service = s
	return &MailDispatcher{
		WorkerPool: make(chan chan MailProcessingJob),
		maxWorkers: s.MaxWorkers,
		JobQueue:   make(chan MailProcessingJob),
	}
}
