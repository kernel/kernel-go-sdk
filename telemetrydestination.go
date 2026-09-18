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

// Stream live telemetry events from a browser session, and manage the destinations
// sessions export them to.
//
// TelemetryDestinationService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTelemetryDestinationService] method instead.
type TelemetryDestinationService struct {
	Options []option.RequestOption
}

// NewTelemetryDestinationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewTelemetryDestinationService(opts ...option.RequestOption) (r TelemetryDestinationService) {
	r = TelemetryDestinationService{}
	r.Options = opts
	return
}

// Create an OTLP export destination in the authenticated organization. Names must
// be unique within the organization. Requires an organization-scoped credential or
// dashboard authentication; project-scoped credentials receive a 403.
func (r *TelemetryDestinationService) New(ctx context.Context, body TelemetryDestinationNewParams, opts ...option.RequestOption) (res *OtlpDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "telemetry/destinations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a customer-visible OTLP destination in the authenticated organization
// by its ID or name. Project-scoped credentials can retrieve these destinations
// for selection by workloads in their project. Non-dashboard reads return header
// values redacted.
func (r *TelemetryDestinationService) Get(ctx context.Context, idOrName string, opts ...option.RequestOption) (res *OtlpDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("telemetry/destinations/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an OTLP destination. Sessions already exporting to it pick up the new
// values without restarting, which makes this the way to rotate credentials
// without interrupting export.
//
// Names must be unique within the organization. Renaming is refused with a 409
// while a managed auth connection selects this destination by name, since that
// connection resolves the name on every login. Every other field, including
// `headers`, stays editable. Requires an organization-scoped credential or
// dashboard authentication; project-scoped credentials receive a 403.
func (r *TelemetryDestinationService) Update(ctx context.Context, idOrName string, body TelemetryDestinationUpdateParams, opts ...option.RequestOption) (res *OtlpDestination, err error) {
	opts = slices.Concat(r.Options, opts)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return nil, err
	}
	path := fmt.Sprintf("telemetry/destinations/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List customer-visible OTLP export destinations in the authenticated
// organization. Project-scoped credentials can list these destinations for
// selection by workloads in their project. Non-dashboard reads return header
// values redacted.
func (r *TelemetryDestinationService) List(ctx context.Context, query TelemetryDestinationListParams, opts ...option.RequestOption) (res *pagination.OffsetPagination[OtlpDestination], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "telemetry/destinations"
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

// List customer-visible OTLP export destinations in the authenticated
// organization. Project-scoped credentials can list these destinations for
// selection by workloads in their project. Non-dashboard reads return header
// values redacted.
func (r *TelemetryDestinationService) ListAutoPaging(ctx context.Context, query TelemetryDestinationListParams, opts ...option.RequestOption) *pagination.OffsetPaginationAutoPager[OtlpDestination] {
	return pagination.NewOffsetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete an OTLP destination. Sessions bound to it are still exporting, so the
// delete is refused with a 409 while any exist; either wait for those sessions to
// end or delete them first. It is refused the same way while a managed auth
// connection still selects it, because that connection re-resolves the destination
// on every login, and while a managed auth login using it is still in progress.
// Requires an organization-scoped credential or dashboard authentication;
// project-scoped credentials receive a 403.
func (r *TelemetryDestinationService) Delete(ctx context.Context, idOrName string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if idOrName == "" {
		err = errors.New("missing required id_or_name parameter")
		return err
	}
	path := fmt.Sprintf("telemetry/destinations/%s", idOrName)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// An OTLP endpoint to export browser session telemetry to. Reference one from
// `telemetry.export.otlp.destination` when creating a browser to export that
// session's captured telemetry to it.
type OtlpDestination struct {
	ID string `json:"id" api:"required"`
	// Failed deliveries since the last success, as observed by the relay process that
	// wrote the latest outcome. Zero means the most recently recorded outcome
	// succeeded.
	ConsecutiveFailures int64     `json:"consecutive_failures" api:"required"`
	CreatedAt           time.Time `json:"created_at" api:"required" format:"date-time"`
	// OTLP/HTTP endpoint telemetry is sent to.
	Endpoint string `json:"endpoint" api:"required"`
	// Headers sent with each export request. Names are returned in canonical form
	// (`Authorization`, not `authorization`). Non-dashboard reads return values
	// redacted as empty strings, so the keys are visible but the credentials are not.
	// Dashboard reads return the stored values.
	Headers map[string]string `json:"headers" api:"required"`
	// Unique within the organization. Usable in place of the ID when selecting a
	// destination, so it cannot be shaped like an ID.
	Name        string    `json:"name" api:"required"`
	UpdatedAt   time.Time `json:"updated_at" api:"required" format:"date-time"`
	Description string    `json:"description"`
	// Sanitized class of the delivery failure recorded at `last_error_at`. It is
	// retained after a later success, so its presence does not mean the destination is
	// currently failing. Response bodies, endpoint URLs, credentials, and raw
	// transport errors are never returned.
	LastError string `json:"last_error"`
	// Timestamp of the most recent failed delivery. It is retained after a later
	// success, so it can predate `last_export_at`. Read `consecutive_failures` to tell
	// whether the destination is currently failing.
	LastErrorAt time.Time `json:"last_error_at" format:"date-time"`
	// Timestamp of the most recent successful delivery. Moves only on success.
	LastExportAt time.Time `json:"last_export_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		ConsecutiveFailures respjson.Field
		CreatedAt           respjson.Field
		Endpoint            respjson.Field
		Headers             respjson.Field
		Name                respjson.Field
		UpdatedAt           respjson.Field
		Description         respjson.Field
		LastError           respjson.Field
		LastErrorAt         respjson.Field
		LastExportAt        respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OtlpDestination) RawJSON() string { return r.JSON.raw }
func (r *OtlpDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TelemetryDestinationNewParams struct {
	// Base endpoint of the OTLP/HTTP collector, without a signal path. Kernel appends
	// the signal path itself, so pass `https://api.honeycomb.io` rather than
	// `https://api.honeycomb.io/v1/logs`. If your provider's docs give you a
	// signal-specific URL, drop the trailing `/v1/logs`, `/v1/traces`, or
	// `/v1/metrics` — an endpoint that already carries one is rejected.
	//
	// Must be http or https, must resolve to a public address, and must carry no query
	// string or fragment. Examples: `https://api.honeycomb.io`,
	// `https://otlp-gateway-prod-us-east-0.grafana.net/otlp`,
	// `https://otlp.datadoghq.com` (Datadog's OTLP intake for US1, not its logs
	// intake).
	Endpoint string `json:"endpoint" api:"required"`
	// Unique within the organization.
	Name        string            `json:"name" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	// Headers sent with each export request, typically an ingestion key. Encrypted at
	// rest and returned redacted. Names and values must be valid HTTP header tokens,
	// and the names and values together cannot exceed 8192 bytes. Names are matched
	// case-insensitively and stored canonicalized, so supplying two spellings of one
	// header is rejected.
	Headers map[string]string `json:"headers,omitzero"`
	paramObj
}

func (r TelemetryDestinationNewParams) MarshalJSON() (data []byte, err error) {
	type shadow TelemetryDestinationNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TelemetryDestinationNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TelemetryDestinationUpdateParams struct {
	Description param.Opt[string] `json:"description,omitzero"`
	// Base endpoint of the OTLP/HTTP collector, without a signal path. Same rules as
	// on create.
	Endpoint param.Opt[string] `json:"endpoint,omitzero"`
	Name     param.Opt[string] `json:"name,omitzero"`
	// Edits stored headers key by key rather than replacing the map. A string value
	// adds or replaces that header, `null` deletes it, and any key you omit is left as
	// it is. Names are matched case-insensitively, so `authorization` replaces a
	// stored `Authorization` rather than adding a second entry. This is the credential
	// rotation path; sessions already exporting pick up the new values without
	// restarting. Names and values must be valid HTTP header tokens, and the names and
	// values together cannot exceed 8192 bytes.
	Headers map[string]string `json:"headers,omitzero"`
	paramObj
}

func (r TelemetryDestinationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow TelemetryDestinationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TelemetryDestinationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TelemetryDestinationListParams struct {
	// Limit the number of destinations to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Exact-match filter on destination name using the database collation. In
	// production, matching is case- and accent-insensitive.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Offset the number of destinations to return.
	Offset param.Opt[int64] `query:"offset,omitzero" json:"-"`
	// Case-insensitive substring match against destination name or endpoint. IDs match
	// by exact value.
	Query param.Opt[string] `query:"query,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [TelemetryDestinationListParams]'s query parameters as
// `url.Values`.
func (r TelemetryDestinationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
