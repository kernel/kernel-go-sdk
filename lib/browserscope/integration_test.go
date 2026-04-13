package browserscope_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestIntegrationBrowserSessionClient(t *testing.T) {
	apiKey := strings.TrimSpace(os.Getenv("KERNEL_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("KERNEL_BASE_URL"))
	if apiKey == "" || baseURL == "" {
		t.Skip("set KERNEL_API_KEY and KERNEL_BASE_URL to run integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := kernel.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	browser, err := client.Browsers.New(ctx, kernel.BrowserNewParams{
		Headless:       kernel.Bool(true),
		TimeoutSeconds: kernel.Int(60),
	})
	if err != nil {
		t.Fatalf("create browser: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Browsers.DeleteByID(context.Background(), browser.SessionID)
	})

	if browser.BaseURL == "" {
		t.Fatal("expected browser base_url to be set")
	}
	if !strings.Contains(browser.BaseURL, "/browser/kernel") {
		t.Fatalf("expected browser base_url to include /browser/kernel, got %q", browser.BaseURL)
	}

	scoped, err := client.ForBrowser(browser)
	if err != nil {
		t.Fatalf("for browser: %v", err)
	}

	execRes, err := scoped.Process.Exec(ctx, kernel.BrowserProcessExecParams{
		Command: "echo",
		Args:    []string{"hello"},
	})
	if err != nil {
		t.Fatalf("process exec: %v", err)
	}
	if execRes.ExitCode != 0 {
		t.Fatalf("expected process exit code 0, got %d", execRes.ExitCode)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := scoped.HTTPClient().Do(req)
	if err != nil {
		t.Fatalf("browser http client: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		t.Fatalf("expected successful browser http response, got %d", resp.StatusCode)
	}
}
