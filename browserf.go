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
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// BrowserFService contains methods and other services that help with interacting
// with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBrowserFService] method instead.
type BrowserFService struct {
	Options []option.RequestOption
	Watch   BrowserFWatchService
}

// NewBrowserFService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBrowserFService(opts ...option.RequestOption) (r BrowserFService) {
	r = BrowserFService{}
	r.Options = opts
	r.Watch = NewBrowserFWatchService(opts...)
	return
}

// Create a new directory
func (r *BrowserFService) NewDirectory(ctx context.Context, id string, body BrowserFNewDirectoryParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/create_directory", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, nil, opts...)
	return
}

// Delete a directory
func (r *BrowserFService) DeleteDirectory(ctx context.Context, id string, body BrowserFDeleteDirectoryParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/delete_directory", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, nil, opts...)
	return
}

// Delete a file
func (r *BrowserFService) DeleteFile(ctx context.Context, id string, body BrowserFDeleteFileParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/delete_file", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, nil, opts...)
	return
}

// Returns a ZIP file containing the contents of the specified directory.
func (r *BrowserFService) DownloadDirZip(ctx context.Context, id string, query BrowserFDownloadDirZipParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/zip")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/download_dir_zip", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Get information about a file or directory
func (r *BrowserFService) FileInfo(ctx context.Context, id string, query BrowserFFileInfoParams, opts ...option.RequestOption) (res *BrowserFFileInfoResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/file_info", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// List files in a directory
func (r *BrowserFService) ListFiles(ctx context.Context, id string, query BrowserFListFilesParams, opts ...option.RequestOption) (res *[]BrowserFListFilesResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/list_files", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Move or rename a file or directory
func (r *BrowserFService) Move(ctx context.Context, id string, body BrowserFMoveParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/move", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, nil, opts...)
	return
}

// Read file contents
func (r *BrowserFService) ReadFile(ctx context.Context, id string, query BrowserFReadFileParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/read_file", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Set file or directory permissions/ownership
func (r *BrowserFService) SetFilePermissions(ctx context.Context, id string, body BrowserFSetFilePermissionsParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/set_file_permissions", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, nil, opts...)
	return
}

// Allows uploading single or multiple files to the remote filesystem.
func (r *BrowserFService) Upload(ctx context.Context, id string, body BrowserFUploadParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/upload", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// Upload a zip file and extract its contents to the specified destination path.
func (r *BrowserFService) UploadZip(ctx context.Context, id string, body BrowserFUploadZipParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/upload_zip", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// Write or create a file
func (r *BrowserFService) WriteFile(ctx context.Context, id string, contents io.Reader, params BrowserFWriteFileParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*"), option.WithRequestBody("application/octet-stream", contents)}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return
	}
	path := fmt.Sprintf("browsers/%s/fs/write_file", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, params, nil, opts...)
	return
}

type BrowserFFileInfoResponse struct {
	// Whether the path is a directory.
	IsDir bool `json:"is_dir,required"`
	// Last modification time.
	ModTime time.Time `json:"mod_time,required" format:"date-time"`
	// File mode bits (e.g., "drwxr-xr-x" or "-rw-r--r--").
	Mode string `json:"mode,required"`
	// Base name of the file or directory.
	Name string `json:"name,required"`
	// Absolute path.
	Path string `json:"path,required"`
	// Size in bytes. 0 for directories.
	SizeBytes int64 `json:"size_bytes,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		IsDir       respjson.Field
		ModTime     respjson.Field
		Mode        respjson.Field
		Name        respjson.Field
		Path        respjson.Field
		SizeBytes   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserFFileInfoResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserFFileInfoResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFListFilesResponse struct {
	// Whether the path is a directory.
	IsDir bool `json:"is_dir,required"`
	// Last modification time.
	ModTime time.Time `json:"mod_time,required" format:"date-time"`
	// File mode bits (e.g., "drwxr-xr-x" or "-rw-r--r--").
	Mode string `json:"mode,required"`
	// Base name of the file or directory.
	Name string `json:"name,required"`
	// Absolute path.
	Path string `json:"path,required"`
	// Size in bytes. 0 for directories.
	SizeBytes int64 `json:"size_bytes,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		IsDir       respjson.Field
		ModTime     respjson.Field
		Mode        respjson.Field
		Name        respjson.Field
		Path        respjson.Field
		SizeBytes   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BrowserFListFilesResponse) RawJSON() string { return r.JSON.raw }
func (r *BrowserFListFilesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFNewDirectoryParams struct {
	// Absolute directory path to create.
	Path string `json:"path,required"`
	// Optional directory mode (octal string, e.g. 755). Defaults to 755.
	Mode param.Opt[string] `json:"mode,omitzero"`
	paramObj
}

func (r BrowserFNewDirectoryParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFNewDirectoryParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFNewDirectoryParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFDeleteDirectoryParams struct {
	// Absolute path to delete.
	Path string `json:"path,required"`
	paramObj
}

func (r BrowserFDeleteDirectoryParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFDeleteDirectoryParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFDeleteDirectoryParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFDeleteFileParams struct {
	// Absolute path to delete.
	Path string `json:"path,required"`
	paramObj
}

func (r BrowserFDeleteFileParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFDeleteFileParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFDeleteFileParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFDownloadDirZipParams struct {
	// Absolute directory path to archive and download.
	Path string `query:"path,required" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserFDownloadDirZipParams]'s query parameters as
// `url.Values`.
func (r BrowserFDownloadDirZipParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserFFileInfoParams struct {
	// Absolute path of the file or directory.
	Path string `query:"path,required" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserFFileInfoParams]'s query parameters as `url.Values`.
func (r BrowserFFileInfoParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserFListFilesParams struct {
	// Absolute directory path.
	Path string `query:"path,required" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserFListFilesParams]'s query parameters as
// `url.Values`.
func (r BrowserFListFilesParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserFMoveParams struct {
	// Absolute destination path.
	DestPath string `json:"dest_path,required"`
	// Absolute source path.
	SrcPath string `json:"src_path,required"`
	paramObj
}

func (r BrowserFMoveParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFMoveParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFMoveParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFReadFileParams struct {
	// Absolute file path to read.
	Path string `query:"path,required" json:"-"`
	paramObj
}

// URLQuery serializes [BrowserFReadFileParams]'s query parameters as `url.Values`.
func (r BrowserFReadFileParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BrowserFSetFilePermissionsParams struct {
	// File mode bits (octal string, e.g. 644).
	Mode string `json:"mode,required"`
	// Absolute path whose permissions are to be changed.
	Path string `json:"path,required"`
	// New group name or GID.
	Group param.Opt[string] `json:"group,omitzero"`
	// New owner username or UID.
	Owner param.Opt[string] `json:"owner,omitzero"`
	paramObj
}

func (r BrowserFSetFilePermissionsParams) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFSetFilePermissionsParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFSetFilePermissionsParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFUploadParams struct {
	Files []BrowserFUploadParamsFile `json:"files,omitzero,required"`
	paramObj
}

func (r BrowserFUploadParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

// The properties DestPath, File are required.
type BrowserFUploadParamsFile struct {
	// Absolute destination path to write the file.
	DestPath string    `json:"dest_path,required"`
	File     io.Reader `json:"file,omitzero,required" format:"binary"`
	paramObj
}

func (r BrowserFUploadParamsFile) MarshalJSON() (data []byte, err error) {
	type shadow BrowserFUploadParamsFile
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserFUploadParamsFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BrowserFUploadZipParams struct {
	// Absolute destination directory to extract the archive to.
	DestPath string    `json:"dest_path,required"`
	ZipFile  io.Reader `json:"zip_file,omitzero,required" format:"binary"`
	paramObj
}

func (r BrowserFUploadZipParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type BrowserFWriteFileParams struct {
	// Destination absolute file path.
	Path string `query:"path,required" json:"-"`
	// Optional file mode (octal string, e.g. 644). Defaults to 644.
	Mode param.Opt[string] `query:"mode,omitzero" json:"-"`
	paramObj
}

func (r BrowserFWriteFileParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

// URLQuery serializes [BrowserFWriteFileParams]'s query parameters as
// `url.Values`.
func (r BrowserFWriteFileParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
