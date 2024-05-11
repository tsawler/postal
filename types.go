package postal

import "html/template"

// MailMessage is a simple wrapper for MailData.
type MailMessage struct {
	MailData MailData
}

// MailData holds all information for a given message.
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
	BCC          []string
	UseHermes    bool
	Attachments  []string
	Data         map[string]interface{}
	InlineImages []string
	ServerURL    string
}
