package lunk

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestJSONEventLoggerLog(t *testing.T) {
	ev := mockEvent{Example: "whee"}

	buf := bytes.NewBuffer(nil)
	logger := NewJSONEventLogger(buf)
	id := NewRootEventID()
	logger.Log(id, ev)

	var e RawJSONEntry
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

	var newEV mockEvent
	if err := json.Unmarshal(e.Event, &newEV); err != nil {
		t.Fatal(err)
	}

	if newEV.Example != "whee" {
		t.Errorf("Bad event data: %v", newEV)
	}
}

var flattenableJSON = `
{
  "root":"03c9c7ed185b429d",
  "id":"9c0b8e3d54a51e8b",
  "schema":"example",
  "time":"2014-05-13T20:34:13.802338093-07:00",
  "host":"lunchbox-hd.local",
  "pid":18826,
  "event":{
    "b": false,
    "n": null,
    "one":1,
    "two": 2.2,
    "three": [3, "3", 3.0],
    "four": {
      "five": "5",
      "six": 6
    }
  }
}
`

func TestFlatJSONEntry(t *testing.T) {
	var e FlatJSONEntry
	if err := json.Unmarshal([]byte(flattenableJSON), &e); err != nil {
		t.Fatal(e)
	}

	actual := make(map[string]string)
	e.Flatten(func(k, v string) {
		actual[k] = v
	})

	expected := map[string]string{
		"n":         "null",
		"one":       "1",
		"four.six":  "6",
		"four.five": "5",
		"b":         "false",
		"two":       "2.2",
		"three.0":   "3",
		"three.1":   "3",
		"three.2":   "3",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func BenchmarkFlatJSONEntryFlatten(b *testing.B) {
	var e FlatJSONEntry
	if err := json.Unmarshal([]byte(flattenableJSON), &e); err != nil {
		b.Fatal(e)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e.Flatten(func(k, v string) {
			n += len(k) + len(v)
		})
	}
}
