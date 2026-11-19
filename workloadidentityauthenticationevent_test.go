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

func TestWorkloadIdentityAuthenticationEventListWithOptionalParams(t *testing.T) {
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
	_, err := client.WorkloadIdentityAuthenticationEvents.List(context.TODO(), kernel.WorkloadIdentityAuthenticationEventListParams{
		EndTime:          kernel.Time(time.Now()),
		FederationRuleID: kernel.String("federation_rule_id"),
		IssuerID:         kernel.String("issuer_id"),
		Limit:            kernel.Int(1),
		Offset:           kernel.Int(0),
		Outcome:          kernel.WorkloadIdentityAuthenticationEventListParamsOutcomeSuccess,
		ProjectID:        kernel.String("project_id"),
		ServiceAccountID: kernel.String("service_account_id"),
		StartTime:        kernel.Time(time.Now()),
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
