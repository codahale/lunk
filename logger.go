package lunk

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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

// FlatJSONEntry allows entry properties to be flatten into dotted-path keys and
// stringified values. Nested objects are given keys like "inner.key", and
// nested arrays are given keys like "array.0".
type FlatJSONEntry struct {
	EventID
	Metadata

	// Event is the parsed event data.
	Event map[string]interface{} `json:"event"`
}

// Flatten calls the given function with every flattened property key and value.
func (e FlatJSONEntry) Flatten(f func(k, v string)) {
	for k, v := range e.Event {
		flatten(k, v, f)
	}
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

func flatten(k string, v interface{}, f func(k, v string)) {
	switch v := v.(type) {
	case float64: // all numeric types are considered float64s, sadly
		f(k, strconv.FormatFloat(v, 'f', -1, 64))
	case []interface{}:
		for i, v2 := range v {
			flatten(k+"."+strconv.Itoa(i), v2, f)
		}
	case map[string]interface{}:
		for k2, v2 := range v {
			flatten(k+"."+k2, v2, f)
		}
	case bool:
		f(k, strconv.FormatBool(v))
	case nil:
		f(k, "null")
	case string:
		f(k, v)
	default:
		f(k, fmt.Sprintf("%v", v))
	}
}
