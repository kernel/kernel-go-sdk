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
	return res, err
}

// Retrieve details for a single profile by its ID or name.
func (r *ProfileService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *Profile, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("profiles/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a profile's name. Names must be unique within the logical project; during
// the default-project migration, unscoped profiles and profiles in the org default
// project are treated as the same project. Duplicate-name conflicts are checked
// before update but are best-effort because there is no backing unique index.
// Renaming a profile while a browser session references it by name may prevent
// that session's changes from saving; prefer renaming when the profile is not in
// use.
func (r *ProfileService) Update(ctx context.Context, idOrName string, body ProfileUpdateParams, opts ...option.RequestOption) (res *Profile, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("profiles/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
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
		return err
	}
	path := fmt.Sprintf("profiles/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Downloads the profile in its stored format by default. Current profiles are
// returned as zstd-compressed tar archives, while legacy profiles remain JSON. Set
// `format=tar` to decompress current profiles during download; legacy profiles
// remain JSON.
func (r *ProfileService) Download(ctx context.Context, idOrName string, query ProfileDownloadParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("profiles/%s/download", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type ProfileNewParams struct {
	// Optional name of the profile. Must be unique within the logical project; during
	// the default-project migration, unscoped profiles and profiles in the org default
	// project are treated as the same project.
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

type ProfileUpdateParams struct {
	// New profile name. Must be unique within the logical project; during the
	// default-project migration, unscoped profiles and profiles in the org default
	// project are treated as the same project.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r ProfileUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ProfileUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ProfileUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProfileListParams struct {
	// Limit the number of profiles to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on profile name using the database collation. In production,
	// matching is case- and accent-insensitive. During the default-project migration,
	// unscoped requests prefer a concrete default-project profile over a legacy
	// unscoped profile with the same name.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Offset the number of profiles to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against profile name or ID.
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

type ProfileDownloadParams struct {
	// Response format for current profile archives. Legacy profiles are always
	// returned as JSON.
	//
	// Any of "tar.zst", "tar".
	Format ProfileDownloadParamsFormat `query:"format,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ProfileDownloadParams]'s query parameters as `url.Values`.
func (r ProfileDownloadParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Response format for current profile archives. Legacy profiles are always
// returned as JSON.
type ProfileDownloadParamsFormat string

const (
	ProfileDownloadParamsFormatTarZst ProfileDownloadParamsFormat = "tar.zst"
	ProfileDownloadParamsFormatTar    ProfileDownloadParamsFormat = "tar"
)
