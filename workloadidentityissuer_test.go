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

func TestWorkloadIdentityIssuerNewWithOptionalParams(t *testing.T) {
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
	_, err := client.WorkloadIdentityIssuers.New(context.TODO(), kernel.WorkloadIdentityIssuerNewParams{
		WorkloadIdentityIssuerConfig: kernel.WorkloadIdentityIssuerConfigParam{
			Name:        "name",
			Description: kernel.String("description"),
			Enabled:     kernel.Bool(true),
			IssuerURL:   kernel.String("https://example.com"),
			Jwks: kernel.WorkloadIdentityKeySourceParam{
				Type:             kernel.WorkloadIdentityKeySourceTypeDiscovery,
				CaCertPem:        kernel.String("ca_cert_pem"),
				DiscoveryBaseURL: kernel.String("https://example.com"),
				Keys: []map[string]any{{
					"foo": "bar",
				}},
				URL: kernel.String("https://example.com"),
			},
			PresetType: kernel.WorkloadIdentityPresetTypeGitHubActions,
			TokenValidation: kernel.WorkloadIdentityTokenValidationParam{
				MaxLifetimeSeconds: kernel.Int(1),
				ReplayProtection:   kernel.WorkloadIdentityTokenValidationReplayProtectionSingleUse,
			},
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

func TestWorkloadIdentityIssuerGet(t *testing.T) {
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
	_, err := client.WorkloadIdentityIssuers.Get(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityIssuerUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.WorkloadIdentityIssuers.Update(
		context.TODO(),
		"id",
		kernel.WorkloadIdentityIssuerUpdateParams{
			ExpectedVersion: 1,
			Description:     kernel.String("description"),
			Enabled:         kernel.Bool(true),
			IssuerURL:       kernel.String("https://example.com"),
			Jwks: kernel.WorkloadIdentityKeySourceParam{
				Type:             kernel.WorkloadIdentityKeySourceTypeDiscovery,
				CaCertPem:        kernel.String("ca_cert_pem"),
				DiscoveryBaseURL: kernel.String("https://example.com"),
				Keys: []map[string]any{{
					"foo": "bar",
				}},
				URL: kernel.String("https://example.com"),
			},
			Name:       kernel.String("name"),
			PresetType: kernel.WorkloadIdentityPresetTypeGitHubActions,
			TokenValidation: kernel.WorkloadIdentityTokenValidationParam{
				MaxLifetimeSeconds: kernel.Int(1),
				ReplayProtection:   kernel.WorkloadIdentityTokenValidationReplayProtectionSingleUse,
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

func TestWorkloadIdentityIssuerListWithOptionalParams(t *testing.T) {
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
	_, err := client.WorkloadIdentityIssuers.List(context.TODO(), kernel.WorkloadIdentityIssuerListParams{
		Enabled: kernel.Bool(true),
		Limit:   kernel.Int(1),
		Offset:  kernel.Int(0),
		Query:   kernel.String("query"),
		Status:  kernel.WorkloadIdentityIssuerListParamsStatusActive,
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityIssuerDelete(t *testing.T) {
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
	err := client.WorkloadIdentityIssuers.Delete(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityIssuerVerify(t *testing.T) {
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
	_, err := client.WorkloadIdentityIssuers.Verify(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
