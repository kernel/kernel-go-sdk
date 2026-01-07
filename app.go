// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/apiquery"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/pagination"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
	"github.com/onkernel/kernel-go-sdk/shared"
	"github.com/onkernel/kernel-go-sdk/shared/constant"
)

// AppService contains methods and other services that help with interacting with
// the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAppService] method instead.
type AppService struct {
	Options []option.RequestOption
}

// NewAppService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAppService(opts ...option.RequestOption) (r AppService) {
	r = AppService{}
	r.Options = opts
	return
}

// List applications. Optionally filter by app name and/or version label.
func (r *AppService) List(ctx context.Context, query AppListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[AppListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "apps"
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

// List applications. Optionally filter by app name and/or version label.
func (r *AppService) ListAutoPaging(ctx context.Context, query AppListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[AppListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Summary of an application version.
type AppListResponse struct {
	// Unique identifier for the app version
	ID string `json:"id,required"`
	// List of actions available on the app
	Actions []shared.AppAction `json:"actions,required"`
	// Name of the application
	AppName string `json:"app_name,required"`
	// Deployment ID
	Deployment string `json:"deployment,required"`
	// Environment variables configured for this app version
	EnvVars map[string]string `json:"env_vars,required"`
	// Deployment region code
	Region constant.AwsUsEast1a `json:"region,required"`
	// Version label for the application
	Version string `json:"version,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Actions     respjson.Field
		AppName     respjson.Field
		Deployment  respjson.Field
		EnvVars     respjson.Field
		Region      respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AppListResponse) RawJSON() string { return r.JSON.raw }
func (r *AppListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AppListParams struct {
	// Filter results by application name.
	AppName param.Opt[string] `query:"app_name,omitzero" json:"-"`
	// Limit the number of apps to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of apps to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Filter results by version label.
	Version param.Opt[string] `query:"version,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AppListParams]'s query parameters as `url.Values`.
func (r AppListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
