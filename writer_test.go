package lunk

import (
	"bytes"
	"testing"
	"time"
)

func TestTextEventWriter(t *testing.T) {
	e := &Event{
		Type:     "type",
		Message:  "msg",
		ID:       ID(100),
		ParentID: ID(150),
		Time:     time.Date(2014, 4, 23, 12, 15, 32, 0, time.UTC),
		Properties: map[string]string{
			"yay": "whee",
		},
	}

	buf := bytes.NewBuffer(nil)
	w := NewTextEventWriter(buf)
	if err := w.Write(e); err != nil {
		t.Fatal(err)
	}

	actual := buf.String()
	expected := `msg type="type" id="0000000000000064" parent="0000000000000096" time="2014-04-23T12:15:32Z" yay="whee" ` + "\n"
	if actual != expected {
		t.Errorf("Was `%s`, but expected `%s`", actual, expected)
	}
}

func TestJSONEventWriter(t *testing.T) {
	e := &Event{
		Type:     "type",
		Message:  "msg",
		ID:       ID(100),
		ParentID: ID(150),
		Time:     time.Date(2014, 4, 23, 12, 15, 32, 0, time.UTC),
		Properties: map[string]string{
			"yay": "whee",
		},
	}

	buf := bytes.NewBuffer(nil)
	w := NewJSONEventWriter(buf)
	if err := w.Write(e); err != nil {
		t.Fatal(err)
	}

	actual := buf.String()
	expected := `{"type":"type","msg":"msg","id":"0000000000000064","parent":"0000000000000096","t":"2014-04-23T12:15:32Z","p":{"yay":"whee"}}` + "\n"
	if actual != expected {
		t.Errorf("Was `%s`, but expected `%s`", actual, expected)
	}
}

func TestAsyncEventWriter(t *testing.T) {
	e := &Event{
		Type:     "type",
		Message:  "msg",
		ID:       ID(100),
		ParentID: ID(150),
		Time:     time.Date(2014, 4, 23, 12, 15, 32, 0, time.UTC),
		Properties: map[string]string{
			"yay": "whee",
		},
	}

	buf := bytes.NewBuffer(nil)
	w := NewJSONEventWriter(buf)
	a := NewAsyncEventWriter(w, 100, func(e *Event, err error) {
		t.Errorf("Error handling %v: %v", e, err)
	})
	a.Start()
	if err := w.Write(e); err != nil {
		t.Fatal(err)
	}
	a.Stop()

	actual := buf.String()
	expected := `{"type":"type","msg":"msg","id":"0000000000000064","parent":"0000000000000096","t":"2014-04-23T12:15:32Z","p":{"yay":"whee"}}` + "\n"
	if actual != expected {
		t.Errorf("Was `%s`, but expected `%s`", actual, expected)
	}
}

func BenchmarkJSONEventWriter(b *testing.B) {
	e := NewEvent("http.request")
	e.Add("one", "two")
	e.Add("three", "four five")

	w := NewJSONEventWriter(nullWriter(0))
	for i := 0; i < b.N; i++ {
		w.Write(e)
	}
}

func BenchmarkAsyncJSONEventWriter(b *testing.B) {
	e := NewEvent("http.request")
	e.Add("one", "two")
	e.Add("three", "four five")

	w := NewJSONEventWriter(nullWriter(0))
	a := NewAsyncEventWriter(w, 100, func(e *Event, err error) {
		b.Errorf("Error handling %v: %v", e, err)
	})
	a.Start()
	for i := 0; i < b.N; i++ {
		w.Write(e)
	}
	a.Stop()
}

func BenchmarkTextEventFormatter(b *testing.B) {
	e := NewEvent("http.request")
	e.Add("one", "two")
	e.Add("three", "four five")

	w := NewTextEventWriter(nullWriter(0))
	for i := 0; i < b.N; i++ {
		w.Write(e)
	}
}

type nullWriter int

func (nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
