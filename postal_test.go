package postal

import "testing"

func TestNew(t *testing.T) {
	s := Service{
		ServerURL:    "http://localhost",
		SMTPServer:   "localhost",
		SMTPPort:     1026,
		SMTPUser:     "",
		SMTPPassword: "",
		ErrorChan:    make(chan error),
		MaxWorkers:   0,
		MaxMessages:  0,
		TemplateDir:  "./testdata/templates",
	}

	_, err := New(s)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	s.SMTPServer = ""
	_, err = New(s)
	if err == nil {
		t.Error("error expected when SMTP server unspecified")
	}

	s.SMTPServer = "localhost"
	s.SMTPPort = 0
	_, err = New(s)
	if err == nil {
		t.Error("error expected when SMTP port unspecified")
	}

	s.SMTPPort = 1026
	s.ErrorChan = nil
	_, err = New(s)
	if err == nil {
		t.Error("error expected when error chan unspecified")
	}

	s.ErrorChan = make(chan error)
	s.Method = MailGun
	s.Domain = ""
	_, err = New(s)
	if err == nil {
		t.Error("error expected when method is mailgun and domain unspecified")
	}

	s.Method = MailGun
	s.Domain = "localhost"
	s.APIKey = ""
	_, err = New(s)
	if err == nil {
		t.Error("error expected when method is mailgun and api key unspecified")
	}
}
