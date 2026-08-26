// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
)

// Resolve browser and proxy recommendations for bot-protected sites.
//
// ConfigRegistryAnalysisService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewConfigRegistryAnalysisService] method instead.
type ConfigRegistryAnalysisService struct {
	Options []option.RequestOption
}

// NewConfigRegistryAnalysisService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewConfigRegistryAnalysisService(opts ...option.RequestOption) (r ConfigRegistryAnalysisService) {
	r = ConfigRegistryAnalysisService{}
	r.Options = opts
	return
}

// Returns a project-scoped historical analysis and the recommendation outcome
// concluded by that run. Later knowledge does not change this response.
func (r *ConfigRegistryAnalysisService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ConfigRegistryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("config-registry/analyses/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists analyses for the selected project, newest first.
func (r *ConfigRegistryAnalysisService) List(ctx context.Context, query ConfigRegistryAnalysisListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[AnalysisSummary], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "config-registry/analyses"
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

// Lists analyses for the selected project, newest first.
func (r *ConfigRegistryAnalysisService) ListAutoPaging(ctx context.Context, query ConfigRegistryAnalysisListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[AnalysisSummary] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

type ConfigRegistryAnalysisListParams struct {
	Limit  param.Opt[int64] `query:"limit,omitzero" json:"-"`
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring search over requested URLs.
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ConfigRegistryAnalysisListParams]'s query parameters as
// `url.Values`.
func (r ConfigRegistryAnalysisListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
