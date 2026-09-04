// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/internal/testutil"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/shared"
)

func TestBrowserNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Browsers.New(context.TODO(), kernel.BrowserNewParams{
		ChromePolicy: map[string]any{
			"foo": "bar",
		},
		Extensions: []shared.BrowserExtensionParam{{
			ID:   kernel.String("id"),
			Name: kernel.String("name"),
		}},
		GPU:          kernel.Bool(false),
		Headless:     kernel.Bool(false),
		InvocationID: kernel.String("rr33xuugxj9h0bkf1rdt2bet"),
		KioskMode:    kernel.Bool(true),
		Memory:       kernel.BrowserMemoryRequest8GiB,
		Name:         kernel.String("checkout-flow-1"),
		Network: kernel.BrowserNetworkConfigParam{
			PrivateHosts: []string{"*.example.ts.net", "100.64.0.0/10"},
		},
		Profile: shared.BrowserProfileParam{
			ID:          kernel.String("id"),
			Name:        kernel.String("name"),
			SaveChanges: kernel.Bool(true),
		},
		Proxy: kernel.BrowserProxyConfigParam{
			ID:   kernel.String("x"),
			Mode: kernel.BrowserProxyModeDirect,
			Name: kernel.String("x"),
		},
		ProxyID:  kernel.String("proxy_id"),
		Region:   kernel.BrowserNewParamsRegionUsEast,
		StartURL: kernel.String("https://example.com"),
		Stealth:  kernel.Bool(true),
		Tags: kernel.Tags{
			"team": "backend",
			"env":  "staging",
		},
		Telemetry: kernel.BrowserNewParamsTelemetry{
			Browser: kernel.BrowserTelemetryCategoriesConfigParam{
				Captcha: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Connection: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Console: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Control: kernel.BrowserTelemetryControlConfigParam{
					Cdp: kernel.BrowserTelemetryCdpControlConfigParam{
						ExcludedMethods: []kernel.BrowserCdpCommandMethod{kernel.BrowserCdpCommandMethodInputDispatchMouseEvent},
					},
					Enabled: kernel.Bool(true),
				},
				Interaction: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Network: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Page: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Platform: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				Screenshot: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				System: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
			},
			Enabled: kernel.Bool(true),
			Export: kernel.BrowserNewParamsTelemetryExport{
				Otlp: kernel.BrowserNewParamsTelemetryExportOtlp{
					Destination: kernel.BrowserNewParamsTelemetryExportOtlpDestination{
						ID:   kernel.String("id"),
						Name: kernel.String("name"),
					},
					Enabled: kernel.Bool(true),
				},
			},
		},
		TimeoutSeconds: kernel.Int(10),
		Vaults: []kernel.VaultReferenceParam{{
			ID:   kernel.String("id"),
			Name: kernel.String("x"),
		}},
		Viewport: shared.BrowserViewportParam{
			Height:      800,
			Width:       1280,
			RefreshRate: kernel.Int(60),
		},
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Browsers.Get(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserGetParams{
			IncludeDeleted: kernel.Bool(true),
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

func TestBrowserUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Browsers.Update(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserUpdateParams{
			DisableDefaultProxy: kernel.Bool(true),
			Name:                kernel.String("checkout-flow-1"),
			Profile: shared.BrowserProfileParam{
				ID:          kernel.String("id"),
				Name:        kernel.String("name"),
				SaveChanges: kernel.Bool(true),
			},
			Proxy: kernel.BrowserProxyConfigParam{
				ID:   kernel.String("x"),
				Mode: kernel.BrowserProxyModeDirect,
				Name: kernel.String("x"),
			},
			ProxyID: kernel.String("proxy_id"),
			Tags: kernel.Tags{
				"team": "backend",
				"env":  "staging",
			},
			Telemetry: kernel.BrowserUpdateParamsTelemetry{
				Browser: kernel.BrowserTelemetryCategoriesConfigParam{
					Captcha: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Connection: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Console: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Control: kernel.BrowserTelemetryControlConfigParam{
						Cdp: kernel.BrowserTelemetryCdpControlConfigParam{
							ExcludedMethods: []kernel.BrowserCdpCommandMethod{kernel.BrowserCdpCommandMethodInputDispatchMouseEvent},
						},
						Enabled: kernel.Bool(true),
					},
					Interaction: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Network: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Page: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Platform: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					Screenshot: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					System: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
				},
				Enabled: kernel.Bool(true),
				Export: kernel.BrowserUpdateParamsTelemetryExport{
					Otlp: kernel.BrowserUpdateParamsTelemetryExportOtlp{
						Destination: kernel.BrowserUpdateParamsTelemetryExportOtlpDestination{
							ID:   kernel.String("id"),
							Name: kernel.String("name"),
						},
						Enabled: kernel.Bool(true),
					},
				},
			},
			Viewport: kernel.BrowserUpdateParamsViewport{
				BrowserViewportParam: shared.BrowserViewportParam{
					Height:      800,
					Width:       1280,
					RefreshRate: kernel.Int(60),
				},
				Force: kernel.Bool(true),
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
}

func TestBrowserListWithOptionalParams(t *testing.T) {
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
	_, err := client.Browsers.List(context.TODO(), kernel.BrowserListParams{
		IncludeDeleted: kernel.Bool(true),
		Limit:          kernel.Int(1),
		Offset:         kernel.Int(0),
		Query:          kernel.String("query"),
		Region:         kernel.BrowserListParamsRegionUsEast,
		Status:         kernel.BrowserListParamsStatusActive,
		Tags: map[string]string{
			"foo": "string",
		},
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserCurlWithOptionalParams(t *testing.T) {
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
	_, err := client.Browsers.Curl(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserCurlParams{
			URL:  "url",
			Body: kernel.String("body"),
			Headers: map[string]string{
				"foo": "string",
			},
			Method:           kernel.BrowserCurlParamsMethodGet,
			ResponseEncoding: kernel.BrowserCurlParamsResponseEncodingUtf8,
			TimeoutMs:        kernel.Int(1000),
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

func TestBrowserDeleteByID(t *testing.T) {
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
	err := client.Browsers.DeleteByID(context.TODO(), "htzv5orfit78e1m2biiifpbv")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserLoadExtensions(t *testing.T) {
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
	err := client.Browsers.LoadExtensions(
		context.TODO(),
		"htzv5orfit78e1m2biiifpbv",
		kernel.BrowserLoadExtensionsParams{
			Extensions: []kernel.BrowserLoadExtensionsParamsExtension{{
				Name:    "name",
				ZipFile: io.Reader(bytes.NewBuffer([]byte("Example data"))),
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
