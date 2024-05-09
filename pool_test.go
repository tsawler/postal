package postal

import "testing"

func Test_Run(t *testing.T) {
	s := Service{
		ServerURL:    "http://localhost",
		SMTPServer:   "localhost",
		SMTPPort:     1026,
		SMTPUser:     "",
		SMTPPassword: "",
		ErrorChan:    make(chan error),
		MaxWorkers:   2,
		MaxMessages:  10,
	}
	dispatcher, _ := New(s)

	dispatcher.Run()

	msg := MailData{
		ToName:      "Me",
		ToAddress:   "me@here.com",
		FromName:    "Jack",
		FromAddress: "jack@there.com",
		Subject:     "Test subject",
		Content:     "Hello, world!",
		CC:          []string{"you@here.com", "him@here.com"},
	}

	dispatcher.Send(msg)
	err := <-service.ErrorChan
	if err != nil {
		t.Error("unexpected error when sending message")
	}
}
