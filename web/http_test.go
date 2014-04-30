package web

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/codahale/lunk"
)

var _ lunk.Event = HTTPRequestEvent{}

func TestSetRequestEventID(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}

	SetRequestEventID(&r, lunk.EventID{
		Root: 100,
		ID:   150,
	})

	actual := r.Header.Get("Event-ID")
	expected := "root=0000000000000064&id=0000000000000096"
	if actual != expected {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func TestGetRequestEventID(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}
	r.Header.Add("Event-ID", "root=0000000000000064&id=0000000000000096")

	id, err := GetRequestEventID(&r)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if id.Root != 100 || id.ID != 150 {
		t.Errorf("Unexpected event ID: %+v", id)
	}
}

func TestGetRequestEventIDMissing(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}

	id, err := GetRequestEventID(&r)

	if id != nil {
		t.Errorf("Unexpected event ID: %+v", id)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestHTTPRequest(t *testing.T) {
	r := &http.Request{
		Method:        "GET",
		RequestURI:    "/woohoo",
		Proto:         "HTTP/1.1",
		RemoteAddr:    "127.0.0.1",
		ContentLength: 0,
		Header: http.Header{
			"Authorization": []string{"Basic seeecret"},
			"Accept":        []string{"application/json"},
		},
		Trailer: http.Header{
			"Authorization": []string{"Basic seeecret"},
			"Connection":    []string{"close"},
		},
	}

	e := HTTPRequest(r)
	e.Status = 200
	e.ElapsedMillis = 4.3

	if e.Schema() != "httprequest" {
		t.Errorf("Unexpected schema: %v", e.Schema())
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(data)
	expected := `{"method":"GET","uri":"/woohoo","proto":"HTTP/1.1","headers":{"Accept":["application/json"],"Authorization":["REDACTED"],"Connection":["close"]},"host":"","remote_addr":"127.0.0.1","content_length":0,"status":200,"elapsed_ms":4.3}`
	if actual != expected {
		t.Errorf("Was `%s`, but expected `%v`", actual, expected)
	}
}
