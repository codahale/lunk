package lunk

import (
	"bytes"
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestJSONEventLoggerLog(t *testing.T) {
	ev := mockEvent{Example: "whee"}

	buf := bytes.NewBuffer(nil)
	logger := NewJSONEventLogger(buf)
	id := NewRootEventID()
	logger.Log(id, ev)

	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatal(err)
	}

	if e.Schema != "example" {
		t.Errorf("Unexpected schema: %v", e.Schema)
	}

	if e.Root != id.Root {
		t.Errorf("Root was %v, but expected %v", e.Root, id.Root)
	}

	if e.ID != id.ID {
		t.Errorf("ID was %v, but expected %v", e.ID, id.ID)
	}

	if e.Parent != 0 {
		t.Errorf("Parent was %v, but expected %v", e.Parent, 0)
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

	expected := map[string]string{
		"example": "whee",
	}
	actual := e.Properties
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Properties were %+v, expected %+v", actual, expected)
	}
}

func TestTextLogger(t *testing.T) {
	ev := mockEvent{Example: "whee"}

	buf := bytes.NewBuffer(nil)
	logger := NewTextEventLogger(buf)
	id := EventID{
		Root:   100,
		Parent: 150,
		ID:     200,
	}
	logger.Log(id, ev)

	expected := regexp.MustCompile(
		`^time="[\d]{4}-[\d]{2}-[\d]{2}T[\d]{2}:[\d]{2}:[\d]{2}Z"` +
			` host="[^"]+"` +
			` pid="[\d]+"` +
			` deploy="[^"]*"` +
			` schema="example"` +
			` id="00000000000000c8"` +
			` root="0000000000000064"` +
			` parent="0000000000000096"` +
			` p:example="whee"` +
			`\n$`,
	)
	actual := buf.String()

	if !expected.MatchString(actual) {
		t.Errorf("Was `%s` but expected to match `%s`", actual, expected)
	}
}

func TestTextLoggerElidedParentID(t *testing.T) {
	ev := mockEvent{Example: "whee"}

	buf := bytes.NewBuffer(nil)
	logger := NewTextEventLogger(buf)
	id := EventID{
		Root: 100,
		ID:   200,
	}
	logger.Log(id, ev)

	expected := regexp.MustCompile(
		`^time="[\d]{4}-[\d]{2}-[\d]{2}T[\d]{2}:[\d]{2}:[\d]{2}Z"` +
			` host="[^"]+"` +
			` pid="[\d]+"` +
			` deploy="[^"]*"` +
			` schema="example"` +
			` id="00000000000000c8"` +
			` root="0000000000000064"` +
			` p:example="whee"` +
			`\n$`,
	)
	actual := buf.String()

	if !expected.MatchString(actual) {
		t.Errorf("Was `%s` but expected to match `%s`", actual, expected)
	}
}

type nullWriter struct{}

func (nullWriter) Write(a []byte) (int, error) {
	return len(a), nil
}

func BenchmarkJSONEventLogger(b *testing.B) {
	ev := mockEvent{Example: "whee"}
	logger := NewJSONEventLogger(nullWriter{})
	id := NewRootEventID()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Log(id, ev)
	}
}
