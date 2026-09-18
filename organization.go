// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"github.com/kernel/kernel-go-sdk/option"
)

// OrganizationService contains methods and other services that help with
// interacting with the kernel API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationService] method instead.
type OrganizationService struct {
	Options []option.RequestOption
	// Read and manage organization-level limits.
	Entitlements OrganizationEntitlementService
	// Read and manage organization-level limits.
	Limits OrganizationLimitService
}

// NewOrganizationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewOrganizationService(opts ...option.RequestOption) (r OrganizationService) {
	r = OrganizationService{}
	r.Options = opts
	r.Entitlements = NewOrganizationEntitlementService(opts...)
	r.Limits = NewOrganizationLimitService(opts...)
	return
}
