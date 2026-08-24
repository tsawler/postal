package postal

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func Test_Run(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Error("Error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	dispatcher.Send(testMsg)
	err = <-dispatcher.ErrorChan
	if err != nil {
		t.Error("unexpected error when sending message", err)
	}
}

func Test_MailDispatcherSend(t *testing.T) {
	dispatcher, _ := New(testService)

	dispatcher.Run()
	defer dispatcher.Close()

	dispatcher.Send(testMsg)
	err := <-dispatcher.ErrorChan
	if err != nil {
		t.Error("unexpected error when sending message", err)
	}

	oldTemplate := testMsg.Template
	testMsg.Template = "{{end}}"
	dispatcher.Send(testMsg)

	err = <-dispatcher.ErrorChan
	if err == nil {
		t.Error("no error with invalid template")
	}
	testMsg.Template = oldTemplate
}

func Test_sendViaMailGun(t *testing.T) {
	// A local copy, so this test cannot leak configuration into the others.
	s := testService
	s.Method = MailGun
	s.SendingFromEU = true

	dispatcher, err := New(s)
	if err != nil {
		t.Error(err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	dispatcher.Send(testMsg)

	err = <-dispatcher.ErrorChan
	if err == nil {
		t.Error("expected error when sending message but did not get one")
	}
}

// =====================================================
// SendAndWait
// =====================================================

func Test_SendAndWait(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, testMsg); err != nil {
		t.Error("unexpected error when sending message", err)
	}
}

func Test_SendAndWaitReportsFailure(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	bad := testMsg
	bad.Template = "{{end}}"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, bad); err == nil {
		t.Error("expected an error for an invalid template")
	}
}

func Test_SendAndWaitRequiresRun(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// No Run, so no worker would ever consume the message. This must report
	// rather than block until the context expires.
	if err := dispatcher.SendAndWait(ctx, testMsg); !errors.Is(err, ErrNotRunning) {
		t.Errorf("expected ErrNotRunning, got %v", err)
	}
}

func Test_SendAndWaitAfterClose(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	dispatcher.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, testMsg); !errors.Is(err, ErrDispatcherClosed) {
		t.Errorf("expected ErrDispatcherClosed, got %v", err)
	}
}

// Test_SendAndWaitResultsAreCorrelated is the regression for cross-attribution:
// concurrent callers must each receive the result for their own message, not
// whichever result happens to arrive first. A single shared channel cannot do
// that, which is why SendAndWait carries a private one.
func Test_SendAndWaitResultsAreCorrelated(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	const senders = 8

	var wg sync.WaitGroup
	results := make([]error, senders)

	for i := 0; i < senders; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			msg := testMsg
			// Odd senders use a broken template, so exactly half must fail.
			if i%2 == 1 {
				msg.Template = "{{end}}"
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			results[i] = dispatcher.SendAndWait(ctx, msg)
		}(i)
	}
	wg.Wait()

	for i, err := range results {
		wantErr := i%2 == 1
		if wantErr && err == nil {
			t.Errorf("sender %d: expected failure, got success (result was misattributed)", i)
		}
		if !wantErr && err != nil {
			t.Errorf("sender %d: expected success, got %v (result was misattributed)", i, err)
		}
	}
}

// =====================================================
// Liveness
// =====================================================

// Test_WorkersSurviveAbandonedResults is the core regression. Results that
// nobody waits for used to park a worker on an unconditional channel send, and
// once every worker was parked the dispatcher silently stopped sending mail.
// Here every result is abandoned, and the pool must still work afterwards.
func Test_WorkersSurviveAbandonedResults(t *testing.T) {
	s := testService
	// Unbuffered and never read from: the worst case for the old code.
	s.ErrorChan = make(chan error)
	// Keep the test quick; a dropped result should not cost the full default.
	s.ErrorChanGrace = 250 * time.Millisecond

	dispatcher, err := New(s)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	// Send more messages than there are workers, reading none of the results.
	for i := 0; i < s.MaxWorkers*3; i++ {
		dispatcher.Send(testMsg)
	}

	// Give the workers time to process and drop those results.
	time.Sleep(2 * time.Second)

	// The pool must still be alive. Under the old behavior every worker would
	// be parked by now and this would block until the context expired.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, testMsg); err != nil {
		t.Errorf("pool did not survive abandoned results: %v", err)
	}
}

