// Command browser-telemetry-routing-smoke is a live smoke test that proves the
// browser telemetry SSE stream is routed directly to the browser VM by default.
//
// It creates a telemetry-enabled browser, opens the telemetry stream, generates
// activity via browsers.curl, and asserts that (a) at least one telemetry event
// arrives and (b) the telemetry stream request was rewritten to the VM base URL
// (host contains "proxy." and ":8443", path ends with "/telemetry/stream",
// Authorization stripped, jwt query param present).
//
// Run with KERNEL_API_KEY set:
//
//	go run ./examples/browser-telemetry-routing-smoke
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

type recorder struct {
	mu             sync.Mutex
	telemetryURL   string
	telemetryAuth  string
	telemetryFound bool
}

func (rec *recorder) middleware(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	// Call next first so the direct-VM routing middleware (registered after this
	// one) has rewritten the in-place request URL before we record it.
	res, err := next(req)
	if strings.Contains(req.URL.Path, "/telemetry/stream") {
		rec.mu.Lock()
		rec.telemetryURL = req.URL.String()
		rec.telemetryAuth = req.Header.Get("Authorization")
		rec.telemetryFound = true
		rec.mu.Unlock()
	}
	return res, err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "SMOKE FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("SMOKE PASS")
}

func run() error {
	// Belt-and-suspenders: telemetry is a default routing subresource, but set it
	// explicitly so the test is robust to env defaults.
	_ = os.Setenv("KERNEL_BROWSER_ROUTING_SUBRESOURCES", "curl,telemetry")

	rec := &recorder{}
	client := kernel.NewClient(option.WithMiddleware(rec.middleware))

	ctx := context.Background()

	browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
		Headless:       kernel.Bool(true),
		TimeoutSeconds: kernel.Int(120),
		Telemetry: kernel.BrowserNewParamsTelemetry{
			Enabled: kernel.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("create browser: %w", err)
	}
	sessionID := browser.SessionID
	fmt.Println("created browser", sessionID, "base_url", browser.BaseURL)

	defer func() {
		if err := client.Browsers.DeleteByID(context.Background(), sessionID); err != nil {
			fmt.Fprintln(os.Stderr, "warning: delete browser:", err)
		} else {
			fmt.Println("deleted browser", sessionID)
		}
	}()

	// Open the telemetry stream with a bounded lifetime.
	streamCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	// Generate VM API activity throughout. The "api" telemetry category emits an
	// event for every VM API call (e.g. browsers.curl), so this produces telemetry
	// events within ~1s.
	go func() {
		for streamCtx.Err() == nil {
			_, cerr := client.Browsers.Curl(ctx, sessionID, kernel.BrowserCurlParams{
				URL:              "https://example.com/",
				Method:           kernel.BrowserCurlParamsMethodGet,
				ResponseEncoding: kernel.BrowserCurlParamsResponseEncodingUtf8,
				TimeoutMs:        kernel.Int(10000),
			})
			if cerr != nil {
				fmt.Fprintln(os.Stderr, "warning: curl:", cerr)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Read until we get at least one event. The SSE stream interleaves keepalive
	// comment frames (": ...") which the generated Stream wrapper surfaces as an
	// "unexpected end of JSON input" error and terminates; when that happens we
	// reopen the stream and keep reading until the deadline. A data frame from VM
	// API activity unmarshals cleanly into an event.
	var got kernel.BrowserTelemetryStreamResponse
	gotEvent := false
	var lastErr error
	for streamCtx.Err() == nil && !gotEvent {
		stream := client.Browsers.Telemetry.StreamStreaming(streamCtx, sessionID, kernel.BrowserTelemetryStreamParams{})
		for stream.Next() {
			got = stream.Current()
			gotEvent = true
			break
		}
		if !gotEvent {
			lastErr = stream.Err()
			_ = stream.Close()
		}
	}

	if !gotEvent {
		if lastErr != nil {
			return fmt.Errorf("no telemetry events received; last stream err: %w", lastErr)
		}
		return fmt.Errorf("no telemetry events received within deadline")
	}
	fmt.Printf("received telemetry event: seq=%d type=%v\n", got.Seq, got.Event.Type)

	// Verify routing.
	rec.mu.Lock()
	url, auth, found := rec.telemetryURL, rec.telemetryAuth, rec.telemetryFound
	rec.mu.Unlock()

	if !found {
		return fmt.Errorf("no telemetry stream request was recorded")
	}
	fmt.Println("recorded telemetry stream URL:", url)

	if strings.Contains(url, "api.onkernel.com") {
		return fmt.Errorf("telemetry request was NOT routed to VM (hit api.onkernel.com): %s", url)
	}
	// Strip query for path/host checks.
	noQuery := url
	if i := strings.Index(noQuery, "?"); i >= 0 {
		noQuery = noQuery[:i]
	}
	if !strings.HasSuffix(noQuery, "/telemetry/stream") {
		return fmt.Errorf("telemetry path does not end with /telemetry/stream: %s", url)
	}
	if !strings.Contains(url, "proxy.") || !strings.Contains(url, ":8443") {
		return fmt.Errorf("telemetry URL does not look like a VM base_url (expected proxy. and :8443): %s", url)
	}
	if auth != "" {
		return fmt.Errorf("Authorization header was NOT stripped: %q", auth)
	}
	if !strings.Contains(url, "jwt=") {
		return fmt.Errorf("jwt query param missing from telemetry URL: %s", url)
	}

	// Sanity check the recorded host matches the session base_url host.
	if base := browser.BaseURL; base != "" {
		bh := base
		bh = strings.TrimPrefix(strings.TrimPrefix(bh, "https://"), "http://")
		if i := strings.IndexAny(bh, "/?"); i >= 0 {
			bh = bh[:i]
		}
		if !strings.Contains(url, bh) {
			return fmt.Errorf("telemetry URL host %q does not match session base_url host %q", url, bh)
		}
	}

	fmt.Println("routing verified: telemetry stream went directly to the VM")
	return nil
}
