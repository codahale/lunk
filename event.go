package lunk

import (
	"os"
	"time"
)

// An Event is a record of the occurrence of something.
type Event interface {
	// Schema returns the schema of the event. This should be constant.
	Schema() string
}

// Metadata is a collection of metadata about an Event.
type Metadata struct {
	// Schema is the schema of the event.
	Schema string `json:"schema"`

	// Span is the ID of the span which contains all of this event's siblings.
	Span ID `json:"span"`

	// ID is an ID uniquely identifying the event.
	ID ID `json:"id"`

	// Parent is the ID of the parent event, if any.
	Parent ID `json:"parent,omitempty"`

	// Time is the timestamp of the event.
	Time time.Time `json:"time"`

	// Host is the name of the host on which the event occurred.
	Host string `json:"host,omitempty"`

	// Deploy is the ID of the deployed artifact, read from the DEPLOY
	// environment variable on startup.
	Deploy string `json:"deploy,omitempty"`

	// Event is the actual event object, to be serialized as a nested
	// object.
	Event Event `json:"event"`
}

// An Entry is the combination of an event and its metadata.
type Entry struct {
	Metadata

	// Event is the actual event object, to be serialized as a object.
	Event Event `json:"event"`
}

// NewEntry creates a new Entry instance for the given event.
func NewEntry(span, parent ID, e Event) *Entry {
	return &Entry{Metadata: Metadata{
		Schema: e.Schema(),
		Span:   span,
		ID:     NewID(),
		Parent: parent,
		Time:   time.Now(),
		Host:   host,
		Deploy: deploy,
		Event:  e,
	},
		Event: e}
}

var (
	host, deploy string
)

func init() {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	host = h
	deploy = os.Getenv("DEPLOY")
}
