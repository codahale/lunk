package lunk

import (
	"encoding/json"
	"net/http"
	"testing"
)

var _ Event = HTTPRequestEvent{}

func TestSetRootID(t *testing.T) {
	root := ID(100)
	parent := ID(150)

	r := http.Request{
		Header: http.Header{},
	}

	SetRootID(&r, root, parent)

	actual := r.Header.Get("Root-ID")
	expected := "0000000000000064:0000000000000096"
	if actual != expected {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func TestGetRootID(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}
	r.Header.Add("Root-ID", "0000000000000064:0000000000000096")

	root, parent, err := GetRootID(&r)
	if root != 100 {
		t.Errorf("Unexpected root ID: %v", root)
	}

	if parent != 150 {
		t.Errorf("Unexpected parent ID: %v", parent)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetRootIDMissing(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}

	root, parent, err := GetRootID(&r)
	if root != 0 {
		t.Errorf("Unexpected root ID: %#v", root)
	}

	if parent != 0 {
		t.Errorf("Unexpected parent ID: %#v", parent)
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetRootIDBadRoot(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}
	r.Header.Add("Root-ID", "000g000000000064:0000000000000096")

	root, parent, err := GetRootID(&r)
	if root != 0 {
		t.Errorf("Unexpected root ID: %v", root)
	}

	if parent != 0 {
		t.Errorf("Unexpected parent ID: %v", parent)
	}

	if err.Error() != `strconv.ParseUint: parsing "000g000000000064": invalid syntax` {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetRootIDBadParent(t *testing.T) {
	r := http.Request{
		Header: http.Header{},
	}
	r.Header.Add("Root-ID", "000000000000064:0000000g000000096")

	root, parent, err := GetRootID(&r)
	if root != 100 {
		t.Errorf("Unexpected root ID: %v", root)
	}

	if parent != 0 {
		t.Errorf("Unexpected parent ID: %v", parent)
	}

	if err.Error() != `strconv.ParseUint: parsing "0000000g000000096": invalid syntax` {
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
