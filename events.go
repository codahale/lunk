package lunk

import (
	"fmt"
	"strconv"
	"time"
)

// An Event is a record of the occurrence of something.
type Event struct {
	// Type is the type of the event.
	Type string `json:"type"`

	// Message is an optional human-readable message.
	Message string `json:"msg,omitempty"`

	// ID is a unique ID.
	ID ID `json:"id"`

	// ParentID is the ID of the parent event, if any.
	ParentID ID `json:"parent,omitempty"`

	// Time is the timestamp of the event.
	Time time.Time `json:"t"`

	// Properties are arbitrary named values.
	Properties map[string]string `json:"p"`
}

// NewEvent returns a new event of the given type.
func NewEvent(eventType string) *Event {
	e := &Event{
		ID:         NewID(),
		Time:       time.Now(),
		Type:       eventType,
		Properties: make(map[string]string),
	}
	return e
}

// Add associates a named value with the event.
func (e *Event) Add(key string, value interface{}) {
	switch value := value.(type) {
	case string:
		e.Properties[key] = value
	case float32, float64:
		e.Properties[key] = fmt.Sprintf("%.2f", value)
	case int:
		e.Properties[key] = strconv.Itoa(value)
	case fmt.Stringer:
		e.Properties[key] = value.String()
	default:
		e.Properties[key] = fmt.Sprintf("%v", value)
	}
}

// ChildEvent returns a child event of the given type.
func (e *Event) ChildEvent(eventType string) *Event {
	child := NewEvent(eventType)
	child.Add("parent", e.ID)
	return child
}

// String returns the event as a single-line log entry of the event's message,
// followed by the properties, formatted as key="value".
func (e Event) String() string {
	s := ""
	if e.Message != "" {
		s += e.Message
		s += " "
	}

	s += fmt.Sprintf(`type=%s `, strconv.Quote(e.Type))
	s += fmt.Sprintf(`id=%s `, strconv.Quote(e.ID.String()))
	if e.ParentID != 0 {
		s += fmt.Sprintf(`parent=%s `, strconv.Quote(e.ParentID.String()))
	}

	s += fmt.Sprintf(`time=%s `, strconv.Quote(e.Time.Format(time.RFC3339)))

	for k, v := range e.Properties {
		s += fmt.Sprintf(`%s=%s `, k, strconv.Quote(v))
	}

	return s
}
