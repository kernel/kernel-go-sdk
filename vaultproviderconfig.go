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

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
	"github.com/kernel/kernel-go-sdk/shared/constant"
)

// VaultProviderConfigService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVaultProviderConfigService] method instead.
type VaultProviderConfigService struct {
	Options []option.RequestOption
}

// NewVaultProviderConfigService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewVaultProviderConfigService(opts ...option.RequestOption) (r VaultProviderConfigService) {
	r = VaultProviderConfigService{}
	r.Options = opts
	return
}

// Register a configuration shared across the organization's projects. Names are
// unique within the organization; duplicate names return 409 without replacing
// credentials. A configuration serves many wallets. Secret credentials are never
// returned. Requires an organization-scoped credential or dashboard
// authentication; project-scoped credentials receive 403.
func (r *VaultProviderConfigService) New(ctx context.Context, body VaultProviderConfigNewParams, opts ...option.RequestOption) (res *VaultProviderConfigUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "vault-provider-configs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Look up a configuration by ID or name. Returns 404 when it does not exist in the
// organization.
func (r *VaultProviderConfigService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *VaultProviderConfigUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("vault-provider-configs/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update the supplied fields; omitted fields remain unchanged. Names must remain
// unique within the organization. Requires an organization-scoped credential or
// dashboard authentication; project-scoped credentials receive 403.
func (r *VaultProviderConfigService) Update(ctx context.Context, idOrName string, body VaultProviderConfigUpdateParams, opts ...option.RequestOption) (res *VaultProviderConfigUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("vault-provider-configs/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Secret credentials are never returned.
func (r *VaultProviderConfigService) List(ctx context.Context, query VaultProviderConfigListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[VaultProviderConfigUnion], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "vault-provider-configs"
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

// Secret credentials are never returned.
func (r *VaultProviderConfigService) ListAutoPaging(ctx context.Context, query VaultProviderConfigListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[VaultProviderConfigUnion] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a configuration in the organization. Returns 409 while any non-deleted
// vault item references the configuration, regardless of connection status. Does
// not delete the external OAuth client or revoke unrelated grants. Requires an
// organization-scoped credential or dashboard authentication; project-scoped
// credentials receive 403.
func (r *VaultProviderConfigService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("vault-provider-configs/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// VaultProviderConfigUnion contains all possible properties and values from
// [VaultProviderConfigLink], [VaultProviderConfigAgentcard].
//
// Use the [VaultProviderConfigUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type VaultProviderConfigUnion struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	// Any of "link", "agentcard".
	Provider  string    `json:"provider"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [VaultProviderConfigAgentcard].
	TestMode bool `json:"test_mode"`
	JSON     struct {
		ID        respjson.Field
		ClientID  respjson.Field
		CreatedAt respjson.Field
		Name      respjson.Field
		Provider  respjson.Field
		UpdatedAt respjson.Field
		TestMode  respjson.Field
		raw       string
	} `json:"-"`
}

// anyVaultProviderConfig is implemented by each variant of
// [VaultProviderConfigUnion] to add type safety for the return type of
// [VaultProviderConfigUnion.AsAny]
type anyVaultProviderConfig interface {
	implVaultProviderConfigUnion()
}

func (VaultProviderConfigLink) implVaultProviderConfigUnion()      {}
func (VaultProviderConfigAgentcard) implVaultProviderConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := VaultProviderConfigUnion.AsAny().(type) {
//	case kernel.VaultProviderConfigLink:
//	case kernel.VaultProviderConfigAgentcard:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u VaultProviderConfigUnion) AsAny() anyVaultProviderConfig {
	switch u.Provider {
	case "link":
		return u.AsLink()
	case "agentcard":
		return u.AsAgentcard()
	}
	return nil
}

func (u VaultProviderConfigUnion) AsLink() (v VaultProviderConfigLink) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u VaultProviderConfigUnion) AsAgentcard() (v VaultProviderConfigAgentcard) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u VaultProviderConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *VaultProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response schema for a Link configuration, without secret credentials. Kernel
// generates the ID and timestamps. Configuration creation uses
// VaultLinkProviderConfigRequest.
type VaultProviderConfigLink struct {
	ID string `json:"id" api:"required"`
	// OAuth client identity; immutable. Secret credentials are never returned.
	ClientID  string    `json:"client_id" api:"required"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Unique within the organization.
	Name      string        `json:"name" api:"required"`
	Provider  constant.Link `json:"provider" default:"link"`
	UpdatedAt time.Time     `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ClientID    respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Provider    respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultProviderConfigLink) RawJSON() string { return r.JSON.raw }
func (r *VaultProviderConfigLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response schema for an AgentCard configuration, without secret credentials.
// Kernel generates the ID and timestamps and introspects test_mode from the
// credentials. Configuration creation uses VaultAgentCardProviderConfigRequest.
type VaultProviderConfigAgentcard struct {
	ID        string             `json:"id" api:"required"`
	ClientID  string             `json:"client_id" api:"required"`
	CreatedAt time.Time          `json:"created_at" api:"required" format:"date-time"`
	Name      string             `json:"name" api:"required"`
	Provider  constant.Agentcard `json:"provider" default:"agentcard"`
	// Introspected mode of the selected credential; true means sandbox objects.
	TestMode  bool      `json:"test_mode" api:"required"`
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ClientID    respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Provider    respjson.Field
		TestMode    respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VaultProviderConfigAgentcard) RawJSON() string { return r.JSON.raw }
func (r *VaultProviderConfigAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultProviderConfigNewParams struct {

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	// Register a customer-owned Link OAuth client for refreshing and revoking imported
	// wallet grants. The customer is responsible for initiating OAuth flows to connect
	// a user's wallet and handling the redirect to obtain a Link access token and
	// refresh token.
	OfLink *VaultProviderConfigNewParamsBodyLink `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	// Register application credentials for AgentCard wallet enrollment and checkout
	// approvals. Kernel obtains application access tokens using client_credentials.
	OfAgentcard *VaultProviderConfigNewParamsBodyAgentcard `json:",inline"`

	paramObj
}

func (u VaultProviderConfigNewParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLink, u.OfAgentcard)
}
func (r *VaultProviderConfigNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Register a customer-owned Link OAuth client for refreshing and revoking imported
// wallet grants. The customer is responsible for initiating OAuth flows to connect
// a user's wallet and handling the redirect to obtain a Link access token and
// refresh token.
//
// The properties Credentials, Name, Provider are required.
type VaultProviderConfigNewParamsBodyLink struct {
	Credentials VaultProviderConfigNewParamsBodyLinkCredentials `json:"credentials,omitzero" api:"required"`
	// Unique within the organization.
	Name string `json:"name" api:"required"`
	// This field can be elided, and will marshal its zero value as "link".
	Provider constant.Link `json:"provider" default:"link"`
	paramObj
}

func (r VaultProviderConfigNewParamsBodyLink) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigNewParamsBodyLink
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigNewParamsBodyLink) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ClientID, ClientSecret are required.
type VaultProviderConfigNewParamsBodyLinkCredentials struct {
	ClientID     string `json:"client_id" api:"required"`
	ClientSecret string `json:"client_secret" api:"required"`
	paramObj
}

func (r VaultProviderConfigNewParamsBodyLinkCredentials) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigNewParamsBodyLinkCredentials
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigNewParamsBodyLinkCredentials) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Register application credentials for AgentCard wallet enrollment and checkout
// approvals. Kernel obtains application access tokens using client_credentials.
//
// The properties Credentials, Name, Provider are required.
type VaultProviderConfigNewParamsBodyAgentcard struct {
	Credentials VaultProviderConfigNewParamsBodyAgentcardCredentials `json:"credentials,omitzero" api:"required"`
	// Unique within the organization.
	Name string `json:"name" api:"required"`
	// This field can be elided, and will marshal its zero value as "agentcard".
	Provider constant.Agentcard `json:"provider" default:"agentcard"`
	paramObj
}

func (r VaultProviderConfigNewParamsBodyAgentcard) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigNewParamsBodyAgentcard
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigNewParamsBodyAgentcard) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ClientID, ClientSecret are required.
type VaultProviderConfigNewParamsBodyAgentcardCredentials struct {
	ClientID     string `json:"client_id" api:"required"`
	ClientSecret string `json:"client_secret" api:"required"`
	paramObj
}

func (r VaultProviderConfigNewParamsBodyAgentcardCredentials) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigNewParamsBodyAgentcardCredentials
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigNewParamsBodyAgentcardCredentials) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultProviderConfigUpdateParams struct {
	// Unique within the organization.
	Name param.Opt[string] `json:"name,omitzero"`
	// Fields to update. Omitted credentials are left unchanged. A rejected update
	// leaves existing credentials unchanged.
	Credentials VaultProviderConfigUpdateParamsCredentials `json:"credentials,omitzero"`
	paramObj
}

func (r VaultProviderConfigUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Fields to update. Omitted credentials are left unchanged. A rejected update
// leaves existing credentials unchanged.
type VaultProviderConfigUpdateParamsCredentials struct {
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r VaultProviderConfigUpdateParamsCredentials) MarshalJSON() (data []byte, err error) {
	type shadow VaultProviderConfigUpdateParamsCredentials
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultProviderConfigUpdateParamsCredentials) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultProviderConfigListParams struct {
	Limit  param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultProviderConfigListParams]'s query parameters as
// `url.Values`.
func (r VaultProviderConfigListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
