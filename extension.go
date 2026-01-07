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

	"github.com/onkernel/kernel-go-sdk/internal/apiform"
	"github.com/onkernel/kernel-go-sdk/internal/apijson"
	"github.com/onkernel/kernel-go-sdk/internal/apiquery"
	"github.com/onkernel/kernel-go-sdk/internal/requestconfig"
	"github.com/onkernel/kernel-go-sdk/option"
	"github.com/onkernel/kernel-go-sdk/packages/param"
	"github.com/onkernel/kernel-go-sdk/packages/respjson"
)

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

// List extensions owned by the caller's organization.
func (r *ExtensionService) List(ctx context.Context, opts ...option.RequestOption) (res *[]ExtensionListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "extensions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Delete an extension by its ID or by its name.
func (r *ExtensionService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("extensions/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

// Download the extension as a ZIP archive by ID or name.
func (r *ExtensionService) Download(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return
	}
	path := fmt.Sprintf("extensions/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Returns a ZIP archive containing the unpacked extension fetched from the Chrome
// Web Store.
func (r *ExtensionService) DownloadFromChromeStore(ctx context.Context, query ExtensionDownloadFromChromeStoreParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	path := "extensions/from_chrome_store"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Upload a zip file containing an unpacked browser extension. Optionally provide a
// unique name for later reference.
func (r *ExtensionService) Upload(ctx context.Context, body ExtensionUploadParams, opts ...option.RequestOption) (res *ExtensionUploadResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "extensions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// A browser extension uploaded to Kernel.
type ExtensionListResponse struct {
	// Unique identifier for the extension
	ID string `json:"id,required"`
	// Timestamp when the extension was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Size of the extension archive in bytes
	SizeBytes int64 `json:"size_bytes,required"`
	// Timestamp when the extension was last used
	LastUsedAt time.Time `json:"last_used_at,nullable" format:"date-time"`
	// Optional, easier-to-reference name for the extension. Must be unique within the
	// organization.
	Name string `json:"name,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		SizeBytes   respjson.Field
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
type ExtensionUploadResponse struct {
	// Unique identifier for the extension
	ID string `json:"id,required"`
	// Timestamp when the extension was created
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// Size of the extension archive in bytes
	SizeBytes int64 `json:"size_bytes,required"`
	// Timestamp when the extension was last used
	LastUsedAt time.Time `json:"last_used_at,nullable" format:"date-time"`
	// Optional, easier-to-reference name for the extension. Must be unique within the
	// organization.
	Name string `json:"name,nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		SizeBytes   respjson.Field
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

type ExtensionDownloadFromChromeStoreParams struct {
	// Chrome Web Store URL for the extension.
	URL string `query:"url,required" json:"-"`
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
	File io.Reader `json:"file,omitzero,required" format:"binary"`
	// Optional unique name within the organization to reference this extension.
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
