// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/kernel/kernel-go-sdk/internal/apijson"
	"github.com/kernel/kernel-go-sdk/internal/requestconfig"
	"github.com/kernel/kernel-go-sdk/option"
	"github.com/kernel/kernel-go-sdk/packages/respjson"
)

// Read and manage organization-level limits.
//
// OrganizationEntitlementService contains methods and other services that help
// with interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationEntitlementService] method instead.
type OrganizationEntitlementService struct {
	Options []option.RequestOption
}

// NewOrganizationEntitlementService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationEntitlementService(opts ...option.RequestOption) (r OrganizationEntitlementService) {
	r = OrganizationEntitlementService{}
	r.Options = opts
	return
}

// Get the authenticated organization's effective feature access and constraints
// after applying its plan, active trial treatment, plan status, and
// organization-specific overrides. Null constraint values mean unlimited.
func (r *OrganizationEntitlementService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrgEntitlements, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "org/entitlements"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Effective feature access and constraints for the authenticated organization.
// Values already include trial treatment, plan status, and organization-specific
// overrides; consumers should use these resolved values instead of comparing plan
// IDs.
type OrgEntitlements struct {
	Features OrgEntitlementsFeatures `json:"features" api:"required"`
	Limits   OrgEntitlementsLimits   `json:"limits" api:"required"`
	Plan     OrgEntitlementsPlan     `json:"plan" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Features    respjson.Field
		Limits      respjson.Field
		Plan        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlements) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlements) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeatures struct {
	BrowserExtensions   OrgEntitlementsFeaturesBrowserExtensions   `json:"browser_extensions" api:"required"`
	BrowserPools        OrgEntitlementsFeaturesBrowserPools        `json:"browser_pools" api:"required"`
	BrowserReplays      OrgEntitlementsFeaturesBrowserReplays      `json:"browser_replays" api:"required"`
	CredentialProviders OrgEntitlementsFeaturesCredentialProviders `json:"credential_providers" api:"required"`
	Credentials         OrgEntitlementsFeaturesCredentials         `json:"credentials" api:"required"`
	CustomProxies       OrgEntitlementsFeaturesCustomProxies       `json:"custom_proxies" api:"required"`
	FileIo              OrgEntitlementsFeaturesFileIo              `json:"file_io" api:"required"`
	GPU                 OrgEntitlementsFeaturesGPU                 `json:"gpu" api:"required"`
	ManagedAuth         OrgEntitlementsFeaturesManagedAuth         `json:"managed_auth" api:"required"`
	ManagedProxies      OrgEntitlementsFeaturesManagedProxies      `json:"managed_proxies" api:"required"`
	Profiles            OrgEntitlementsFeaturesProfiles            `json:"profiles" api:"required"`
	ProxyBypassHosts    OrgEntitlementsFeaturesProxyBypassHosts    `json:"proxy_bypass_hosts" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BrowserExtensions   respjson.Field
		BrowserPools        respjson.Field
		BrowserReplays      respjson.Field
		CredentialProviders respjson.Field
		Credentials         respjson.Field
		CustomProxies       respjson.Field
		FileIo              respjson.Field
		GPU                 respjson.Field
		ManagedAuth         respjson.Field
		ManagedProxies      respjson.Field
		Profiles            respjson.Field
		ProxyBypassHosts    respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeatures) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeatures) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesBrowserExtensions struct {
	// Whether browser extensions are available.
	Enabled bool `json:"enabled" api:"required"`
	// Maximum active custom extensions the organization may store. Null means
	// unlimited. Loading stored extensions into a browser is not plan-limited.
	MaxStoredPerOrg int64 `json:"max_stored_per_org" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled         respjson.Field
		MaxStoredPerOrg respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesBrowserExtensions) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesBrowserExtensions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesBrowserPools struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesBrowserPools) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesBrowserPools) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesBrowserReplays struct {
	// Whether browser replay recording is available.
	Enabled bool `json:"enabled" api:"required"`
	// Number of days browser replays are retained, matching the replay reaper policy.
	RetentionDays int64 `json:"retention_days" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled       respjson.Field
		RetentionDays respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesBrowserReplays) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesBrowserReplays) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesCredentialProviders struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesCredentialProviders) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesCredentialProviders) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesCredentials struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesCredentials) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesCredentials) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesCustomProxies struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesCustomProxies) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesCustomProxies) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesFileIo struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesFileIo) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesFileIo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesGPU struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesGPU) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesGPU) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesManagedAuth struct {
	// Whether managed auth is available.
	Enabled bool `json:"enabled" api:"required"`
	// Effective interval in seconds used when a connection is created without an
	// explicit health-check interval.
	HealthCheckIntervalDefaultSeconds int64 `json:"health_check_interval_default_seconds" api:"required"`
	// Maximum accepted managed auth health-check interval in seconds.
	HealthCheckIntervalMaxSeconds int64 `json:"health_check_interval_max_seconds" api:"required"`
	// Minimum accepted managed auth health-check interval in seconds.
	HealthCheckIntervalMinSeconds int64 `json:"health_check_interval_min_seconds" api:"required"`
	// Maximum active managed auth connections in the organization. Null means
	// unlimited.
	MaxConnections int64 `json:"max_connections" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled                           respjson.Field
		HealthCheckIntervalDefaultSeconds respjson.Field
		HealthCheckIntervalMaxSeconds     respjson.Field
		HealthCheckIntervalMinSeconds     respjson.Field
		MaxConnections                    respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesManagedAuth) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesManagedAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesManagedProxies struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesManagedProxies) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesManagedProxies) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesProfiles struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesProfiles) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesProfiles) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsFeaturesProxyBypassHosts struct {
	// Whether the organization is entitled to use this feature.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsFeaturesProxyBypassHosts) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsFeaturesProxyBypassHosts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsLimits struct {
	// Effective org-level default concurrent invocation ceiling for apps without an
	// app-specific override. App-specific overrides are not represented here.
	DefaultMaxConcurrentInvocationsPerApp int64 `json:"default_max_concurrent_invocations_per_app" api:"required"`
	// Effective organization-wide ceiling shared by on-demand browsers and browser
	// pool reservations.
	MaxConcurrentBrowsers int64 `json:"max_concurrent_browsers" api:"required"`
	// Effective organization-wide concurrent app invocation ceiling.
	MaxConcurrentInvocations int64 `json:"max_concurrent_invocations" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DefaultMaxConcurrentInvocationsPerApp respjson.Field
		MaxConcurrentBrowsers                 respjson.Field
		MaxConcurrentInvocations              respjson.Field
		ExtraFields                           map[string]respjson.Field
		raw                                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsLimits) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsLimits) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrgEntitlementsPlan struct {
	// The organization's contractual plan identifier. Use the resolved feature and
	// limit values, not this field, for access decisions.
	//
	// Any of "FREE", "HOBBYIST", "START_UP", "ENTERPRISE".
	ID string `json:"id" api:"required"`
	// The plan used to resolve plan-based access. Active trials resolve to START_UP
	// regardless of the contractual plan.
	//
	// Any of "FREE", "HOBBYIST", "START_UP", "ENTERPRISE".
	EffectiveID string `json:"effective_id" api:"required"`
	// Whether the organization is currently within its trial period.
	IsTrialing bool `json:"is_trialing" api:"required"`
	// Current billing status of the contractual plan, or null when no billing status
	// is recorded. Status-sensitive feature values already account for it.
	//
	// Any of "NEEDS_PAYMENT_METHOD", "ACTIVE", "CANCELED", "UNPAID".
	Status string `json:"status" api:"required"`
	// Configured trial end timestamp, or null when the organization has no trial. A
	// past timestamp may be returned when is_trialing is false.
	TrialEndsAt time.Time `json:"trial_ends_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		EffectiveID respjson.Field
		IsTrialing  respjson.Field
		Status      respjson.Field
		TrialEndsAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgEntitlementsPlan) RawJSON() string { return r.JSON.raw }
func (r *OrgEntitlementsPlan) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
