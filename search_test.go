// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/internal/testutil"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestSearchNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Search.New(context.TODO(), kernel.SearchNewParams{
		Request: kernel.RequestParam{
			Query: "x",
			Content: kernel.RequestContentUnionParam{
				OfRequestContentBoolean: kernel.Bool(true),
			},
			Country:        kernel.String("se"),
			EndDate:        kernel.Time(time.Now()),
			ExcludeDomains: []string{"string"},
			IncludeDomains: []string{"string"},
			IncludeRaw:     kernel.Bool(true),
			Language:       kernel.String("language"),
			MaxResults:     kernel.Int(1),
			Recency:        kernel.RequestRecencyHour,
			SafeSearch:     kernel.RequestSafeSearchOff,
			StartDate:      kernel.Time(time.Now()),
			Strategy: kernel.StrategyUnionParam{
				OfAuto: &kernel.StrategyAutoParam{
					FallbackOn: []string{"error"},
					ProviderOptions: []kernel.ProviderTargetUnionParam{{
						OfBrave: &kernel.ProviderTargetBraveParam{
							Options: kernel.ProviderTargetBraveOptionsParam{
								Count:                kernel.Int(1),
								ExtraSnippets:        kernel.Bool(true),
								Goggles:              kernel.String("goggles"),
								GogglesID:            kernel.String("goggles_id"),
								IncludeFetchMetadata: kernel.Bool(true),
								Offset:               kernel.Int(0),
								Operators:            kernel.String("operators"),
								ResultFilter:         kernel.String("result_filter"),
								SearchLang:           kernel.String("search_lang"),
								Spellcheck:           kernel.Bool(true),
								UiLang:               kernel.String("ui_lang"),
								Units:                "metric",
							},
						},
					}},
				},
			},
			StrictParams: kernel.Bool(true),
			TimeoutMs:    kernel.Int(1000),
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

func TestSearchGet(t *testing.T) {
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
	_, err := client.Search.Get(context.TODO(), "srch_abc123")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
