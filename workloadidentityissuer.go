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
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// WorkloadIdentityIssuerService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkloadIdentityIssuerService] method instead.
type WorkloadIdentityIssuerService struct {
	Options []option.RequestOption
}

// NewWorkloadIdentityIssuerService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewWorkloadIdentityIssuerService(opts ...option.RequestOption) (r WorkloadIdentityIssuerService) {
	r = WorkloadIdentityIssuerService{}
	r.Options = opts
	return
}

// Organization credential administrators only. No external namespace-ownership or
// staff approval is required. New trust resources default disabled.
func (r *WorkloadIdentityIssuerService) New(ctx context.Context, body WorkloadIdentityIssuerNewParams, opts ...option.RequestOption) (res *WorkloadIdentityIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/workload_identity_issuers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityIssuerService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *WorkloadIdentityIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/workload_identity_issuers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a workload identity issuer
func (r *WorkloadIdentityIssuerService) Update(ctx context.Context, id string, body WorkloadIdentityIssuerUpdateParams, opts ...option.RequestOption) (res *WorkloadIdentityIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/workload_identity_issuers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityIssuerService) List(ctx context.Context, query WorkloadIdentityIssuerListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[WorkloadIdentityIssuer], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/workload_identity_issuers"
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

// Any signed-in member of the organization can read this. API keys cannot.
func (r *WorkloadIdentityIssuerService) ListAutoPaging(ctx context.Context, query WorkloadIdentityIssuerListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[WorkloadIdentityIssuer] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Retains history and invalidates associated authority within the documented
// status bound. Does not delete created browser resources.
func (r *WorkloadIdentityIssuerService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/workload_identity_issuers/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Read-only. Does not enable trust, approve namespace ownership, create a grant or
// issue credentials.
func (r *WorkloadIdentityIssuerService) Verify(ctx context.Context, id string, opts ...option.RequestOption) (res *WorkloadIdentityIssuerKeyCheck, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/workload_identity_issuers/%s/verify", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type WorkloadIdentityIssuer struct {
	ID              string    `json:"id" api:"required"`
	ArchivedAt      time.Time `json:"archived_at" api:"required" format:"date-time"`
	CreatedAt       time.Time `json:"created_at" api:"required" format:"date-time"`
	CreatedByUserID string    `json:"created_by_user_id" api:"required"`
	Description     string    `json:"description" api:"required"`
	Enabled         bool      `json:"enabled" api:"required"`
	// Exact expected signed issuer. Optional for github_actions and google_cloud,
	// which each have one fixed issuer.
	IssuerURL string                    `json:"issuer_url" api:"required" format:"uri"`
	Jwks      WorkloadIdentityKeySource `json:"jwks" api:"required"`
	Name      string                    `json:"name" api:"required"`
	OrgID     string                    `json:"org_id" api:"required"`
	// The kind of identity provider. GitHub Actions, Google Cloud and Microsoft Entra
	// issuers are recognized by their issuer URL and must use their own type. Defaults
	// to the type the issuer URL belongs to, or custom_oidc.
	//
	// Any of "github_actions", "aws", "google_cloud", "microsoft_entra", "kubernetes",
	// "custom_oidc".
	PresetType WorkloadIdentityPresetType `json:"preset_type" api:"required"`
	// Authority revision. Cosmetic updates do not widen permissions.
	SecurityVersion int64                           `json:"security_version" api:"required"`
	TokenValidation WorkloadIdentityTokenValidation `json:"token_validation" api:"required"`
	UpdatedAt       time.Time                       `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		ArchivedAt      respjson.Field
		CreatedAt       respjson.Field
		CreatedByUserID respjson.Field
		Description     respjson.Field
		Enabled         respjson.Field
		IssuerURL       respjson.Field
		Jwks            respjson.Field
		Name            respjson.Field
		OrgID           respjson.Field
		PresetType      respjson.Field
		SecurityVersion respjson.Field
		TokenValidation respjson.Field
		UpdatedAt       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityIssuer) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityIssuer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Name is required.
type WorkloadIdentityIssuerConfigParam struct {
	Name        string            `json:"name" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	Enabled     param.Opt[bool]   `json:"enabled,omitzero"`
	// Exact expected signed issuer. Optional for github_actions and google_cloud,
	// which each have one fixed issuer.
	IssuerURL param.Opt[string]              `json:"issuer_url,omitzero" format:"uri"`
	Jwks      WorkloadIdentityKeySourceParam `json:"jwks,omitzero"`
	// The kind of identity provider. GitHub Actions, Google Cloud and Microsoft Entra
	// issuers are recognized by their issuer URL and must use their own type. Defaults
	// to the type the issuer URL belongs to, or custom_oidc.
	//
	// Any of "github_actions", "aws", "google_cloud", "microsoft_entra", "kubernetes",
	// "custom_oidc".
	PresetType      WorkloadIdentityPresetType           `json:"preset_type,omitzero"`
	TokenValidation WorkloadIdentityTokenValidationParam `json:"token_validation,omitzero"`
	paramObj
}

func (r WorkloadIdentityIssuerConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityIssuerConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityIssuerConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityIssuerKeyCheck struct {
	IssuerID string   `json:"issuer_id" api:"required"`
	KeyIDs   []string `json:"key_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		IssuerID    respjson.Field
		KeyIDs      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityIssuerKeyCheck) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityIssuerKeyCheck) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityKeySource struct {
	// Any of "discovery", "explicit_url", "inline".
	Type WorkloadIdentityKeySourceType `json:"type" api:"required"`
	// Optional public CA certificates. When set, only these CAs are trusted for key
	// retrieval, in place of the system roots. Private keys and private-network
	// retrieval are not supported.
	CaCertPem        string           `json:"ca_cert_pem"`
	DiscoveryBaseURL string           `json:"discovery_base_url" format:"uri"`
	Keys             []map[string]any `json:"keys"`
	URL              string           `json:"url" format:"uri"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type             respjson.Field
		CaCertPem        respjson.Field
		DiscoveryBaseURL respjson.Field
		Keys             respjson.Field
		URL              respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityKeySource) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityKeySource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this WorkloadIdentityKeySource to a
// WorkloadIdentityKeySourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WorkloadIdentityKeySourceParam.Overrides()
func (r WorkloadIdentityKeySource) ToParam() WorkloadIdentityKeySourceParam {
	return param.Override[WorkloadIdentityKeySourceParam](json.RawMessage(r.RawJSON()))
}

type WorkloadIdentityKeySourceType string

const (
	WorkloadIdentityKeySourceTypeDiscovery   WorkloadIdentityKeySourceType = "discovery"
	WorkloadIdentityKeySourceTypeExplicitURL WorkloadIdentityKeySourceType = "explicit_url"
	WorkloadIdentityKeySourceTypeInline      WorkloadIdentityKeySourceType = "inline"
)

// The property Type is required.
type WorkloadIdentityKeySourceParam struct {
	// Any of "discovery", "explicit_url", "inline".
	Type WorkloadIdentityKeySourceType `json:"type,omitzero" api:"required"`
	// Optional public CA certificates. When set, only these CAs are trusted for key
	// retrieval, in place of the system roots. Private keys and private-network
	// retrieval are not supported.
	CaCertPem        param.Opt[string] `json:"ca_cert_pem,omitzero"`
	DiscoveryBaseURL param.Opt[string] `json:"discovery_base_url,omitzero" format:"uri"`
	URL              param.Opt[string] `json:"url,omitzero" format:"uri"`
	Keys             []map[string]any  `json:"keys,omitzero"`
	paramObj
}

func (r WorkloadIdentityKeySourceParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityKeySourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityKeySourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The kind of identity provider. GitHub Actions, Google Cloud and Microsoft Entra
// issuers are recognized by their issuer URL and must use their own type. Defaults
// to the type the issuer URL belongs to, or custom_oidc.
type WorkloadIdentityPresetType string

const (
	WorkloadIdentityPresetTypeGitHubActions  WorkloadIdentityPresetType = "github_actions"
	WorkloadIdentityPresetTypeAws            WorkloadIdentityPresetType = "aws"
	WorkloadIdentityPresetTypeGoogleCloud    WorkloadIdentityPresetType = "google_cloud"
	WorkloadIdentityPresetTypeMicrosoftEntra WorkloadIdentityPresetType = "microsoft_entra"
	WorkloadIdentityPresetTypeKubernetes     WorkloadIdentityPresetType = "kubernetes"
	WorkloadIdentityPresetTypeCustomOidc     WorkloadIdentityPresetType = "custom_oidc"
)

type WorkloadIdentityTokenValidation struct {
	// Maximum incoming signed exp minus iat, independent of outgoing token lifetime.
	MaxLifetimeSeconds int64 `json:"max_lifetime_seconds"`
	// Whether an incoming token can be exchanged more than once before it expires.
	// Defaults to single_use for github_actions, aws, microsoft_entra and custom_oidc,
	// and to allow_reuse for kubernetes, because a pod reads the same projected token
	// until it is rotated. google_cloud is always allow_reuse because Google identity
	// tokens carry no token ID. allow_reuse can be set on microsoft_entra and
	// custom_oidc, and single_use on kubernetes. github_actions and aws are always
	// single_use.
	//
	// Any of "single_use", "allow_reuse".
	ReplayProtection WorkloadIdentityTokenValidationReplayProtection `json:"replay_protection"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MaxLifetimeSeconds respjson.Field
		ReplayProtection   respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityTokenValidation) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityTokenValidation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this WorkloadIdentityTokenValidation to a
// WorkloadIdentityTokenValidationParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// WorkloadIdentityTokenValidationParam.Overrides()
func (r WorkloadIdentityTokenValidation) ToParam() WorkloadIdentityTokenValidationParam {
	return param.Override[WorkloadIdentityTokenValidationParam](json.RawMessage(r.RawJSON()))
}

// Whether an incoming token can be exchanged more than once before it expires.
// Defaults to single_use for github_actions, aws, microsoft_entra and custom_oidc,
// and to allow_reuse for kubernetes, because a pod reads the same projected token
// until it is rotated. google_cloud is always allow_reuse because Google identity
// tokens carry no token ID. allow_reuse can be set on microsoft_entra and
// custom_oidc, and single_use on kubernetes. github_actions and aws are always
// single_use.
type WorkloadIdentityTokenValidationReplayProtection string

const (
	WorkloadIdentityTokenValidationReplayProtectionSingleUse  WorkloadIdentityTokenValidationReplayProtection = "single_use"
	WorkloadIdentityTokenValidationReplayProtectionAllowReuse WorkloadIdentityTokenValidationReplayProtection = "allow_reuse"
)

type WorkloadIdentityTokenValidationParam struct {
	// Maximum incoming signed exp minus iat, independent of outgoing token lifetime.
	MaxLifetimeSeconds param.Opt[int64] `json:"max_lifetime_seconds,omitzero"`
	// Whether an incoming token can be exchanged more than once before it expires.
	// Defaults to single_use for github_actions, aws, microsoft_entra and custom_oidc,
	// and to allow_reuse for kubernetes, because a pod reads the same projected token
	// until it is rotated. google_cloud is always allow_reuse because Google identity
	// tokens carry no token ID. allow_reuse can be set on microsoft_entra and
	// custom_oidc, and single_use on kubernetes. github_actions and aws are always
	// single_use.
	//
	// Any of "single_use", "allow_reuse".
	ReplayProtection WorkloadIdentityTokenValidationReplayProtection `json:"replay_protection,omitzero"`
	paramObj
}

func (r WorkloadIdentityTokenValidationParam) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityTokenValidationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityTokenValidationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityIssuerNewParams struct {
	WorkloadIdentityIssuerConfig WorkloadIdentityIssuerConfigParam
	paramObj
}

func (r WorkloadIdentityIssuerNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.WorkloadIdentityIssuerConfig)
}
func (r *WorkloadIdentityIssuerNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityIssuerUpdateParams struct {
	ExpectedVersion int64             `json:"expected_version" api:"required"`
	Description     param.Opt[string] `json:"description,omitzero"`
	Enabled         param.Opt[bool]   `json:"enabled,omitzero"`
	// Exact expected signed issuer. Optional for github_actions and google_cloud,
	// which each have one fixed issuer.
	IssuerURL param.Opt[string]              `json:"issuer_url,omitzero" format:"uri"`
	Name      param.Opt[string]              `json:"name,omitzero"`
	Jwks      WorkloadIdentityKeySourceParam `json:"jwks,omitzero"`
	// The kind of identity provider. GitHub Actions, Google Cloud and Microsoft Entra
	// issuers are recognized by their issuer URL and must use their own type. Defaults
	// to the type the issuer URL belongs to, or custom_oidc.
	//
	// Any of "github_actions", "aws", "google_cloud", "microsoft_entra", "kubernetes",
	// "custom_oidc".
	PresetType      WorkloadIdentityPresetType           `json:"preset_type,omitzero"`
	TokenValidation WorkloadIdentityTokenValidationParam `json:"token_validation,omitzero"`
	paramObj
}

func (r WorkloadIdentityIssuerUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow WorkloadIdentityIssuerUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WorkloadIdentityIssuerUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityIssuerListParams struct {
	Enabled param.Opt[bool]  `query:"enabled,omitzero" json:"-"`
	Limit   param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset  param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against name or description; exact ID match.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	// Any of "active", "archived", "all".
	Status WorkloadIdentityIssuerListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [WorkloadIdentityIssuerListParams]'s query parameters as
// `url.Values`.
func (r WorkloadIdentityIssuerListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WorkloadIdentityIssuerListParamsStatus string

const (
	WorkloadIdentityIssuerListParamsStatusActive   WorkloadIdentityIssuerListParamsStatus = "active"
	WorkloadIdentityIssuerListParamsStatusArchived WorkloadIdentityIssuerListParamsStatus = "archived"
	WorkloadIdentityIssuerListParamsStatusAll      WorkloadIdentityIssuerListParamsStatus = "all"
)
