package main

import (
	"context"
	"fmt"
	"os"
	"time"

	kernel "github.com/kernel/kernel-go-sdk"
)

func main() {
	ctx := context.Background()
	client := kernel.NewClient()

	// Create a browser with telemetry enabled so the VM emits telemetry events.
	browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
		Headless: kernel.Bool(true),
		Telemetry: kernel.BrowserNewParamsTelemetry{
			Enabled: kernel.Bool(true),
		},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Browsers.DeleteByID(context.Background(), browser.SessionID)
	}()

	// Generate activity in the background. The "api" telemetry category emits an
	// event per VM API call, so these curls produce a steady stream of telemetry
	// events within ~1s.
	activityCtx, stopActivity := context.WithCancel(ctx)
	defer stopActivity()
	go func() {
		for activityCtx.Err() == nil {
			_, _ = client.Browsers.Curl(activityCtx, browser.SessionID, kernel.BrowserCurlParams{
				URL: "https://example.com",
			})
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Telemetry is a default direct-to-VM routing subresource, so the stream goes
	// straight to the browser VM automatically. Read a few events, print them, then
	// stop.
	stream := client.Browsers.Telemetry.StreamStreaming(ctx, browser.SessionID, kernel.BrowserTelemetryStreamParams{})
	defer stream.Close()

	for printed := 0; printed < 3 && stream.Next(); printed++ {
		event := stream.Current()
		fmt.Fprintf(os.Stdout, "telemetry event seq=%d type=%s\n", event.Seq, event.Event.Type)
	}
}
