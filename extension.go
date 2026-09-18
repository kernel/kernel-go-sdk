// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apiform"
	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/apiquery"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/pagination"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Create, list, retrieve, and delete browser extensions.
//
// ExtensionService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewExtensionService] method instead.
type ExtensionService struct {
	Options []option.RequestOption
}

// NewExtensionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewExtensionService(opts ...option.RequestOption) (r ExtensionService) {
	r = ExtensionService{}
	r.Options = opts
	return
}

// List extensions in the resolved project.
func (r *ExtensionService) List(ctx context.Context, query ExtensionListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[ExtensionListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "extensions"
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

// List extensions in the resolved project.
func (r *ExtensionService) ListAutoPaging(ctx context.Context, query ExtensionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[ExtensionListResponse] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete an extension by its ID or by its name.
func (r *ExtensionService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("extensions/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Download the extension as a ZIP archive by ID or name.
func (r *ExtensionService) Download(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("extensions/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Returns a ZIP archive containing the unpacked extension fetched from the Chrome
// Web Store.
func (r *ExtensionService) DownloadFromChromeStore(ctx context.Context, query ExtensionDownloadFromChromeStoreParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	path := "extensions/from_chrome_store"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get an extension's metadata (name, size, timestamps) by ID or name, without
// downloading the archive.
func (r *ExtensionService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *ExtensionGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("extensions/%s/metadata", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Upload a zip file containing an unpacked browser extension. Optionally provide a
// unique name for later reference.
func (r *ExtensionService) Upload(ctx context.Context, body ExtensionUploadParams, opts ...option.RequestOption) (res *ExtensionUploadResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "extensions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// A browser extension uploaded to Kernel.
type ExtensionListResponse struct {
	// Unique identifier for the extension
	ID string `json:"id" api:"required"`
	// Timestamp when the extension was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Size of the extension archive in bytes
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// SHA-256 checksum, encoded as lowercase hexadecimal, of the exact uploaded
	// extension archive bytes. This is not a normalized checksum of the extension
	// contents; archive metadata, file ordering, and compression can change the
	// checksum for otherwise identical contents. Omitted for legacy rows and
	// server-repackaged Chrome Web Store extensions.
	Checksum string `json:"checksum" api:"nullable"`
	// Timestamp when the extension was last used
	LastUsedAt time.Time `json:"last_used_at" api:"nullable" format:"date-time"`
	// Optional, easier-to-reference name for the extension. Must be unique within the
	// project.
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		SizeBytes   respjson.Field
		Checksum    respjson.Field
		LastUsedAt  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExtensionListResponse) RawJSON() string { return r.JSON.raw }
func (r *ExtensionListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser extension uploaded to Kernel.
type ExtensionGetResponse struct {
	// Unique identifier for the extension
	ID string `json:"id" api:"required"`
	// Timestamp when the extension was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Size of the extension archive in bytes
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// SHA-256 checksum, encoded as lowercase hexadecimal, of the exact uploaded
	// extension archive bytes. This is not a normalized checksum of the extension
	// contents; archive metadata, file ordering, and compression can change the
	// checksum for otherwise identical contents. Omitted for legacy rows and
	// server-repackaged Chrome Web Store extensions.
	Checksum string `json:"checksum" api:"nullable"`
	// Timestamp when the extension was last used
	LastUsedAt time.Time `json:"last_used_at" api:"nullable" format:"date-time"`
	// Optional, easier-to-reference name for the extension. Must be unique within the
	// project.
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		SizeBytes   respjson.Field
		Checksum    respjson.Field
		LastUsedAt  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExtensionGetResponse) RawJSON() string { return r.JSON.raw }
func (r *ExtensionGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A browser extension uploaded to Kernel.
type ExtensionUploadResponse struct {
	// Unique identifier for the extension
	ID string `json:"id" api:"required"`
	// Timestamp when the extension was created
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Size of the extension archive in bytes
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// SHA-256 checksum, encoded as lowercase hexadecimal, of the exact uploaded
	// extension archive bytes. This is not a normalized checksum of the extension
	// contents; archive metadata, file ordering, and compression can change the
	// checksum for otherwise identical contents. Omitted for legacy rows and
	// server-repackaged Chrome Web Store extensions.
	Checksum string `json:"checksum" api:"nullable"`
	// Timestamp when the extension was last used
	LastUsedAt time.Time `json:"last_used_at" api:"nullable" format:"date-time"`
	// Optional, easier-to-reference name for the extension. Must be unique within the
	// project.
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		SizeBytes   respjson.Field
		Checksum    respjson.Field
		LastUsedAt  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExtensionUploadResponse) RawJSON() string { return r.JSON.raw }
func (r *ExtensionUploadResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExtensionListParams struct {
	// Limit the number of extensions to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on extension name using the database collation. In
	// production, matching is case- and accent-insensitive. During the default-project
	// migration, unscoped requests prefer a concrete default-project extension over a
	// legacy unscoped extension with the same name.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Offset the number of extensions to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against extension name. IDs match by exact
	// value.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ExtensionListParams]'s query parameters as `url.Values`.
func (r ExtensionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ExtensionDownloadFromChromeStoreParams struct {
	// Chrome Web Store URL for the extension.
	URL string `query:"url" api:"required" json:"-"`
	// Target operating system for the extension package. Defaults to linux.
	//
	// Any of "win", "mac", "linux".
	Os ExtensionDownloadFromChromeStoreParamsOs `query:"os,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ExtensionDownloadFromChromeStoreParams]'s query parameters
// as `url.Values`.
func (r ExtensionDownloadFromChromeStoreParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Target operating system for the extension package. Defaults to linux.
type ExtensionDownloadFromChromeStoreParamsOs string

const (
	ExtensionDownloadFromChromeStoreParamsOsWin   ExtensionDownloadFromChromeStoreParamsOs = "win"
	ExtensionDownloadFromChromeStoreParamsOsMac   ExtensionDownloadFromChromeStoreParamsOs = "mac"
	ExtensionDownloadFromChromeStoreParamsOsLinux ExtensionDownloadFromChromeStoreParamsOs = "linux"
)

type ExtensionUploadParams struct {
	// ZIP file containing the browser extension.
	File io.Reader `json:"file,omitzero" api:"required" format:"binary"`
	// Optional unique name within the project to reference this extension.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r ExtensionUploadParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}
