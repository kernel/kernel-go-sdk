// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
)

// Create, list, retrieve, and delete browser profiles.
//
// ProfileService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProfileService] method instead.
type ProfileService struct {
	Options []option.RequestOption
}

// NewProfileService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewProfileService(opts ...option.RequestOption) (r ProfileService) {
	r = ProfileService{}
	r.Options = opts
	return
}

// Create a browser profile that can be used to load state into future browser
// sessions.
func (r *ProfileService) New(ctx context.Context, body ProfileNewParams, opts ...option.RequestOption) (res *Profile, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "profiles"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve details for a single profile by its ID or name.
func (r *ProfileService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *Profile, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("profiles/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List profiles with optional filtering and pagination.
func (r *ProfileService) List(ctx context.Context, query ProfileListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[Profile], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "profiles"
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

// List profiles with optional filtering and pagination.
func (r *ProfileService) ListAutoPaging(ctx context.Context, query ProfileListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[Profile] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a profile by its ID or by its name.
func (r *ProfileService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("profiles/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Download the profile. Profiles are JSON files containing the pieces of state
// that we save.
func (r *ProfileService) Download(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("profiles/%s/download", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type ProfileNewParams struct {
	// Optional name of the profile. Must be unique within the organization.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ProfileNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ProfileNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProfileNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProfileListParams struct {
	// Limit the number of profiles to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Offset the number of profiles to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Search profiles by name or ID.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ProfileListParams]'s query parameters as `url.Values`.
func (r ProfileListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
