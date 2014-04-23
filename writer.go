package lunk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

// An EventWriter writes events to an underlying stream in a particular format.
type EventWriter interface {
	Write(e *Event) error
}

// A TextEventWriter writes events to an io.Writer as human-readable lines of
// text.
type TextEventWriter struct {
	w io.Writer
}

// NewTextEventWriter returns a TextEventWriter for the given io.Writer.
func NewTextEventWriter(w io.Writer) TextEventWriter {
	return TextEventWriter{
		w: w,
	}
}

func (f TextEventWriter) Write(e *Event) error {
	_, err := fmt.Fprintln(f.w, e.String())
	return err
}

// A JSONEventWriter writes events to an io.Writer as a stream of JSON
// objects.
type JSONEventWriter struct {
	w *bufio.Writer
	e *json.Encoder
}

// NewJSONEventWriter returns a JSONEventWriter for the given io.Writer.
func NewJSONEventWriter(w io.Writer) JSONEventWriter {
	return JSONEventWriter{
		e: json.NewEncoder(w),
	}
}

func (w JSONEventWriter) Write(e *Event) error {
	return w.e.Encode(e)
}

// AsyncErrorHandler is a callback function for errors which occur during
// asynchronous writing of events.
type AsyncErrorHandler func(*Event, error)

// AsyncEventWriter handles the writing of events in a separate goroutine.
type AsyncEventWriter struct {
	c  chan *Event
	w  EventWriter
	wg *sync.WaitGroup
	h  AsyncErrorHandler
}

// NewAsyncEventWriter creates a new AsyncEventWriter given an underlying
// EventWriter, a buffer size, and an error handler.
func NewAsyncEventWriter(w EventWriter, bufSize int, h AsyncErrorHandler) *AsyncEventWriter {
	return &AsyncEventWriter{
		c:  make(chan *Event, bufSize),
		w:  w,
		wg: new(sync.WaitGroup),
		h:  h,
	}
}

func (w *AsyncEventWriter) Write(e *Event) error {
	w.c <- e
	return nil
}

// Start creates a writer goroutine. This must be called for events to be written
// to the underlying EventWriter.
func (w *AsyncEventWriter) Start() {
	w.wg.Add(1)
	go func() {
		for e := range w.c {
			if err := w.w.Write(e); err != nil {
				w.h(e, err)
			}
		}
		w.wg.Done()
	}()
}

// Stop prevents the writer from accepting any further events and blocks until
// alls events have been processed.
func (w *AsyncEventWriter) Stop() {
	close(w.c)
	w.wg.Wait()
}
