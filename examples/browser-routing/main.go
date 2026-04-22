package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	kernel "github.com/kernel/kernel-go-sdk"
)

func main() {
	ctx := context.Background()
	client := kernel.NewClient(
		kernel.WithBrowserRouting(kernel.BrowserRoutingConfig{
			Enabled:      true,
			Subresources: []string{"process"},
		}),
	)

	browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
		Headless: kernel.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Browsers.DeleteByID(context.Background(), browser.SessionID)
	}()

	httpClient, err := client.Browsers.HTTPClient(browser.SessionID)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Fprintln(os.Stdout, "status", resp.StatusCode)
}
