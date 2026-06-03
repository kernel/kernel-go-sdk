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

// Create and manage API keys for organization and project-scoped access.
//
// APIKeyService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAPIKeyService] method instead.
type APIKeyService struct {
	Options []option.RequestOption
}

// NewAPIKeyService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAPIKeyService(opts ...option.RequestOption) (r APIKeyService) {
	r = APIKeyService{}
	r.Options = opts
	return
}

// Create a new API key within the authenticated organization.
func (r *APIKeyService) New(ctx context.Context, body APIKeyNewParams, opts ...option.RequestOption) (res *CreatedAPIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/api_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an API key by ID for the authenticated organization. API keys are
// masked.
func (r *APIKeyService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an API key's name.
func (r *APIKeyService) Update(ctx context.Context, id string, body APIKeyUpdateParams, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List API keys for the authenticated organization. API keys are masked.
func (r *APIKeyService) List(ctx context.Context, query APIKeyListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[APIKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "org/api_keys"
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

// List API keys for the authenticated organization. API keys are masked.
func (r *APIKeyService) ListAutoPaging(ctx context.Context, query APIKeyListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[APIKey] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete an API key.
func (r *APIKeyService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return err
	}
	path := fmt.Sprintf("org/api_keys/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type APIKey struct {
	// Unique API key identifier
	ID string `json:"id" api:"required"`
	// When the API key was created
	CreatedAt time.Time       `json:"created_at" api:"required" format:"date-time"`
	CreatedBy APIKeyCreatedBy `json:"created_by" api:"required"`
	// When the API key expires
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Masked version of the API key
	MaskedKey string `json:"masked_key" api:"required"`
	// API key name
	Name string `json:"name" api:"required"`
	// Project identifier for project-scoped API keys. Null means org-wide.
	ProjectID string `json:"project_id" api:"required"`
	// Project name for project-scoped API keys. Null means the key is org-wide or the
	// project name is unavailable.
	ProjectName string `json:"project_name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		CreatedBy   respjson.Field
		ExpiresAt   respjson.Field
		MaskedKey   respjson.Field
		Name        respjson.Field
		ProjectID   respjson.Field
		ProjectName respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKey) RawJSON() string { return r.JSON.raw }
func (r *APIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyCreatedBy struct {
	// Kernel user ID of the creator.
	ID string `json:"id" api:"required"`
	// Email address of the creator.
	Email string `json:"email" api:"required" format:"email"`
	// Display name of the creator, if available.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyCreatedBy) RawJSON() string { return r.JSON.raw }
func (r *APIKeyCreatedBy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// API key returned immediately after creation. Includes the plaintext key once.
type CreatedAPIKey struct {
	// Plaintext API key. Only returned once when the key is created.
	Key string `json:"key" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Key         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	APIKey
}

// Returns the unmodified JSON received from the API
func (r CreatedAPIKey) RawJSON() string { return r.JSON.raw }
func (r *CreatedAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyNewParams struct {
	// API key name (1-255 characters)
	Name string `json:"name" api:"required"`
	// Number of days until expiry, up to 3650. Use null for never.
	DaysToExpire param.Opt[int64] `json:"days_to_expire,omitzero"`
	// Unique project identifier
	ProjectID param.Opt[string] `json:"project_id,omitzero"`
	paramObj
}

func (r APIKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyUpdateParams struct {
	// New API key name
	Name string `json:"name" api:"required"`
	paramObj
}

func (r APIKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow APIKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *APIKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyListParams struct {
	// Maximum number of results to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Number of results to skip
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against API key name, creator, and project. API
	// key identifiers and masked keys match by exact value or prefix.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	// Field to sort API keys by.
	//
	// Any of "created_at", "name", "expires_at".
	SortBy APIKeyListParamsSortBy `query:"sort_by,omitzero" json:"-"`
	// Sort direction for API keys.
	//
	// Any of "asc", "desc".
	SortDirection APIKeyListParamsSortDirection `query:"sort_direction,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [APIKeyListParams]'s query parameters as `url.Values`.
func (r APIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Field to sort API keys by.
type APIKeyListParamsSortBy string

const (
	APIKeyListParamsSortByCreatedAt APIKeyListParamsSortBy = "created_at"
	APIKeyListParamsSortByName      APIKeyListParamsSortBy = "name"
	APIKeyListParamsSortByExpiresAt APIKeyListParamsSortBy = "expires_at"
)

// Sort direction for API keys.
type APIKeyListParamsSortDirection string

const (
	APIKeyListParamsSortDirectionAsc  APIKeyListParamsSortDirection = "asc"
	APIKeyListParamsSortDirectionDesc APIKeyListParamsSortDirection = "desc"
)
