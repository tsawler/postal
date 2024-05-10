package postal

import "html/template"

type MailMessage struct {
	MailData MailData
}

// MailData holds info for sending an email
type MailData struct {
	ToName       string
	ToAddress    string
	FromName     string
	FromAddress  string
	AdditionalTo []string
	Subject      string
	Content      template.HTML
	Template     string
	CC           []string
	UseHermes    bool
	Attachments  []string
	Data         map[string]interface{}
	InlineImages []string
	ServerURL    string
}
