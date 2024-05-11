package postal

import "html/template"

// MailMessage is a simple wrapper for MailData.
type MailMessage struct {
	MailData MailData
}

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
