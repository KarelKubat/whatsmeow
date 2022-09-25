package handlers

import "testing"

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
