package lunk

import (
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	e := NewEvent("request")

	if e.Type != "request" {
		t.Errorf("Unexpected event type: %v", e.Type)
	}

	if e.ID == 0 {
		t.Errorf("Zero ID")
	}

	if e.Time.Sub(time.Now()) > 1*time.Millisecond {
		t.Errorf("Unexpectedly old timestamp: %v", e.Time)
	}

	if e.Properties == nil {
		t.Errorf("Nil properties")
	}
}

func TestEventAddString(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", "two")

	expected := map[string]string{
		"one": "two",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventAddFloat32(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", float32(2.2))

	expected := map[string]string{
		"one": "2.20",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventAddFloat64(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", float64(2.2))

	expected := map[string]string{
		"one": "2.20",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventAddInt(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", 20)

	expected := map[string]string{
		"one": "20",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventAddStringer(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", mockStringer(0))

	expected := map[string]string{
		"one": "mock!",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventAddAnything(t *testing.T) {
	e := NewEvent("blah")
	e.Add("one", mockNonStringer(0))

	expected := map[string]string{
		"one": "0",
	}
	if !reflect.DeepEqual(e.Properties, expected) {
		t.Errorf("Was %v, but expected %v", e.Properties, expected)
	}
}

func TestEventChildEvent(t *testing.T) {
	p := NewEvent("blah")
	c := p.ChildEvent("woo")

	if p.ID.String() != c.Properties["parent"] {
		t.Errorf("No parent ID set: %v", c)
	}
}

type mockStringer int

func (mockStringer) String() string {
	return "mock!"
}

type mockNonStringer int

func TestEventString(t *testing.T) {
	e := Event{
		Type:     "type",
		Message:  "msg",
		ID:       ID(100),
		ParentID: ID(150),
		Time:     time.Date(2014, 4, 23, 12, 15, 32, 0, time.UTC),
		Properties: map[string]string{
			"yay": "whee",
		},
	}

	actual := e.String()
	expected := `msg type="type" id="0000000000000064" parent="0000000000000096" time="2014-04-23T12:15:32Z" yay="whee" `

	if actual != expected {
		t.Errorf("Was `%s` but expected `%s`", actual, expected)
	}
}
