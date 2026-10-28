package kernel_test

import (
	"encoding/json"
	"testing"

	"github.com/kernel/kernel-go-sdk"
)

func TestVaultItemUnionSpecMetadata(t *testing.T) {
	data := []byte(`{"id":"item","key":"card","type":"card","spec":{"provider":"link","wallet":"wallet","metadata":{"order":"42"}}}`)
	var item kernel.VaultItemUnion
	if err := json.Unmarshal(data, &item); err != nil {
		t.Fatal(err)
	}
	if got := item.Spec.Metadata["order"]; got != "42" {
		t.Fatalf("VaultItemUnion.Spec.Metadata[order] = %q, want 42", got)
	}
	var response kernel.VaultItemOperationResponseUnion
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	if got := response.Spec.Metadata["order"]; got != "42" {
		t.Fatalf("VaultItemOperationResponseUnion.Spec.Metadata[order] = %q, want 42", got)
	}
}
