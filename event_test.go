package lunk

import (
	"testing"
	"time"
)

var (
	root   = ID(1)
	parent = ID(2)
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

	if e.Event != ev {
		t.Errorf("Unexpected properties: %v", e.Event)
	}
}

func TestNewEntry(t *testing.T) {
	ev := mockEvent{Example: "yay"}
	e := NewEntry(root, parent, ev)

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Root != root {
		t.Errorf("Unexpected root ID: %v", e.Root)
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
