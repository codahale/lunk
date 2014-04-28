package lunk

import (
	"encoding/json"
	"io"
)

// An EventLogger logs events and their metadata.
type EventLogger interface {
	// Log adds the given event to the log stream, returning the logged event's
	// unique ID.
	Log(root, parent ID, e Event) ID

	// LogRoot adds the given root event to the log stream, returning the logged
	// event's unique ID. This should only be used to generate entries for events
	// caused exclusively by events which are outside of your system as a whole
	// (e.g., a root entry for the first time you see a user request).
	LogRoot(e Event) ID
}

// NewJSONEventLogger returns an EventLogger which writes entries as streaming
// JSON to the given writer.
func NewJSONEventLogger(w io.Writer) EventLogger {
	return jsonEventLogger{json.NewEncoder(w)}
}

// RawJSONEntry allows entries to be parsed without any knowledge of the schemas.
type RawJSONEntry struct {
	Metadata

	// Event is the unparsed event data.
	Event json.RawMessage `json:"event"`
}

type jsonEventLogger struct {
	*json.Encoder
}

func (l jsonEventLogger) Log(root, parent ID, e Event) ID {
	entry := NewEntry(root, parent, e)
	if err := l.Encode(entry); err != nil {
		panic(err)
	}
	return entry.ID
}

func (l jsonEventLogger) LogRoot(e Event) ID {
	entry := NewRootEntry(e)
	if err := l.Encode(entry); err != nil {
		panic(err)
	}
	return entry.ID
}
