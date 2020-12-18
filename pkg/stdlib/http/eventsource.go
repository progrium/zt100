package http

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type EventSink func([]byte) error

type EventSource struct {
	sinks   sync.Map
	counter uint64
}

func (c *EventSource) Broadcast(d []byte) {
	c.sinks.Range(func(k, v interface{}) bool {
		s := v.(EventSink)
		if err := s(d); err != nil {
			c.Remove(k.(uint64))
		}
		return true
	})
}

func (c *EventSource) Add(sink EventSink) uint64 {
	c.counter++
	c.sinks.Store(c.counter, sink)
	return c.counter
}

func (c *EventSource) Remove(id uint64) {
	c.sinks.Delete(id)
}

func (c *EventSource) IsEventStream(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/event-stream"
}

func (c *EventSource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sink := EventSink(func(d []byte) error {
		if w == nil {
			return nil
		}
		defer flusher.Flush()
		_, err := io.WriteString(w, fmt.Sprintf("data: %s\n\n", string(d)))
		return err
	})

	id := c.Add(sink)
	<-r.Context().Done()
	c.Remove(id)
}
