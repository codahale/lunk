package lunk

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

// An EntryRecorder records entries, e.g. to a streaming processing system or an
// OLAP database.
type EntryRecorder interface {
	io.Closer
	Record(Entry) error
}

// NewNormalizedCSVEntryRecorder returns an EntryRecorder which writes events to
// one CSV file and properties to another.
func NewNormalizedCSVEntryRecorder(events, props io.Writer) EntryRecorder {
	return nCSVRecorder{
		events: csv.NewWriter(events),
		props:  csv.NewWriter(props),
	}
}

// NewDenormalizedCSVEntryRecorder returns an EntryRecorder which writes events
// and their properties to a single CSV file, duplicating event data when
// necessary.
func NewDenormalizedCSVEntryRecorder(w io.Writer) EntryRecorder {
	return dCSVRecorder{
		w: csv.NewWriter(w),
	}
}

type nCSVRecorder struct {
	events *csv.Writer
	props  *csv.Writer
}

func (r nCSVRecorder) Record(e Entry) error {
	root, id, parent := e.Root.String(), e.ID.String(), e.Parent.String()

	if err := r.events.Write([]string{
		root,
		id,
		parent,
		e.Schema,
		e.Time.Format(time.RFC3339Nano),
		e.Host,
		strconv.Itoa(e.PID),
		e.Deploy,
	}); err != nil {
		return err
	}

	for k, v := range e.Properties {
		if err := r.props.Write([]string{
			root,
			id,
			parent,
			k,
			v,
		}); err != nil {
			return err
		}

	}

	return nil
}

func (r nCSVRecorder) Close() error {
	r.events.Flush()
	r.props.Flush()
	return nil
}

type dCSVRecorder struct {
	w *csv.Writer
}

func (r dCSVRecorder) Record(e Entry) error {
	root, id, parent := e.Root.String(), e.ID.String(), e.Parent.String()
	time := e.Time.Format(time.RFC3339Nano)
	pid := strconv.Itoa(e.PID)

	for k, v := range e.Properties {
		if err := r.w.Write([]string{
			root,
			id,
			parent,
			e.Schema,
			time,
			e.Host,
			pid,
			e.Deploy,
			k,
			v,
		}); err != nil {
			return err
		}

	}

	return nil
}

func (r dCSVRecorder) Close() error {
	r.w.Flush()
	return nil
}
