package lunk

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"
)

// An Event is a record of the occurrence of something.
type Event interface {
	// Schema returns the schema of the event. This should be constant.
	Schema() string
}

var (
	// ErrBadEventID is returned when the event ID cannot be parsed.
	ErrBadEventID = errors.New("bad event ID")
)

// EventID is the ID of an event, its parent event, and its root event.
type EventID struct {
	// Root is the root ID of the tree which contains all of the events related
	// to this one.
	Root ID `json:"root"`

	// ID is an ID uniquely identifying the event.
	ID ID `json:"id"`

	// Parent is the ID of the parent event, if any.
	Parent ID `json:"parent,omitempty"`
}

// String returns the EventID as a URL-encoded set of parameters.
func (id EventID) String() string {
	return fmt.Sprintf("root=%s&id=%s", id.Root, id.ID)
}

// Format formats according to a format specifier and returns the resulting
// string. The receiver's string representation is the first argument.
func (id EventID) Format(s string, args ...interface{}) string {
	args = append([]interface{}{id.String()}, args...)
	return fmt.Sprintf(s, args...)
}

// NewRootEventID generates a new event ID for a root event. This should only be
// used to generate entries for events caused exclusively by events which are
// outside of your system as a whole (e.g., a root event for the first time you
// see a user request).
func NewRootEventID() EventID {
	id := NewID()
	return EventID{
		Root: id,
		ID:   id,
	}
}

// NewEventID returns a new ID for an event which is the child of the given
// parent ID. This should be used to track causal relationships between events.
func NewEventID(parent EventID) EventID {
	return EventID{
		Root:   parent.Root,
		ID:     NewID(),
		Parent: parent.ID,
	}
}

// ParseEventID parses the given string as a URL-encoded set of parameters.
func ParseEventID(s string) (*EventID, error) {
	u, err := url.ParseQuery(s)
	if err != nil {
		return nil, ErrBadEventID
	}

	root, err := ParseID(u.Get("root"))
	if err != nil {
		return nil, ErrBadEventID
	}

	id, err := ParseID(u.Get("id"))
	if err != nil {
		return nil, ErrBadEventID
	}

	return &EventID{
		Root: root,
		ID:   id,
	}, nil
}

// Metadata is a collection of metadata about an Event.
type Metadata struct {
	// Schema is the schema of the event.
	Schema string `json:"schema"`

	// Time is the timestamp of the event.
	Time time.Time `json:"time"`

	// Host is the name of the host on which the event occurred.
	Host string `json:"host,omitempty"`

	// Deploy is the ID of the deployed artifact, read from the DEPLOY
	// environment variable on startup.
	Deploy string `json:"deploy,omitempty"`

	// PID is the process ID which generated the event.
	PID int `json:"pid"`
}

// NewMetdata returns a populated Metadata instance for the given event.
func NewMetadata(e Event) Metadata {
	return Metadata{
		Schema: e.Schema(),
		Time:   time.Now(),
		Host:   host,
		Deploy: deploy,
		PID:    pid,
	}
}

// An Entry is the combination of an event and its metadata.
type Entry struct {
	EventID
	Metadata

	// Event is the actual event object, to be serialized to JSON.
	Event Event `json:"event"`
}

var (
	host, deploy string
	pid          int
)

func init() {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	host = h
	deploy = os.Getenv("DEPLOY")
	pid = os.Getpid()
}
