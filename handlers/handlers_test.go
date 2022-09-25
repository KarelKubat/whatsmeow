package handlers

import (
	"errors"
	"sync"
	"testing"

	"go.mau.fi/whatsmeow/types/events"
)

// TestEventTypeString checks that there are strings for all event types. If not, the test barfs with
// an index-out-of-range error.
func TestEventTypeString(t *testing.T) {
	for tp := firstEventType + 1; tp < lastEventType; tp++ {
		t.Log(int(tp), tp.String())
	}
}

// TestDispatchErrorTypeString checks that there are strings for all dispatcher error types.
func TestDispatchErrorTypeString(t *testing.T) {
	for de := firstDispatchError + 1; de < lastDispatchError; de++ {
		t.Log(int(de), de.String())
	}
}

type dummyHandler struct{}

func (d *dummyHandler) Handle(ev interface{}) error { return errors.New("fail") }

// TestAsyncRegistration checks that in-parallel registration doesn't break.
func TestAsyncRegistration(t *testing.T) {
	registry = make(map[EventType][]handler)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Register(UnknownCallEvent, &dummyHandler{})
		}()
	}
	wg.Wait()
	if l := len(registry[UnknownCallEvent]); l != 1000 {
		t.Errorf("TestAsyncRegistration: %v handlers registered, want 1000", l)
	}
}

// TestDispatchError checks that Dispatch() returns a correct error type.
func TestDispatchError(t *testing.T) {
	registry = make(map[EventType][]handler)
	Register(UndecryptableMessage, &dummyHandler{})

	for _, test := range []struct {
		description   string
		event         interface{}
		wantErrorType dispatchErrorType
	}{
		{
			description:   "unregistered handler",
			event:         &events.AppState{},
			wantErrorType: NoHandlerFound,
		},
		{
			description:   "handler error",
			event:         &events.UndecryptableMessage{},
			wantErrorType: HandlerFailed,
		},
		{
			description:   "unknown event",
			event:         &struct{}{},
			wantErrorType: UnknownEvent,
		},
	} {
		err := Dispatch(test.event)
		if err == nil {
			t.Fatalf("%v: Dispatch(_) = nil, need error", test.description)
		}
		if err.Type != test.wantErrorType {
			t.Errorf("%v: Dispatch(_) = %v, want type %v", test.description, err, test.wantErrorType)
		}
	}
}
