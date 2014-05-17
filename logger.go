package lunk

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
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

// NewTextEventLogger returns an EventLogger which writes entries as single
// lines of attr="value" formatted text.
func NewTextEventLogger(w io.Writer) EventLogger {
	return textEventLogger{w: w}
}

type textEventLogger struct {
	w io.Writer
}

func (l textEventLogger) Log(id EventID, e Event) {
	entry := NewEntry(id, e)

	props := []string{
		fmt.Sprintf("time=%s", strconv.Quote(entry.Time.Format(time.RFC3339))),
		fmt.Sprintf("host=%s", strconv.Quote(entry.Host)),
		fmt.Sprintf(`pid="%d"`, entry.PID),
		fmt.Sprintf("deploy=%s", strconv.Quote(entry.Deploy)),
		fmt.Sprintf("schema=%s", strconv.Quote(entry.Schema)),
		fmt.Sprintf("id=%s", strconv.Quote(entry.ID.String())),
		fmt.Sprintf("root=%s", strconv.Quote(entry.Root.String())),
	}

	if entry.Parent != 0 {
		s := fmt.Sprintf("parent=%s", strconv.Quote(entry.Parent.String()))
		props = append(props, s)
	}

	for _, k := range sortedKeys(entry.Properties) {
		s := fmt.Sprintf("p:%s=%s", k, strconv.Quote(entry.Properties[k]))
		props = append(props, s)
	}

	if _, err := fmt.Fprintln(l.w, strings.Join(props, " ")); err != nil {
		panic(err)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type jsonEventLogger struct {
	*json.Encoder
}

func (l jsonEventLogger) Log(id EventID, e Event) {
	if err := l.Encode(NewEntry(id, e)); err != nil {
		panic(err)
	}
}
