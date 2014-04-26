package lunk

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestJSONEventLogger(t *testing.T) {
	root := NewID()
	parent := NewID()
	ev := mockEvent{Example: "whee"}

	buf := bytes.NewBuffer(nil)
	logger := NewJSONEventLogger(buf)
	id := logger.Log(root, parent, ev)

	var e RawJSONEntry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatal(err)
	}

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Root != root {
		t.Errorf("Root was %v, but expected %v", e.Root, root)
	}

	if e.ID != id {
		t.Errorf("ID was %v, but expected %v", e.ID, id)
	}

	if e.Parent != parent {
		t.Errorf("Parent was %v, but expected %v", e.Parent, parent)
	}

	if time.Now().Sub(e.Time) > 5*time.Millisecond {
		t.Errorf("Unexpectedly old timestamp: %v", e.Time)
	}

	if e.Host == "" {
		t.Errorf("Blank hostname for meta data")
	}

	var newEV mockEvent
	if err := json.Unmarshal(e.Event, &newEV); err != nil {
		t.Fatal(err)
	}

	if newEV.Example != "whee" {
		t.Errorf("Bad event data: %v", newEV)
	}
}
