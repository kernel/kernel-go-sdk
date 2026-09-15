package kernel_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
)

type vaultRetryTransport func(*http.Request) (*http.Response, error)

func (f vaultRetryTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestVaultFillDoesNotRetry(t *testing.T) {
	for _, status := range []int{0, 409, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			client := kernel.NewClient(option.WithAPIKey("test"), option.WithMaxRetries(1), option.WithHTTPClient(&http.Client{Transport: vaultRetryTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				if status == 0 {
					return nil, io.ErrUnexpectedEOF
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader("{}")), Request: r}, nil
			})}))
			_, err := client.Vaults.Items.PerformOperation(context.Background(), "login", kernel.VaultItemPerformOperationParams{
				IDOrName: "vault", OfFill: &kernel.FillVaultItemOperationRequestParam{BrowserID: "browser", Fields: []kernel.VaultFillFieldParam{{Field: "password", Selector: "#password"}}},
			})
			if err == nil {
				t.Fatal("expected error")
			}
			if calls != 1 {
				t.Fatalf("got %d attempts, want 1", calls)
			}
		})
	}
}

func TestVaultReadKeepsClientRetries(t *testing.T) {
	calls := 0
	client := kernel.NewClient(option.WithAPIKey("test"), option.WithMaxRetries(1), option.WithHTTPClient(&http.Client{Transport: vaultRetryTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		return nil, io.ErrUnexpectedEOF
	})}))
	_, err := client.Vaults.Items.Get(context.Background(), "login", kernel.VaultItemGetParams{IDOrName: "vault"})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 2 {
		t.Fatalf("got %d attempts, want 2", calls)
	}
}
