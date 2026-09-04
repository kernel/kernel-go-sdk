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

func TestVaultItemGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Vaults.Items.Get(
		context.TODO(),
		"x",
		kernel.VaultItemGetParams{
			IDOrName: "id_or_name",
			Expand:   []string{"payment_methods"},
			Wait:     kernel.Int(0),
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

func TestVaultItemUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Vaults.Items.Update(
		context.TODO(),
		"x",
		kernel.VaultItemUpdateParams{
			IDOrName: "id_or_name",
			Spec: kernel.CardVaultItemSpecUnionParam{
				OfLink: &kernel.CardVaultItemSpecLinkParam{
					Amount:          1,
					Context:         "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Currency:        "bFx",
					MerchantName:    "x",
					MerchantURL:     "https://example.com",
					PaymentMethodID: "x",
					Test:            true,
					Wallet:          "wallet",
					ExpiresAt:       kernel.Int(0),
					LineItems: []kernel.CardVaultItemSpecLinkLineItemParam{{
						Name:        "name",
						Description: kernel.String("description"),
						ImageURL:    kernel.String("image_url"),
						ProductURL:  kernel.String("product_url"),
						Quantity:    kernel.Int(1),
						SKU:         kernel.String("sku"),
						Totals: []kernel.CardVaultItemSpecLinkLineItemTotalParam{{
							Amount:      0,
							DisplayText: "display_text",
							Type:        "type",
						}},
						UnitAmount: kernel.Int(0),
						URL:        kernel.String("url"),
					}},
					Metadata: map[string]string{
						"foo": "string",
					},
					Totals: []kernel.CardVaultItemSpecLinkTotalParam{{
						Amount:      0,
						DisplayText: "display_text",
						Type:        "type",
					}},
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

func TestVaultItemList(t *testing.T) {
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
	_, err := client.Vaults.Items.List(context.TODO(), "id_or_name")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVaultItemDelete(t *testing.T) {
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
	err := client.Vaults.Items.Delete(
		context.TODO(),
		"x",
		kernel.VaultItemDeleteParams{
			IDOrName: "id_or_name",
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

func TestVaultItemEventsWithOptionalParams(t *testing.T) {
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
	_, err := client.Vaults.Items.Events(
		context.TODO(),
		"key",
		kernel.VaultItemEventsParams{
			IDOrName: "id_or_name",
			After:    kernel.String("after"),
			Wait:     kernel.Int(0),
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

func TestVaultItemPerformOperation(t *testing.T) {
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
	_, err := client.Vaults.Items.PerformOperation(
		context.TODO(),
		"key",
		kernel.VaultItemPerformOperationParams{
			IDOrName: "id_or_name",
			Type:     kernel.VaultItemPerformOperationParamsTypeAuthorize,
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

func TestVaultItemUpsertWithOptionalParams(t *testing.T) {
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
	_, err := client.Vaults.Items.Upsert(
		context.TODO(),
		"x",
		kernel.VaultItemUpsertParams{
			IDOrName: "id_or_name",
			OfWallet: &kernel.VaultItemUpsertParamsBodyWallet{
				Spec: kernel.WalletVaultItemSpecUnionParam{
					OfLink: &kernel.WalletVaultItemSpecLinkParam{
						Authorization: kernel.WalletVaultItemSpecLinkAuthorizationParam{
							Client: kernel.WalletVaultItemSpecLinkAuthorizationClientParam{
								Type: "kernel_managed",
							},
							Method: "oauth",
						},
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
