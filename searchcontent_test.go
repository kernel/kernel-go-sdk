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

func TestSearchContentFetchWithOptionalParams(t *testing.T) {
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
	err := client.Search.Contents.Fetch(
		context.TODO(),
		"srch_abc123",
		kernel.SearchContentFetchParams{
			FetchRequest: kernel.FetchRequestParam{
				Content: kernel.FetchRequestContentParam{
					Browser: kernel.FetchRequestContentBrowserParam{
						BrowserID: kernel.String("browser_id"),
						Mode:      "curl",
					},
					Format:      "markdown",
					MaxAgeHours: kernel.Int(0),
					MaxChars:    kernel.Int(100),
					Source:      "auto",
					TimeoutMs:   kernel.Int(1000),
				},
				Limit:     kernel.Int(1),
				ResultIDs: []string{"string"},
				TimeoutMs: kernel.Int(1000),
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