// Test_ExpiredCallerDoesNotParkWorker covers the same hazard from the
// SendAndWait side: a caller whose context expires must not strand the worker
// still delivering its message.
func Test_ExpiredCallerDoesNotParkWorker(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	defer dispatcher.Close()

	// A context so short it is certain to expire before delivery completes.
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, testMsg); err == nil {
		t.Log("send completed before the context expired; inconclusive but harmless")
	}

	// Every worker must be back in the pool and usable.
	for i := 0; i < testService.MaxWorkers; i++ {
		okCtx, okCancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := dispatcher.SendAndWait(okCtx, testMsg)
		okCancel()
		if err != nil {
			t.Fatalf("worker was stranded by an abandoned send: %v", err)
		}
	}
}

// Test_CloseIsIdempotent guards shutdown against double-close panics.
func Test_CloseIsIdempotent(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	dispatcher.Close()
	dispatcher.Close()
}

// Test_RunIsIdempotent ensures a second Run does not start a second pool.
func Test_RunIsIdempotent(t *testing.T) {
	dispatcher, err := New(testService)
	if err != nil {
		t.Fatal("error getting dispatcher", err)
	}

	dispatcher.Run()
	dispatcher.Run()
	defer dispatcher.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := dispatcher.SendAndWait(ctx, testMsg); err != nil {
		t.Errorf("unexpected error after a repeated Run: %v", err)
	}
}

// =====================================================
// Independent dispatchers
// =====================================================

// Test_MultipleDispatchersAreIndependent is the regression for the package
// global that New used to overwrite. A second dispatcher used to redirect the
// first one's results to its own channel and replace its configuration, so the
// first one's sends failed or timed out unconditionally.
func Test_MultipleDispatchersAreIndependent(t *testing.T) {
	first, err := New(testService)
	if err != nil {
		t.Fatal("error creating first dispatcher", err)
	}
	first.Run()
	defer first.Close()

	// A second dispatcher with a different method and its own channel.
	secondService := testService
	secondService.Method = MailGun
	secondService.ErrorChan = make(chan error, 1)

	second, err := New(secondService)
	if err != nil {
		t.Fatal("error creating second dispatcher", err)
	}
	second.Run()
	defer second.Close()

	// The first dispatcher must still be the SMTP one, and still succeed.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := first.SendAndWait(ctx, testMsg); err != nil {
		t.Errorf("first dispatcher broke after a second was created: %v", err)
	}

	// And the second must still be the Mailgun one, which fails with the
	// invalid test credentials.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	if err := second.SendAndWait(ctx2, testMsg); err == nil {
		t.Error("second dispatcher did not use its own Mailgun configuration")
	}
}

// Test_TemplateCacheIsPerDispatcher confirms template state moved off the
// package global along with the rest of the configuration.
func Test_TemplateCacheIsPerDispatcher(t *testing.T) {
	a, err := New(testService)
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(testService)
	if err != nil {
		t.Fatal(err)
	}

	if a.templateMap == nil || b.templateMap == nil {
		t.Fatal("expected each dispatcher to have its own template cache")
	}

	a.mapLock.Lock()
	a.templateMap["sentinel"] = nil
	a.mapLock.Unlock()

	b.mapLock.Lock()
	_, leaked := b.templateMap["sentinel"]
	b.mapLock.Unlock()

	if leaked {
		t.Error("template cache is shared between dispatchers")
	}
}

// =====================================================
// Configuration
// =====================================================

func Test_TimeoutDefaults(t *testing.T) {
	var s Service

	if got := s.connectTimeout(); got != DefaultConnectTimeout {
		t.Errorf("connectTimeout = %v, want %v", got, DefaultConnectTimeout)
	}
	if got := s.sendTimeout(); got != DefaultSendTimeout {
		t.Errorf("sendTimeout = %v, want %v", got, DefaultSendTimeout)
	}
	if got := s.errorChanGrace(); got != DefaultErrorChanGrace {
		t.Errorf("errorChanGrace = %v, want %v", got, DefaultErrorChanGrace)
	}

	want := DefaultConnectTimeout + DefaultSendTimeout
	if got := s.MaxSendDuration(); got != want {
		t.Errorf("MaxSendDuration = %v, want %v", got, want)
	}
}

func Test_TimeoutOverrides(t *testing.T) {
	s := Service{
		ConnectTimeout: 3 * time.Second,
		SendTimeout:    7 * time.Second,
		ErrorChanGrace: time.Second,
	}

	if got := s.connectTimeout(); got != 3*time.Second {
		t.Errorf("connectTimeout = %v, want 3s", got)
	}
	if got := s.sendTimeout(); got != 7*time.Second {
		t.Errorf("sendTimeout = %v, want 7s", got)
	}
	if got := s.errorChanGrace(); got != time.Second {
		t.Errorf("errorChanGrace = %v, want 1s", got)
	}
	if got := s.MaxSendDuration(); got != 10*time.Second {
		t.Errorf("MaxSendDuration = %v, want 10s", got)
	}
}
