package lunk

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// RootID is the name of the HTTP header by which the root and parent IDs are
	// passed to internal services.
	RootID = "Root-ID"
)

// SetRootID sets the Root-ID header on the request.
func SetRootID(r *http.Request, root, parent ID) {
	r.Header.Set(RootID, fmt.Sprintf("%s:%s", root, parent))
}

// GetRootID returns the root and parent IDs for the request.
func GetRootID(r *http.Request) (root, parent ID, err error) {
	s := r.Header.Get(RootID)
	if s != "" {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) == 2 {
			root, err = ParseID(parts[0])
			if err != nil {
				return
			}

			parent, err = ParseID(parts[1])
			if err != nil {
				return
			}
		}
	}
	return
}

var (
	// RedactedHeaders is a slice of header names whose values should be
	// entirely redacted from logs.
	RedactedHeaders = []string{"Authorization"}
)

// HTTPRequest returns an event which records various aspects of an HTTP request.
// The returned value is incomplete, and should have the response status, size,
// and the elapsed time set before being logged.
func HTTPRequest(r *http.Request) *HTTPRequestEvent {
	return &HTTPRequestEvent{
		Method:        r.Method,
		URI:           r.RequestURI,
		Proto:         r.Proto,
		Headers:       redactHeaders(r),
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		ContentLength: r.ContentLength,
	}
}

// BUG(1.3): HTTPRequestEvent does not record whether the request was over TLS.
// BUG(1.3): HTTPRequestEvent does not record the identity of the TLS peer.

// HTTPRequestEvent records
type HTTPRequestEvent struct {
	Method        string      `json:"method"`
	URI           string      `json:"uri"`
	Proto         string      `json:"proto"`
	Headers       http.Header `json:"headers"`
	Host          string      `json:"host"`
	RemoteAddr    string      `json:"remote_addr"`
	ContentLength int64       `json:"content_length"`
	Status        int         `json:"status"`
	ElapsedMillis float64     `json:"elapsed_ms"`
}

// Schema returns the constant "httprequest".
func (HTTPRequestEvent) Schema() string {
	return "httprequest"
}

var (
	redacted = []string{"REDACTED"}
)

func redactHeaders(r *http.Request) http.Header {
	h := make(http.Header, len(r.Header)+len(r.Trailer))
	for k, v := range r.Header {
		if isRedacted(k) {
			h[k] = redacted
		} else {
			h[k] = v
		}
	}
	for k, v := range r.Trailer {
		if isRedacted(k) {
			h[k] = redacted
		} else {
			h[k] = append(h[k], v...)
		}
	}
	return h
}

func isRedacted(name string) bool {
	for _, v := range RedactedHeaders {
		if strings.EqualFold(name, v) {
			return true
		}
	}
	return false
}
