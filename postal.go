package postal

type Service struct {
	ServerURL    string
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	ErrorChan    chan error
}

var (
	ServerURL    string
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	ErrorChan    chan error
)

func New(serverURL, smtp, user, password string, port int) *Service {
	ServerURL = serverURL
	SMTPServer = smtp
	SMTPPort = port
	SMTPPassword = password
	ErrorChan = make(chan error)

	return &Service{
		ServerURL:    serverURL,
		SMTPServer:   smtp,
		SMTPPort:     port,
		SMTPUser:     user,
		SMTPPassword: password,
		ErrorChan:    ErrorChan,
	}
}
