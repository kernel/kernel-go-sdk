// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel

import (
	"github.com/kernel/kernel-go-sdk/internal/apierror"
	"github.com/kernel/kernel-go-sdk/packages/param"
	"github.com/kernel/kernel-go-sdk/shared"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type Error = apierror.Error

// An action available on the app
//
// This is an alias to an internal type.
type AppAction = shared.AppAction

// Extension selection for the browser session. Provide either id or name of an
// extension uploaded to Kernel.
//
// This is an alias to an internal type.
type BrowserExtension = shared.BrowserExtension

// Extension selection for the browser session. Provide either id or name of an
// extension uploaded to Kernel.
//
// This is an alias to an internal type.
type BrowserExtensionParam = shared.BrowserExtensionParam

// Profile selection for the browser session. Provide either id or name. If
// specified, the matching profile will be loaded into the browser session.
// Profiles must be created beforehand.
//
// This is an alias to an internal type.
type BrowserProfileParam = shared.BrowserProfileParam

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). For GPU images, the default is
// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
// GPU images, recommended presets use one of these resolutions with refresh rates
// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
// or recording behavior. If refresh_rate is not provided, it will be automatically
// determined based on the resolution (higher resolutions use lower refresh rates
// to keep bandwidth reasonable).
//
// This is an alias to an internal type.
type BrowserViewport = shared.BrowserViewport

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). For GPU images, the default is
// 1920x1080@60. Arbitrary viewport dimensions and refresh rates are accepted.
// Known-good presets include: 2560x1440@10, 1920x1080@25, 1920x1200@25,
// 1440x900@25, 1280x800@60, 1024x768@60, 1200x800@60, 768x1024@60, 390x844@60. For
// GPU images, recommended presets use one of these resolutions with refresh rates
// 60, 30, 25, or 10: 800x600, 960x720, 1024x576, 1024x768, 1152x648, 1200x800,
// 1280x720, 1368x768, 1440x900, 1600x900, 1920x1080, 1920x1200, 390x844, 360x250,
// 768x1024, 800x1600. Viewports outside this list may exhibit unstable live view
// or recording behavior. If refresh_rate is not provided, it will be automatically
// determined based on the resolution (higher resolutions use lower refresh rates
// to keep bandwidth reasonable).
//
// This is an alias to an internal type.
type BrowserViewportParam = shared.BrowserViewportParam

// This is an alias to an internal type.
type ErrorDetail = shared.ErrorDetail

// An error event from the application.
//
// This is an alias to an internal type.
type ErrorEvent = shared.ErrorEvent

// This is an alias to an internal type.
type ErrorModel = shared.ErrorModel

// Heartbeat event sent periodically to keep SSE connection alive.
//
// This is an alias to an internal type.
type HeartbeatEvent = shared.HeartbeatEvent

// A log entry from the application.
//
// This is an alias to an internal type.
type LogEvent = shared.LogEvent
