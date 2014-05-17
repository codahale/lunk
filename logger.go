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

type jsonEventLogger struct {
	*json.Encoder
}

func (l jsonEventLogger) Log(id EventID, e Event) {
	if err := l.Encode(NewEntry(id, e)); err != nil {
		panic(err)
	}
}
