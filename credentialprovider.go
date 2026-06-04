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
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Configure external credential providers like 1Password.
//
// CredentialProviderService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCredentialProviderService] method instead.
type CredentialProviderService struct {
	Options []option.RequestOption
}

// NewCredentialProviderService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewCredentialProviderService(opts ...option.RequestOption) (r CredentialProviderService) {
	r = CredentialProviderService{}
	r.Options = opts
	return
}

// Configure an external credential provider (e.g., 1Password) for automatic
// credential lookup.
func (r *CredentialProviderService) New(ctx context.Context, body CredentialProviderNewParams, opts ...option.RequestOption) (res *CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/credential_providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a credential provider by its ID.
func (r *CredentialProviderService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/credential_providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a credential provider's configuration.
func (r *CredentialProviderService) Update(ctx context.Context, id string, body CredentialProviderUpdateParams, opts ...option.RequestOption) (res *CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/credential_providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List external credential providers configured for the organization.
func (r *CredentialProviderService) List(ctx context.Context, query CredentialProviderListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[CredentialProvider], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/credential_providers"
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

// List external credential providers configured for the organization.
func (r *CredentialProviderService) ListAutoPaging(ctx context.Context, query CredentialProviderListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[CredentialProvider] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a credential provider by its ID.
func (r *CredentialProviderService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/credential_providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Returns available credential items (e.g., 1Password login items) from the
// provider.
func (r *CredentialProviderService) ListItems(ctx context.Context, id string, opts ...option.RequestOption) (res *CredentialProviderListItemsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/credential_providers/%s/items", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Validate the credential provider's token and list accessible vaults.
func (r *CredentialProviderService) Test(ctx context.Context, id string, opts ...option.RequestOption) (res *CredentialProviderTestResult, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/credential_providers/%s/test", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Request to create an external credential provider
//
// The properties Token, Name, ProviderType are required.
type CreateCredentialProviderRequestParam struct {
	// Service account token for the provider (e.g., 1Password service account token)
	Token string `json:"token" api:"required"`
	// Human-readable name for this provider instance (unique per org)
	Name string `json:"name" api:"required"`
	// Type of credential provider
	//
	// Any of "onepassword".
	ProviderType CreateCredentialProviderRequestProviderType `json:"provider_type,omitzero" api:"required"`
	// How long to cache credential lists (default 300 seconds)
	CacheTtlSeconds param.Opt[int64] `json:"cache_ttl_seconds,omitzero"`
	paramObj
}

func (r CreateCredentialProviderRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateCredentialProviderRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CreateCredentialProviderRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of credential provider
type CreateCredentialProviderRequestProviderType string

const (
	CreateCredentialProviderRequestProviderTypeOnepassword CreateCredentialProviderRequestProviderType = "onepassword"
)

// An external credential provider (e.g., 1Password) for automatic credential
// lookup
type CredentialProvider struct {
	// Unique identifier for the credential provider
	ID string `json:"id" api:"required"`
	// When the credential provider was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Whether the provider is enabled for credential lookups
	Enabled bool `json:"enabled" api:"required"`
	// Human-readable name for this provider instance
	Name string `json:"name" api:"required"`
	// Priority order for credential lookups (lower numbers are checked first)
	Priority int64 `json:"priority" api:"required"`
	// Type of credential provider
	//
	// Any of "onepassword".
	ProviderType CredentialProviderProviderType `json:"provider_type" api:"required"`
	// When the credential provider was last updated
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		CreatedAt    respjson.Field
		Enabled      respjson.Field
		Name         respjson.Field
		Priority     respjson.Field
		ProviderType respjson.Field
		UpdatedAt    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialProvider) RawJSON() string { return r.JSON.raw }
func (r *CredentialProvider) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of credential provider
type CredentialProviderProviderType string

const (
	CredentialProviderProviderTypeOnepassword CredentialProviderProviderType = "onepassword"
)

// A credential item from an external provider (e.g., a 1Password login item)
type CredentialProviderItem struct {
	// Unique identifier for the item within the provider
	ID string `json:"id" api:"required"`
	// Path to reference this item (VaultName/ItemTitle format)
	Path string `json:"path" api:"required"`
	// Display name of the credential item
	Title string `json:"title" api:"required"`
	// ID of the vault containing this item
	VaultID string `json:"vault_id" api:"required"`
	// Name of the vault containing this item
	VaultName string `json:"vault_name" api:"required"`
	// URLs associated with this credential
	URLs []string `json:"urls"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Path        respjson.Field
		Title       respjson.Field
		VaultID     respjson.Field
		VaultName   respjson.Field
		URLs        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialProviderItem) RawJSON() string { return r.JSON.raw }
func (r *CredentialProviderItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Result of testing a credential provider connection
type CredentialProviderTestResult struct {
	// Whether the connection test was successful
	Success bool `json:"success" api:"required"`
	// List of vaults accessible by the service account
	Vaults []CredentialProviderTestResultVault `json:"vaults" api:"required"`
	// Error message if the test failed
	Error string `json:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Success     respjson.Field
		Vaults      respjson.Field
		Error       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialProviderTestResult) RawJSON() string { return r.JSON.raw }
func (r *CredentialProviderTestResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialProviderTestResultVault struct {
	// Vault ID
	ID string `json:"id" api:"required"`
	// Vault name
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialProviderTestResultVault) RawJSON() string { return r.JSON.raw }
func (r *CredentialProviderTestResultVault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to update a credential provider
type UpdateCredentialProviderRequestParam struct {
	// New service account token (to rotate credentials)
	Token param.Opt[string] `json:"token,omitzero"`
	// How long to cache credential lists
	CacheTtlSeconds param.Opt[int64] `json:"cache_ttl_seconds,omitzero"`
	// Whether the provider is enabled for credential lookups
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Human-readable name for this provider instance
	Name param.Opt[string] `json:"name,omitzero"`
	// Priority order for credential lookups (lower numbers are checked first)
	Priority param.Opt[int64] `json:"priority,omitzero"`
	paramObj
}

func (r UpdateCredentialProviderRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateCredentialProviderRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateCredentialProviderRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialProviderListItemsResponse struct {
	Items []CredentialProviderItem `json:"items"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Items       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialProviderListItemsResponse) RawJSON() string { return r.JSON.raw }
func (r *CredentialProviderListItemsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialProviderNewParams struct {
	// Request to create an external credential provider
	CreateCredentialProviderRequest CreateCredentialProviderRequestParam
	paramObj
}

func (r CredentialProviderNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.CreateCredentialProviderRequest)
}
func (r *CredentialProviderNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialProviderUpdateParams struct {
	// Request to update a credential provider
	UpdateCredentialProviderRequest UpdateCredentialProviderRequestParam
	paramObj
}

func (r CredentialProviderUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateCredentialProviderRequest)
}
func (r *CredentialProviderUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialProviderListParams struct {
	// Limit the number of credential providers to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of credential providers to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CredentialProviderListParams]'s query parameters as
// `url.Values`.
func (r CredentialProviderListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
