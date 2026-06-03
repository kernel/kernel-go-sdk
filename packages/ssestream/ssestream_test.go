package ssestream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func newStringResponse(body string) *http.Response {
	return &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   io.NopCloser(strings.NewReader(body)),
	}
}

// The server sends a ":\n\n" keepalive comment frame every ~15s on an idle stream. It must be
// skipped, not surfaced as an empty-data event that fails to JSON-decode and ends the stream.
func TestStreamSkipsKeepaliveCommentFrames(t *testing.T) {
	body := "id: 1\ndata: {\"seq\":1}\n\n" + ":\n\n" + "data: {\"seq\":2}\n\n"

	type payload struct {
		Seq int `json:"seq"`
	}
	stream := NewStream[payload](NewDecoder(newStringResponse(body)), nil)

	var got []int
	for stream.Next() {
		got = append(got, stream.Current().Seq)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("expected events [1 2], got %v", got)
	}
}

// A standalone keepalive comment (no surrounding events) must not yield anything or error.
func TestStreamKeepaliveOnlyYieldsNothing(t *testing.T) {
	stream := NewStream[map[string]any](NewDecoder(newStringResponse(":\n\n")), nil)
	if stream.Next() {
		t.Fatalf("expected no events, got %v", stream.Current())
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
}
