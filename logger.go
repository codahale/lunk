package lunk

import (
	"encoding/json"
	"io"
)

// An EventLogger logs events and their metadata.
type EventLogger interface {
	// Log adds the given event to the log stream.
	Log(id EventID, e Event)
}

// NewJSONEventLogger returns an EventLogger which writes entries as streaming
// JSON to the given writer.
func NewJSONEventLogger(w io.Writer) EventLogger {
	return jsonEventLogger{json.NewEncoder(w)}
}

// RawJSONEntry allows entries to be parsed without any knowledge of the schemas.
type RawJSONEntry struct {
	EventID
	Metadata

	// Event is the unparsed event data.
	Event json.RawMessage `json:"event"`
}

type jsonEventLogger struct {
	*json.Encoder
}

func (l jsonEventLogger) Log(id EventID, e Event) {
	entry := Entry{
		EventID:  id,
		Metadata: NewMetadata(e),
		Event:    e,
	}
	if err := l.Encode(entry); err != nil {
		panic(err)
	}
}
