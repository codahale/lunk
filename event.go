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

// EventID is the ID of an event and its root event.
type EventID struct {
	// Root is the root ID of the tree which contains all of the events related
	// to this one.
	Root ID `json:"root"`

	// ID is an ID uniquely identifying the event.
	ID ID `json:"id"`
}

// Metadata is a collection of metadata about an Event.
type Metadata struct {
	EventID

	// Schema is the schema of the event.
	Schema string `json:"schema"`

	// Parent is the ID of the parent event, if any.
	Parent ID `json:"parent,omitempty"`

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

// An Entry is the combination of an event and its metadata.
type Entry struct {
	Metadata

	// Event is the actual event object, to be serialized as a object.
	Event Event `json:"event"`
}

// NewRootEntry creates a new Entry instance for the root event in the given
// tree. This should only be used to generate entries for events caused
// exclusively by events which are outside of your system as a whole (e.g., a
// root entry for the first time you see a user request).
func NewRootEntry(e Event) *Entry {
	id := NewID()
	return &Entry{
		Metadata: Metadata{
			EventID: EventID{
				Root: id,
				ID:   id,
			},
			Schema: e.Schema(),
			Time:   time.Now(),
			Host:   host,
			Deploy: deploy,
			PID:    pid,
		},
		Event: e}
}

// NewEntry creates a new Entry instance for the given event in the given tree
// with the given parent.
func NewEntry(parent EventID, e Event) *Entry {
	return &Entry{
		Metadata: Metadata{
			EventID: EventID{
				Root: parent.Root,
				ID:   NewID(),
			},
			Schema: e.Schema(),
			Parent: parent.ID,
			Time:   time.Now(),
			Host:   host,
			Deploy: deploy,
			PID:    pid,
		},
		Event: e}
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
