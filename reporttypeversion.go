// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Discover report types and submit typed reports against tasks.
//
// ReportTypeVersionService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewReportTypeVersionService] method instead.
type ReportTypeVersionService struct {
	Options []option.RequestOption
}

// NewReportTypeVersionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewReportTypeVersionService(opts ...option.RequestOption) (r ReportTypeVersionService) {
	r = ReportTypeVersionService{}
	r.Options = opts
	return
}

// Retrieve one schema version of a report type, including retired versions.
func (r *ReportTypeVersionService) Get(ctx context.Context, version int64, query ReportTypeVersionGetParams, opts ...option.RequestOption) (res *ReportTypeVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if query.Type == "" {
		err = errors.New("missing required type parameter")
		return nil, err
	}
	path := fmt.Sprintf("report_types/%s/versions/%v", query.Type, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// One immutable schema version of a report type.
type ReportTypeVersion struct {
	// Whether the version is retired. Retired versions are no longer accepted for new
	// reports; existing reports keep them.
	Retired bool `json:"retired" api:"required"`
	// JSON Schema (draft 2020-12) that a report's data must match. Its field
	// descriptions and examples describe what to send.
	Schema map[string]any `json:"schema" api:"required"`
	// Report type name, such as task.result
	Type string `json:"type" api:"required"`
	// Schema version number, starting at 1.
	Version int64 `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Retired     respjson.Field
		Schema      respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ReportTypeVersion) RawJSON() string { return r.JSON.raw }
func (r *ReportTypeVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ReportTypeVersionGetParams struct {
	Type string `path:"type" api:"required" json:"-"`
	paramObj
}
