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

// ServiceAccountService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewServiceAccountService] method instead.
type ServiceAccountService struct {
	Options []option.RequestOption
}

// NewServiceAccountService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewServiceAccountService(opts ...option.RequestOption) (r ServiceAccountService) {
	r = ServiceAccountService{}
	r.Options = opts
	return
}

// Organization credential administrators only. No external namespace-ownership or
// staff approval is required. New trust resources default disabled.
func (r *ServiceAccountService) New(ctx context.Context, body ServiceAccountNewParams, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/service_accounts"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *ServiceAccountService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/service_accounts/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a service account
func (r *ServiceAccountService) Update(ctx context.Context, id string, body ServiceAccountUpdateParams, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/service_accounts/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Any signed-in member of the organization can read this. API keys cannot.
func (r *ServiceAccountService) List(ctx context.Context, query ServiceAccountListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ServiceAccount], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/service_accounts"
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
func (r *ServiceAccountService) ListAutoPaging(ctx context.Context, query ServiceAccountListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ServiceAccount] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Retains history and invalidates associated authority within the documented
// status bound. Does not delete created browser resources.
func (r *ServiceAccountService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/service_accounts/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type ServiceAccount struct {
	ID              string    `json:"id" api:"required"`
	ArchivedAt      time.Time `json:"archived_at" api:"required" format:"date-time"`
	CreatedAt       time.Time `json:"created_at" api:"required" format:"date-time"`
	CreatedByUserID string    `json:"created_by_user_id" api:"required"`
	Description     string    `json:"description" api:"required"`
	// When one of the account's API keys last authenticated or one of its rules last
	// issued a token. Recorded at most once every 15 minutes, so it can lag. Null if
	// the account has never been used.
	LastUsedAt time.Time `json:"last_used_at" api:"required" format:"date-time"`
	Name       string    `json:"name" api:"required"`
	OrgID      string    `json:"org_id" api:"required"`
	// Authority revision. Cosmetic updates do not widen permissions.
	SecurityVersion int64     `json:"security_version" api:"required"`
	UpdatedAt       time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		ArchivedAt      respjson.Field
		CreatedAt       respjson.Field
		CreatedByUserID respjson.Field
		Description     respjson.Field
		LastUsedAt      respjson.Field
		Name            respjson.Field
		OrgID           respjson.Field
		SecurityVersion respjson.Field
		UpdatedAt       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServiceAccount) RawJSON() string { return r.JSON.raw }
func (r *ServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Name is required.
type ServiceAccountConfigParam struct {
	Name        string            `json:"name" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r ServiceAccountConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ServiceAccountConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServiceAccountConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ServiceAccountNewParams struct {
	ServiceAccountConfig ServiceAccountConfigParam
	paramObj
}

func (r ServiceAccountNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.ServiceAccountConfig)
}
func (r *ServiceAccountNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ServiceAccountUpdateParams struct {
	ExpectedVersion int64             `json:"expected_version" api:"required"`
	Description     param.Opt[string] `json:"description,omitzero"`
	Name            param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ServiceAccountUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ServiceAccountUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServiceAccountUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ServiceAccountListParams struct {
	Limit  param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against name or description; exact ID match.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	// Any of "active", "archived", "all".
	Status ServiceAccountListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ServiceAccountListParams]'s query parameters as
// `url.Values`.
func (r ServiceAccountListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ServiceAccountListParamsStatus string

const (
	ServiceAccountListParamsStatusActive   ServiceAccountListParamsStatus = "active"
	ServiceAccountListParamsStatusArchived ServiceAccountListParamsStatus = "archived"
	ServiceAccountListParamsStatusAll      ServiceAccountListParamsStatus = "all"
)
