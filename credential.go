// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/apiquery"
	shimjson "github.com/onkernel/kernel-go-sdk/internal/encoding/json"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/pagination"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

// CredentialService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCredentialService] method instead.
type CredentialService struct {
	Options []option.RequestOption
}

// NewCredentialService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewCredentialService(opts ...option.RequestOption) (r CredentialService) {
	r = CredentialService{}
	r.Options = opts
	return
}

// Create a new credential for storing login information.
func (r *CredentialService) New(ctx context.Context, body CredentialNewParams, opts ...option.RequestOption) (res *Credential, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "credentials"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve a credential by its ID or name. Credential values are not returned.
func (r *CredentialService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *Credential, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("credentials/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update a credential's name or values. When values are provided, they are merged
// with existing values (new keys are added, existing keys are overwritten).
func (r *CredentialService) Update(ctx context.Context, idOrName string, body CredentialUpdateParams, opts ...option.RequestOption) (res *Credential, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("credentials/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List credentials owned by the caller's organization. Credential values are not
// returned.
func (r *CredentialService) List(ctx context.Context, query CredentialListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Credential], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "credentials"
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

// List credentials owned by the caller's organization. Credential values are not
// returned.
func (r *CredentialService) ListAutoPaging(ctx context.Context, query CredentialListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Credential] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a credential by its ID or name.
func (r *CredentialService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("credentials/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Returns the current 6-digit TOTP code for a credential with a configured
// totp_secret. Use this to complete 2FA setup on sites or when you need a fresh
// code.
func (r *CredentialService) TotpCode(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *CredentialTotpCodeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("credentials/%s/totp-code", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Request to create a new credential
//
// The properties Domain, Name, Values are required.
type CreateCredentialRequestParam struct {
	// Target domain this credential is for
	Domain string `json:"domain,required"`
	// Unique name for the credential within the organization
	Name string `json:"name,required"`
	// Field name to value mapping (e.g., username, password)
	Values map[string]string `json:"values,omitzero,required"`
	// If set, indicates this credential should be used with the specified SSO provider
	// (e.g., google, github, microsoft). When the target site has a matching SSO
	// button, it will be clicked first before filling credential values on the
	// identity provider's login page.
	SSOProvider param.Opt[string] `json:"sso_provider,omitzero"`
	// Base32-encoded TOTP secret for generating one-time passwords. Used for automatic
	// 2FA during login.
	TotpSecret param.Opt[string] `json:"totp_secret,omitzero"`
	paramObj
}

func (r CreateCredentialRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateCredentialRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CreateCredentialRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A stored credential for automatic re-authentication
type Credential struct {
	// Unique identifier for the credential
	ID string `json:"id,required"`
	// When the credential was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Target domain this credential is for
	Domain string `json:"domain,required"`
	// Unique name for the credential within the organization
	Name string `json:"name,required"`
	// When the credential was last updated
	UpdatedAt time.Time `json:"updated_at,required" format:"date-time"`
	// Whether this credential has a TOTP secret configured for automatic 2FA
	HasTotpSecret bool `json:"has_totp_secret"`
	// If set, indicates this credential should be used with the specified SSO provider
	// (e.g., google, github, microsoft). When the target site has a matching SSO
	// button, it will be clicked first before filling credential values on the
	// identity provider's login page.
	SSOProvider string `json:"sso_provider,nullable"`
	// Current 6-digit TOTP code. Only included in create/update responses when
	// totp_secret was just set.
	TotpCode string `json:"totp_code"`
	// When the totp_code expires. Only included when totp_code is present.
	TotpCodeExpiresAt time.Time `json:"totp_code_expires_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		Domain            respjson.Field
		Name              respjson.Field
		UpdatedAt         respjson.Field
		HasTotpSecret     respjson.Field
		SSOProvider       respjson.Field
		TotpCode          respjson.Field
		TotpCodeExpiresAt respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Credential) RawJSON() string { return r.JSON.raw }
func (r *Credential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to update an existing credential
type UpdateCredentialRequestParam struct {
	// If set, indicates this credential should be used with the specified SSO
	// provider. Set to empty string or null to remove.
	SSOProvider param.Opt[string] `json:"sso_provider,omitzero"`
	// New name for the credential
	Name param.Opt[string] `json:"name,omitzero"`
	// Base32-encoded TOTP secret for generating one-time passwords. Spaces and
	// formatting are automatically normalized. Set to empty string to remove.
	TotpSecret param.Opt[string] `json:"totp_secret,omitzero"`
	// Field name to value mapping. Values are merged with existing values (new keys
	// added, existing keys overwritten).
	Values map[string]string `json:"values,omitzero"`
	paramObj
}

func (r UpdateCredentialRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow UpdateCredentialRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UpdateCredentialRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialTotpCodeResponse struct {
	// Current 6-digit TOTP code
	Code string `json:"code,required"`
	// When this code expires (ISO 8601 timestamp)
	ExpiresAt time.Time `json:"expires_at,required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		ExpiresAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialTotpCodeResponse) RawJSON() string { return r.JSON.raw }
func (r *CredentialTotpCodeResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CredentialNewParams struct {
	// Request to create a new credential
	CreateCredentialRequest CreateCredentialRequestParam
	paramObj
}

func (r CredentialNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.CreateCredentialRequest)
}
func (r *CredentialNewParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.CreateCredentialRequest)
}

type CredentialUpdateParams struct {
	// Request to update an existing credential
	UpdateCredentialRequest UpdateCredentialRequestParam
	paramObj
}

func (r CredentialUpdateParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.UpdateCredentialRequest)
}
func (r *CredentialUpdateParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.UpdateCredentialRequest)
}

type CredentialListParams struct {
	// Filter by domain
	Domain param.Opt[string] `query:"domain,omitzero" json:"-"`
	// Maximum number of results to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of results to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CredentialListParams]'s query parameters as `url.Values`.
func (r CredentialListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
