package lunk

import (
	"encoding/json"
	"testing"
)

func TestMessage(t *testing.T) {
	e := Message("whee")

	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}

	actual := string(j)
	expected := `{"msg":"whee"}`

	if actual != expected {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}
