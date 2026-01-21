// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

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
	path := "org/credential-providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve a credential provider by its ID.
func (r *CredentialProviderService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("org/credential-providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update a credential provider's configuration.
func (r *CredentialProviderService) Update(ctx context.Context, id string, body CredentialProviderUpdateParams, opts ...option.RequestOption) (res *CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("org/credential-providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List external credential providers configured for the organization.
func (r *CredentialProviderService) List(ctx context.Context, opts ...option.RequestOption) (res *[]CredentialProvider, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/credential-providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Delete a credential provider by its ID.
func (r *CredentialProviderService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("org/credential-providers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Validate the credential provider's token and list accessible vaults.
func (r *CredentialProviderService) Test(ctx context.Context, id string, opts ...option.RequestOption) (res *CredentialProviderTestResult, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("org/credential-providers/%s/test", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Request to create an external credential provider
//
// The properties Token, ProviderType are required.
type CreateCredentialProviderRequestParam struct {
	// Service account token for the provider (e.g., 1Password service account token)
	Token string `json:"token,required"`
	// Type of credential provider
	//
	// Any of "onepassword".
	ProviderType CreateCredentialProviderRequestProviderType `json:"provider_type,omitzero,required"`
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
	ID string `json:"id,required"`
	// When the credential provider was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Whether the provider is enabled for credential lookups
	Enabled bool `json:"enabled,required"`
	// Priority order for credential lookups (lower numbers are checked first)
	Priority int64 `json:"priority,required"`
	// Type of credential provider
	//
	// Any of "onepassword".
	ProviderType CredentialProviderProviderType `json:"provider_type,required"`
	// When the credential provider was last updated
	UpdatedAt time.Time `json:"updated_at,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		CreatedAt    respjson.Field
		Enabled      respjson.Field
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

// Result of testing a credential provider connection
type CredentialProviderTestResult struct {
	// Whether the connection test was successful
	Success bool `json:"success,required"`
	// List of vaults accessible by the service account
	Vaults []CredentialProviderTestResultVault `json:"vaults,required"`
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
	ID string `json:"id,required"`
	// Vault name
	Name string `json:"name,required"`
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

type CredentialProviderNewParams struct {
	// Request to create an external credential provider
	CreateCredentialProviderRequest CreateCredentialProviderRequestParam
	paramObj
}

func (r CredentialProviderNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.CreateCredentialProviderRequest)
}
func (r *CredentialProviderNewParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.CreateCredentialProviderRequest)
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
	return json.Unmarshal(data, &r.UpdateCredentialProviderRequest)
}
