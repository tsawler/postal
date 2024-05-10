package postal

import "testing"

func Test_Run(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Error("Error getting dispatcher", err)
	}

	dispatcher.Run()

	dispatcher.Send(testMsg)
	err = <-service.ErrorChan
	if err != nil {
		t.Error("unexpected error when sending message", err)
	}
}

func Test_MailDispatcherSend(t *testing.T) {
	dispatcher, _ := New(testService)

	dispatcher.Run()
	dispatcher.Send(testMsg)
	err := <-service.ErrorChan
	if err != nil {
		t.Error("unexpected error when sending message", err)
	}

	oldTemplate := testMsg.Template
	testMsg.Template = "{{end}}"
	dispatcher.Send(testMsg)

	err = <-service.ErrorChan
	if err == nil {
		t.Error("no error with invalid template")
	}
	testMsg.Template = oldTemplate
}

func Test_sendViaMailGun(t *testing.T) {
	testService.Method = MailGun

	dispatcher, err := New(testService)
	if err != nil {
		t.Error(err)
	}

	t.Log("Got dispatcher")

	dispatcher.Run()

	t.Log("Ran dispatcher")
	dispatcher.Send(testMsg)

	t.Log("Sent msg")

	err = <-service.ErrorChan
	if err == nil {
		t.Error("expected error when sending message but did not get one")
	}
	testService.Method = SMTP
}
