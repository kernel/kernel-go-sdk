// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/internal/testutil"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/shared"
)

func TestBrowserPoolNewWithOptionalParams(t *testing.T) {
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
	_, err := client.BrowserPools.New(context.TODO(), kernel.BrowserPoolNewParams{
		Size: 10,
		ChromePolicy: map[string]any{
			"foo": "bar",
		},
		Extensions: []shared.BrowserExtensionParam{{
			ID:   kernel.String("id"),
			Name: kernel.String("name"),
		}},
		FillRatePerMinute: kernel.Int(0),
		Headless:          kernel.Bool(false),
		KioskMode:         kernel.Bool(true),
		Name:              kernel.String("my-pool"),
		Network: kernel.BrowserNetworkConfigParam{
			PrivateHosts: []string{"*.example.ts.net", "100.64.0.0/10"},
		},
		Profile: kernel.BrowserPoolNewParamsProfile{
			ID:   kernel.String("id"),
			Name: kernel.String("name"),
		},
		ProxyID:                kernel.String("proxy_id"),
		RefreshOnProfileUpdate: kernel.Bool(true),
		Region:                 kernel.BrowserPoolNewParamsRegionUsEast,
		StartURL:               kernel.String("https://example.com"),
		Stealth:                kernel.Bool(true),
		Telemetry: kernel.BrowserPoolNewParamsTelemetry{
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
				Control: kernel.BrowserTelemetryCategoryConfigParam{
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
				Screenshot: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
				System: kernel.BrowserTelemetryCategoryConfigParam{
					Enabled: kernel.Bool(true),
				},
			},
			Enabled: kernel.Bool(true),
			Export: kernel.BrowserPoolNewParamsTelemetryExport{
				Otlp: kernel.BrowserPoolNewParamsTelemetryExportOtlp{
					Destination: kernel.BrowserPoolNewParamsTelemetryExportOtlpDestination{
						ID:   kernel.String("id"),
						Name: kernel.String("name"),
					},
					Enabled: kernel.Bool(true),
				},
			},
		},
		TimeoutSeconds: kernel.Int(10),
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

func TestBrowserPoolGet(t *testing.T) {
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
	_, err := client.BrowserPools.Get(context.TODO(), "id_or_name")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserPoolUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.BrowserPools.Update(
		context.TODO(),
		"id_or_name",
		kernel.BrowserPoolUpdateParams{
			ChromePolicy: map[string]any{
				"foo": "bar",
			},
			DiscardAllIdle: kernel.Bool(false),
			Extensions: []shared.BrowserExtensionParam{{
				ID:   kernel.String("id"),
				Name: kernel.String("name"),
			}},
			FillRatePerMinute: kernel.Int(0),
			Headless:          kernel.Bool(false),
			KioskMode:         kernel.Bool(true),
			Name:              kernel.String("my-pool"),
			Network: kernel.BrowserNetworkConfigParam{
				PrivateHosts: []string{"*.example.ts.net", "100.64.0.0/10"},
			},
			Profile: kernel.BrowserPoolUpdateParamsProfile{
				ID:   kernel.String("id"),
				Name: kernel.String("name"),
			},
			ProxyID:                kernel.String("proxy_id"),
			RefreshOnProfileUpdate: kernel.Bool(true),
			Size:                   kernel.Int(10),
			StartURL:               kernel.String("https://example.com"),
			Stealth:                kernel.Bool(true),
			Telemetry: kernel.BrowserPoolUpdateParamsTelemetry{
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
					Control: kernel.BrowserTelemetryCategoryConfigParam{
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
					Screenshot: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					System: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
				},
				Enabled: kernel.Bool(true),
				Export: kernel.BrowserPoolUpdateParamsTelemetryExport{
					Otlp: kernel.BrowserPoolUpdateParamsTelemetryExportOtlp{
						Destination: kernel.BrowserPoolUpdateParamsTelemetryExportOtlpDestination{
							ID:   kernel.String("id"),
							Name: kernel.String("name"),
						},
						Enabled: kernel.Bool(true),
					},
				},
			},
			TimeoutSeconds: kernel.Int(10),
			Viewport: shared.BrowserViewportParam{
				Height:      800,
				Width:       1280,
				RefreshRate: kernel.Int(60),
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

func TestBrowserPoolListWithOptionalParams(t *testing.T) {
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
	_, err := client.BrowserPools.List(context.TODO(), kernel.BrowserPoolListParams{
		Limit:  kernel.Int(1),
		Name:   kernel.String("name"),
		Offset: kernel.Int(0),
		Query:  kernel.String("query"),
		Region: kernel.BrowserPoolListParamsRegionUsEast,
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserPoolDeleteWithOptionalParams(t *testing.T) {
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
	err := client.BrowserPools.Delete(
		context.TODO(),
		"id_or_name",
		kernel.BrowserPoolDeleteParams{
			Force: kernel.Bool(true),
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

func TestBrowserPoolAcquireWithOptionalParams(t *testing.T) {
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
	_, err := client.BrowserPools.Acquire(
		context.TODO(),
		"id_or_name",
		kernel.BrowserPoolAcquireParams{
			AcquireTimeoutSeconds: kernel.Int(0),
			Name:                  kernel.String("checkout-flow-1"),
			StartURL:              kernel.String("https://example.com"),
			Tags: kernel.Tags{
				"team": "backend",
				"env":  "staging",
			},
			Telemetry: kernel.BrowserPoolAcquireParamsTelemetry{
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
					Control: kernel.BrowserTelemetryCategoryConfigParam{
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
					Screenshot: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
					System: kernel.BrowserTelemetryCategoryConfigParam{
						Enabled: kernel.Bool(true),
					},
				},
				Enabled: kernel.Bool(true),
				Export: kernel.BrowserPoolAcquireParamsTelemetryExport{
					Otlp: kernel.BrowserPoolAcquireParamsTelemetryExportOtlp{
						Destination: kernel.BrowserPoolAcquireParamsTelemetryExportOtlpDestination{
							ID:   kernel.String("id"),
							Name: kernel.String("name"),
						},
						Enabled: kernel.Bool(true),
					},
				},
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

func TestBrowserPoolFlush(t *testing.T) {
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
	err := client.BrowserPools.Flush(context.TODO(), "id_or_name")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBrowserPoolReleaseWithOptionalParams(t *testing.T) {
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
	err := client.BrowserPools.Release(
		context.TODO(),
		"id_or_name",
		kernel.BrowserPoolReleaseParams{
			SessionID: "ts8iy3sg25ibheguyni2lg9t",
			Reuse:     kernel.Bool(false),
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
