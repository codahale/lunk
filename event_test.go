package lunk

import (
	"testing"
	"time"
)

type mockEvent struct {
	Example string `json:"example"`
}

func (mockEvent) Schema() string {
	return "example"
}

func TestNewRootEntry(t *testing.T) {
	ev := mockEvent{Example: "yay"}
	e := NewRootEntry(ev)

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Root != e.ID {
		t.Errorf("Unexpected root ID: %v", e.Root)
	}

	if e.ID == 0 {
		t.Errorf("Zero ID for entry")
	}

	if e.Parent != 0 {
		t.Errorf("Unexpected parent: %v", e.Parent)
	}

	if time.Now().Sub(e.Time) > 5*time.Millisecond {
		t.Errorf("Unexpectedly old timestamp: %v", e.Time)
	}

	if e.Host == "" {
		t.Errorf("Blank hostname for meta data")
	}

	if e.PID == 0 {
		t.Errorf("Blank PID for meta data")
	}

	if e.Event != ev {
		t.Errorf("Unexpected properties: %v", e.Event)
	}
}

func TestNewEntry(t *testing.T) {
	id := EventID{
		Root: 1,
		ID:   2,
	}
	ev := mockEvent{Example: "yay"}
	e := NewEntry(id, ev)

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Root != id.Root {
		t.Errorf("Unexpected root ID: %v", e.Root)
	}

	if e.ID == 0 {
		t.Errorf("Zero ID for entry")
	}

	if e.Parent != id.ID {
		t.Errorf("Unexpected parent: %v", e.Parent)
	}

	if time.Now().Sub(e.Time) > 5*time.Millisecond {
		t.Errorf("Unexpectedly old timestamp: %v", e.Time)
	}

	if e.Host == "" {
		t.Errorf("Blank hostname for meta data")
	}

	if e.PID == 0 {
		t.Errorf("Blank PID for meta data")
	}

	if e.Event != ev {
		t.Errorf("Unexpected properties: %v", e.Event)
	}
}
