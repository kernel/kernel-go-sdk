package main

import (
	"context"
	"fmt"
	"os"

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

	// Telemetry is a default direct-to-VM routing subresource, so the stream goes
	// straight to the browser VM automatically.
	stream := client.Browsers.Telemetry.StreamStreaming(ctx, browser.SessionID, kernel.BrowserTelemetryStreamParams{})
	defer stream.Close()

	// Make a few browser activity calls to generate events. The "api" telemetry
	// category emits an event per VM API call, so events arrive within ~1s.
	for i := 0; i < 3; i++ {
		if _, err := client.Browsers.Curl(ctx, browser.SessionID, kernel.BrowserCurlParams{
			URL: "https://example.com",
		}); err != nil {
			panic(err)
		}
	}

	// Print a few events, then stop.
	for printed := 0; printed < 3 && stream.Next(); printed++ {
		event := stream.Current()
		fmt.Fprintf(os.Stdout, "telemetry event seq=%d type=%s\n", event.Seq, event.Event.Type)
	}
}
