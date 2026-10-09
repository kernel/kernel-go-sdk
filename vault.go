// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// VaultService contains methods and other services that help with interacting with
// the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVaultService] method instead.
type VaultService struct {
	Options []option.RequestOption
	Items   VaultItemService
}

// NewVaultService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewVaultService(opts ...option.RequestOption) (r VaultService) {
	r = VaultService{}
	r.Options = opts
	r.Items = NewVaultItemService(opts...)
	return
}

// Get a vault
func (r *VaultService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *Vault, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List vaults in the current project
func (r *VaultService) List(ctx context.Context, query VaultListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Vault], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "vaults"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List vaults in the current project
func (r *VaultService) ListAutoPaging(ctx context.Context, query VaultListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Vault] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Unresolved payment operations block deletion. Deleting a connected Kernel wallet
// first blocks new payments on it, then removes its enrolled card. If that fails,
// the wallet is kept and keeps refusing payments; retry the deletion. Reconcile
// the original attempt with the provider or support first; deleting or recreating
// an item is not proof that a payment did not occur.
func (r *VaultService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("vaults/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Returns the public key that custom credential collection web apps use to encrypt
// values in the browser. Use this only if you run your own credential collection
// web app and want values encrypted in the browser, sent to your backend still
// encrypted, and forwarded to Kernel's API still encrypted. In every other case,
// including server-side code that already holds the plaintext, use value. The page
// encrypts each value to this key, your backend forwards the ciphertext unchanged
// as encrypted_value when creating or updating a credential item in this vault,
// and Kernel decrypts it.
func (r *VaultService) GetEncryptionKey(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *VaultEncryptionKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/encryption_key", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Free organizations can store up to 3 non-deleted vaults across all projects.
// Paid plans and active trials have no vault cap. Retrieving an existing vault by
// name succeeds even at the limit.
func (r *VaultService) Upsert(ctx context.Context, body VaultUpsertParams, opts ...option.RequestOption) (res *Vault, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "vaults"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type Vault struct {
	ID        string    `json:"id" api:"required"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Immutable name assigned when the vault is created.
	Name      string    `json:"name" api:"required"`
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Vault) RawJSON() string { return r.JSON.raw }
func (r *Vault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Public key for encrypted_value on credential fields. Use this only if you run
// your own credential collection web app and want values encrypted in the browser,
// sent to your backend still encrypted, and forwarded to Kernel's API still
// encrypted. In every other case, including server-side code that already holds
// the plaintext, use value. Each vault has its own key; a value encrypted for one
// vault is rejected by every other vault. The key is created on first request and
// stays the same for the vault's lifetime, so it may be cached.
type VaultEncryptionKey struct {
	// JWE key management algorithm.
	//
	// Any of "ECDH-ES".
	Alg VaultEncryptionKeyAlg `json:"alg" api:"required"`
	// JWE content encryption algorithm.
	//
	// Any of "A256GCM".
	Enc VaultEncryptionKeyEnc `json:"enc" api:"required"`
	// P-256 public key in JWK form.
	Jwk VaultEncryptionKeyJwk `json:"jwk" api:"required"`
	// Key ID. Set it as the kid protected header of every encrypted_value.
	Kid string `json:"kid" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Alg         respjson.Field
		Enc         respjson.Field
		Jwk         respjson.Field
		Kid         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultEncryptionKey) RawJSON() string { return r.JSON.raw }
func (r *VaultEncryptionKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWE key management algorithm.
type VaultEncryptionKeyAlg string

const (
	VaultEncryptionKeyAlgEcdhEs VaultEncryptionKeyAlg = "ECDH-ES"
)

// JWE content encryption algorithm.
type VaultEncryptionKeyEnc string

const (
	VaultEncryptionKeyEncA256Gcm VaultEncryptionKeyEnc = "A256GCM"
)

// P-256 public key in JWK form.
type VaultEncryptionKeyJwk struct {
	// Any of "P-256".
	Crv string `json:"crv" api:"required"`
	// Any of "EC".
	Kty string `json:"kty" api:"required"`
	// Base64url-encoded x coordinate.
	X string `json:"x" api:"required"`
	// Base64url-encoded y coordinate.
	Y string `json:"y" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Crv         respjson.Field
		Kty         respjson.Field
		X           respjson.Field
		Y           respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultEncryptionKeyJwk) RawJSON() string { return r.JSON.raw }
func (r *VaultEncryptionKeyJwk) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultListParams struct {
	Limit  param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against vault name. IDs match by exact value.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultListParams]'s query parameters as `url.Values`.
func (r VaultListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type VaultUpsertParams struct {
	// Immutable name used to create or retrieve the vault.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r VaultUpsertParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultUpsertParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultUpsertParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
