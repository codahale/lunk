package lunk

import (
	"testing"
	"time"
)

var (
	span   = ID(1)
	parent = ID(2)
)

type mockEvent struct {
	Example string `json:"example"`
}

func (mockEvent) Schema() string {
	return "example"
}

func TestNewEventMetadata(t *testing.T) {
	ev := mockEvent{Example: "yay"}
	e := NewEntry(span, parent, ev)

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Span != span {
		t.Errorf("Unexpected span: %v", e.Span)
	}

	if e.ID == 0 {
		t.Errorf("Zero ID for entry")
	}

	if e.Parent != parent {
		t.Errorf("Unexpected parent: %v", e.Parent)
	}

	if time.Now().Sub(e.Time) > 5*time.Millisecond {
		t.Errorf("Unexpectedly old timestamp: %v", e.Time)
	}

	if e.Host == "" {
		t.Errorf("Blank hostname for meta data")
	}

	if e.Event != ev {
		t.Errorf("Unexpected properties: %v", e.Event)
	}
}
