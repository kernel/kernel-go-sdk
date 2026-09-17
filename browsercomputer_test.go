// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/internal/testutil"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestBrowserComputerBatch(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.Batch(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerBatchParams{
			Actions: []kernel.BrowserComputerBatchParamsAction{{
				Type: "click_mouse",
				ClickMouse: kernel.BrowserComputerBatchParamsActionClickMouse{
					X:         0,
					Y:         0,
					Button:    "left",
					ClickType: "down",
					HoldKeys:  []string{"string"},
					NumClicks: kernel.Int(0),
				},
				DragMouse: kernel.BrowserComputerBatchParamsActionDragMouse{
					Path:            [][]int64{{0, 0}, {0, 0}},
					Button:          "left",
					Delay:           kernel.Int(0),
					DurationMs:      kernel.Int(50),
					HoldKeys:        []string{"string"},
					Smooth:          kernel.Bool(true),
					StepDelayMs:     kernel.Int(0),
					StepsPerSegment: kernel.Int(1),
				},
				MoveMouse: kernel.BrowserComputerBatchParamsActionMoveMouse{
					X:          0,
					Y:          0,
					DurationMs: kernel.Int(50),
					HoldKeys:   []string{"string"},
					Smooth:     kernel.Bool(true),
				},
				PressKey: kernel.BrowserComputerBatchParamsActionPressKey{
					Keys:     []string{"string"},
					Duration: kernel.Int(0),
					HoldKeys: []string{"string"},
				},
				Scroll: kernel.BrowserComputerBatchParamsActionScroll{
					X:        0,
					Y:        0,
					DeltaX:   kernel.Int(0),
					DeltaY:   kernel.Int(0),
					HoldKeys: []string{"string"},
				},
				SetCursor: kernel.BrowserComputerBatchParamsActionSetCursor{
					Hidden: true,
				},
				Sleep: kernel.BrowserComputerBatchParamsActionSleep{
					DurationMs: 0,
				},
				TypeText: kernel.BrowserComputerBatchParamsActionTypeText{
					Text:  "text",
					Delay: kernel.Int(0),
				},
			}},
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerCaptureScreenshotWithOptionalParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("abc"))
	}))
	defer server.Close()
	baseURL := server.URL
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	resp, err := client.Browsers.Computer.CaptureScreenshot(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerCaptureScreenshotParams{
			Region: kernel.BrowserComputerCaptureScreenshotParamsRegion{
				Height: 0,
				Width:  0,
				X:      0,
				Y:      0,
			},
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if !bytes.Equal(b, []byte("abc")) {
		t.Fatalf("return value not %s: %s", "abc", b)
	}
}

func TestBrowserComputerClickMouseWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.ClickMouse(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerClickMouseParams{
			X:         0,
			Y:         0,
			Button:    kernel.BrowserComputerClickMouseParamsButtonLeft,
			ClickType: kernel.BrowserComputerClickMouseParamsClickTypeDown,
			HoldKeys:  []string{"string"},
			NumClicks: kernel.Int(0),
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerDragMouseWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.DragMouse(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerDragMouseParams{
			Path:            [][]int64{{0, 0}, {0, 0}},
			Button:          kernel.BrowserComputerDragMouseParamsButtonLeft,
			Delay:           kernel.Int(0),
			DurationMs:      kernel.Int(50),
			HoldKeys:        []string{"string"},
			Smooth:          kernel.Bool(true),
			StepDelayMs:     kernel.Int(0),
			StepsPerSegment: kernel.Int(1),
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerGetMousePosition(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Browsers.Computer.GetMousePosition(context.TODO(), "htzv5orfit78e1m2biiifpbv")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerMoveMouseWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.MoveMouse(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerMoveMouseParams{
			X:          0,
			Y:          0,
			DurationMs: kernel.Int(50),
			HoldKeys:   []string{"string"},
			Smooth:     kernel.Bool(true),
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerPressKeyWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.PressKey(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerPressKeyParams{
			Keys:     []string{"string"},
			Duration: kernel.Int(0),
			HoldKeys: []string{"string"},
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerReadClipboard(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Browsers.Computer.ReadClipboard(context.TODO(), "htzv5orfit78e1m2biiifpbv")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerScrollWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.Scroll(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerScrollParams{
			X:        0,
			Y:        0,
			DeltaX:   kernel.Int(0),
			DeltaY:   kernel.Int(0),
			HoldKeys: []string{"string"},
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerSetCursorVisibility(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Browsers.Computer.SetCursorVisibility(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerSetCursorVisibilityParams{
			Hidden: true,
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerTypeTextWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.TypeText(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerTypeTextParams{
			Text:  "text",
			Delay: kernel.Int(0),
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserComputerWriteClipboard(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.Browsers.Computer.WriteClipboard(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserComputerWriteClipboardParams{
			Text: "text",
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
