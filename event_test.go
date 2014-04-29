package lunk

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestEventIDString(t *testing.T) {
	id := EventID{
		Root: 100,
		ID:   200,
	}

	actual := id.String()
	expected := "0000000000000064/00000000000000c8"
	if actual != expected {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func TestEventIDFormat(t *testing.T) {
	id := EventID{
		Root: 100,
		ID:   200,
	}

	actual := id.Format("/* %s/%s */ %s", "SELECT 1")
	expected := "/* 0000000000000064/00000000000000c8 */ SELECT 1"
	if actual != expected {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func ExampleEventID_Format() {
	// Assume we're connected to a database.
	var (
		event  EventID
		db     *sql.DB
		userID int
	)

	// This passes the root ID and the parent event ID to the database, which
	// allows us to correlate, for example, slow queries with the web requests
	// which caused them.
	query := event.Format(`/* %s/%s */ %s`, `SELECT email FROM users WHERE id = ?`)
	r := db.QueryRow(query, userID)
	if r == nil {
		panic("user not found")
	}

	var email string
	if err := r.Scan(&email); err != nil {
		panic("couldn't read email")
	}

	fmt.Printf("User's email: %s\n", email)
}

func TestParseEventID(t *testing.T) {
	id, err := ParseEventID("0000000000000064/0000000000000096")

	if id.Root != 100 || id.ID != 150 {
		t.Errorf("Unexpected event ID: %+v", id)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseEventIDBadRoot(t *testing.T) {
	id, err := ParseEventID("000g000000000064/0000000000000096")

	if id != nil {
		t.Errorf("Unexpected event ID: %+v", id)
	}

	if err != ErrBadEventID {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseEventIDBadID(t *testing.T) {
	id, err := ParseEventID("000000000000064/0000000g000000096")

	if id != nil {
		t.Errorf("Unexpected event ID: %+v", id)
	}

	if err != ErrBadEventID {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestParseEventIDTooFew(t *testing.T) {
	id, err := ParseEventID("000000000000064")

	if id != nil {
		t.Errorf("Unexpected event ID: %+v", id)
	}

	if err != ErrBadEventID {
		t.Errorf("Unexpected error: %v", err)
	}
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

type mockEvent struct {
	Example string `json:"example"`
}

func (mockEvent) Schema() string {
	return "example"
}
