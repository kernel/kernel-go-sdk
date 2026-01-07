// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"context"
	"os"
	"testing"

	"github.com/onkernel/kernel-go-sdk"
	"github.com/onkernel/kernel-go-sdk/internal/testutil"
	"github.com/onkernel/kernel-go-sdk/option"
)

func TestManualPagination(t *testing.T) {
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
	page, err := client.Deployments.List(context.TODO(), kernel.DeploymentListParams{
		AppName: kernel.String("YOUR_APP"),
		Limit:   kernel.Int(2),
	})
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	for _, deployment := range page.Items {
		t.Logf("%+v\n", deployment.ID)
	}
	// Prism mock isn't going to give us real pagination
	page, err = page.GetNextPage()
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if page != nil {
		for _, deployment := range page.Items {
			t.Logf("%+v\n", deployment.ID)
		}
	}
}
