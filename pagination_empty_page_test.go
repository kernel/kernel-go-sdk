package kernel_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

type emptyPageResponse struct {
	body    string
	headers map[string]string
}

// countEventsAcrossPages auto-pages the telemetry events endpoint against a
// server that serves pages keyed by the request's offset query value.
func countEventsAcrossPages(t *testing.T, pages map[string]emptyPageResponse) (events, requests int, err error) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		page, ok := pages[r.URL.Query().Get("offset")]
		if !ok {
			t.Errorf("unexpected request for offset %q", r.URL.Query().Get("offset"))
			http.Error(w, "unexpected offset", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		for k, v := range page.headers {
			w.Header().Set(k, v)
		}
		fmt.Fprint(w, page.body)
	}))
	defer srv.Close()

	client := kernel.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("test-api-key"), option.WithMaxRetries(0))
	iter := client.Browsers.Telemetry.EventsAutoPaging(context.Background(), "browser", kernel.BrowserTelemetryEventsParams{})
	for iter.Next() {
		events++
	}
	return events, requests, iter.Err()
}

const emptyPageEvent = `{"seq":1,"event":{"ts":1,"type":"page_load","category":"page","source":{"kind":"kernel_api"},"data":{}}}`

func TestOffsetPaginationContinuesPastEmptyPages(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"":   {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "20"}},
		"20": {"[]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "40"}},
		"40": {"[]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "60"}},
		"60": {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "false", "X-Next-Offset": "0"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 2 || requests != 4 {
		t.Fatalf("got %d events in %d requests, want 2 events in 4 requests", events, requests)
	}
}

func TestOffsetPaginationEmptyFirstPageContinues(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"":   {"[]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "20"}},
		"20": {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "false", "X-Next-Offset": "0"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 1 || requests != 2 {
		t.Fatalf("got %d events in %d requests, want 1 event in 2 requests", events, requests)
	}
}

func TestOffsetPaginationStopsOnEmptyLastPage(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"": {"[]", map[string]string{"X-Has-More": "false", "X-Next-Offset": "0"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 0 || requests != 1 {
		t.Fatalf("got %d events in %d requests, want 0 events in 1 request", events, requests)
	}
}

func TestOffsetPaginationStopsOnEmptyPageWithoutCursor(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"": {"[]", nil},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 0 || requests != 1 {
		t.Fatalf("got %d events in %d requests, want 0 events in 1 request", events, requests)
	}
}

func TestOffsetPaginationStopsWhenCursorDoesNotAdvance(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"":   {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "20"}},
		"20": {"[]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "20"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 1 || requests != 2 {
		t.Fatalf("got %d events in %d requests, want 1 event in 2 requests", events, requests)
	}
}

func TestOffsetPaginationFollowsDescendingCursorAcrossEmptyPage(t *testing.T) {
	events, requests, err := countEventsAcrossPages(t, map[string]emptyPageResponse{
		"":    {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "300"}},
		"300": {"[]", map[string]string{"X-Has-More": "true", "X-Next-Offset": "200"}},
		"200": {"[" + emptyPageEvent + "]", map[string]string{"X-Has-More": "false", "X-Next-Offset": "0"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != 2 || requests != 3 {
		t.Fatalf("got %d events in %d requests, want 2 events in 3 requests", events, requests)
	}
}
