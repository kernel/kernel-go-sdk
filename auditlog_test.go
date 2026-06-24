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

func TestAuditLogListWithOptionalParams(t *testing.T) {
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
	_, err := client.AuditLogs.List(context.TODO(), kernel.AuditLogListParams{
		End:           time.Now(),
		Start:         time.Now(),
		AuthStrategy:  kernel.String("auth_strategy"),
		ExcludeMethod: kernel.String("exclude_method"),
		Limit:         kernel.Int(1),
		Method:        kernel.String("method"),
		PageToken:     kernel.String("page_token"),
		Search:        kernel.String("search"),
		SearchUserID:  []string{"string"},
		Service:       kernel.String("service"),
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
