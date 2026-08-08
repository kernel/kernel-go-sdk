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
)

func TestAuthConnectionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.New(context.TODO(), kernel.AuthConnectionNewParams{
		ManagedAuthCreateRequest: kernel.ManagedAuthCreateRequestParam{
			Domain:         "netflix.com",
			ProfileName:    "user-123",
			AllowedDomains: []string{"login.netflix.com", "auth.netflix.com"},
			AutoReauth:     kernel.Bool(true),
			Browser: kernel.ManagedAuthBrowserConfigParam{
				Proxy: kernel.BrowserProxyConfigParam{
					ID:   kernel.String("x"),
					Mode: kernel.BrowserProxyModeDirect,
					Name: kernel.String("x"),
				},
				Stealth: kernel.Bool(false),
				Telemetry: kernel.ManagedAuthBrowserConfigTelemetryParam{
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
					Export: kernel.ManagedAuthBrowserConfigTelemetryExportParam{
						Otlp: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpParam{
							Destination: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam{
								ID:   kernel.String("id"),
								Name: kernel.String("name"),
							},
							Enabled: kernel.Bool(true),
						},
					},
				},
			},
			BrowserTelemetry: kernel.ManagedAuthCreateRequestBrowserTelemetryParam{
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
				Export: kernel.ManagedAuthCreateRequestBrowserTelemetryExportParam{
					Otlp: kernel.ManagedAuthCreateRequestBrowserTelemetryExportOtlpParam{
						Destination: kernel.ManagedAuthCreateRequestBrowserTelemetryExportOtlpDestinationParam{
							ID:   kernel.String("id"),
							Name: kernel.String("name"),
						},
						Enabled: kernel.Bool(true),
					},
				},
			},
			Credential: kernel.ManagedAuthCreateRequestCredentialParam{
				Auto:     kernel.Bool(true),
				Name:     kernel.String("my-netflix-creds"),
				Path:     kernel.String("Personal/Netflix"),
				Provider: kernel.String("my-1p"),
			},
			HealthCheckInterval: kernel.Int(3600),
			HealthChecks:        kernel.Bool(true),
			LoginURL:            kernel.String("https://netflix.com/login"),
			Proxy: kernel.ManagedAuthCreateRequestProxyParam{
				ID:   kernel.String("id"),
				Name: kernel.String("name"),
			},
			RecordSession:   kernel.Bool(false),
			SaveCredentials: kernel.Bool(true),
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

func TestAuthConnectionGet(t *testing.T) {
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
	_, err := client.Auth.Connections.Get(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAuthConnectionUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.Update(
		context.TODO(),
		"id",
		kernel.AuthConnectionUpdateParams{
			ManagedAuthUpdateRequest: kernel.ManagedAuthUpdateRequestParam{
				AllowedDomains: []string{"login.netflix.com", "auth.netflix.com"},
				AutoReauth:     kernel.Bool(true),
				Browser: kernel.ManagedAuthBrowserConfigParam{
					Proxy: kernel.BrowserProxyConfigParam{
						ID:   kernel.String("x"),
						Mode: kernel.BrowserProxyModeDirect,
						Name: kernel.String("x"),
					},
					Stealth: kernel.Bool(false),
					Telemetry: kernel.ManagedAuthBrowserConfigTelemetryParam{
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
						Export: kernel.ManagedAuthBrowserConfigTelemetryExportParam{
							Otlp: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpParam{
								Destination: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam{
									ID:   kernel.String("id"),
									Name: kernel.String("name"),
								},
								Enabled: kernel.Bool(true),
							},
						},
					},
				},
				BrowserTelemetry: kernel.ManagedAuthUpdateRequestBrowserTelemetryParam{
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
					Export: kernel.ManagedAuthUpdateRequestBrowserTelemetryExportParam{
						Otlp: kernel.ManagedAuthUpdateRequestBrowserTelemetryExportOtlpParam{
							Destination: kernel.ManagedAuthUpdateRequestBrowserTelemetryExportOtlpDestinationParam{
								ID:   kernel.String("id"),
								Name: kernel.String("name"),
							},
							Enabled: kernel.Bool(true),
						},
					},
				},
				Credential: kernel.ManagedAuthUpdateRequestCredentialParam{
					Auto:     kernel.Bool(true),
					Name:     kernel.String("my-netflix-creds"),
					Path:     kernel.String("Personal/Netflix"),
					Provider: kernel.String("my-1p"),
				},
				HealthCheckInterval: kernel.Int(3600),
				HealthChecks:        kernel.Bool(true),
				LoginURL:            kernel.String("https://netflix.com/login"),
				Proxy: kernel.ManagedAuthUpdateRequestProxyParam{
					ID:   kernel.String("id"),
					Name: kernel.String("name"),
				},
				RecordSession:   kernel.Bool(false),
				SaveCredentials: kernel.Bool(true),
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

func TestAuthConnectionListWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.List(context.TODO(), kernel.AuthConnectionListParams{
		Domain:      kernel.String("domain"),
		Limit:       kernel.Int(100),
		Offset:      kernel.Int(0),
		ProfileName: kernel.String("profile_name"),
		Query:       kernel.String("query"),
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAuthConnectionDelete(t *testing.T) {
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
	err := client.Auth.Connections.Delete(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAuthConnectionLoginWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.Login(
		context.TODO(),
		"id",
		kernel.AuthConnectionLoginParams{
			Browser: kernel.ManagedAuthBrowserConfigParam{
				Proxy: kernel.BrowserProxyConfigParam{
					ID:   kernel.String("x"),
					Mode: kernel.BrowserProxyModeDirect,
					Name: kernel.String("x"),
				},
				Stealth: kernel.Bool(false),
				Telemetry: kernel.ManagedAuthBrowserConfigTelemetryParam{
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
					Export: kernel.ManagedAuthBrowserConfigTelemetryExportParam{
						Otlp: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpParam{
							Destination: kernel.ManagedAuthBrowserConfigTelemetryExportOtlpDestinationParam{
								ID:   kernel.String("id"),
								Name: kernel.String("name"),
							},
							Enabled: kernel.Bool(true),
						},
					},
				},
			},
			BrowserTelemetry: kernel.AuthConnectionLoginParamsBrowserTelemetry{
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
				Export: kernel.AuthConnectionLoginParamsBrowserTelemetryExport{
					Otlp: kernel.AuthConnectionLoginParamsBrowserTelemetryExportOtlp{
						Destination: kernel.AuthConnectionLoginParamsBrowserTelemetryExportOtlpDestination{
							ID:   kernel.String("id"),
							Name: kernel.String("name"),
						},
						Enabled: kernel.Bool(true),
					},
				},
			},
			Proxy: kernel.AuthConnectionLoginParamsProxy{
				ID:   kernel.String("id"),
				Name: kernel.String("name"),
			},
			RecordSession: kernel.Bool(true),
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

func TestAuthConnectionSubmitWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.Submit(
		context.TODO(),
		"id",
		kernel.AuthConnectionSubmitParams{
			SubmitFieldsRequest: kernel.SubmitFieldsRequestParam{
				FieldValues: map[string]string{
					"field_email":    "user@example.com",
					"field_password": "secret",
				},
				Fields: map[string]string{
					"email":    "user@example.com",
					"password": "secret",
				},
				MfaOptionID:       kernel.String("sms"),
				SelectedChoiceID:  kernel.String("google"),
				SignInOptionID:    kernel.String("work-account"),
				SSOButtonSelector: kernel.String("xpath=//button[contains(text(), 'Continue with Google')]"),
				SSOProvider:       kernel.String("google"),
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

func TestAuthConnectionTimelineWithOptionalParams(t *testing.T) {
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
	_, err := client.Auth.Connections.Timeline(
		context.TODO(),
		"id",
		kernel.AuthConnectionTimelineParams{
			Limit:  kernel.Int(100),
			Offset: kernel.Int(0),
			Type:   kernel.AuthConnectionTimelineParamsTypeLogin,
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
