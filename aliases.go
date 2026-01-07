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
type BrowserProfile = shared.BrowserProfile

// Profile selection for the browser session. Provide either id or name. If
// specified, the matching profile will be loaded into the browser session.
// Profiles must be created beforehand.
//
// This is an alias to an internal type.
type BrowserProfileParam = shared.BrowserProfileParam

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). Only specific viewport configurations are
// supported. The server will reject unsupported combinations. Supported
// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
// automatically determined from the width and height if they match a supported
// configuration exactly. Note: Higher resolutions may affect the responsiveness of
// live view browser
//
// This is an alias to an internal type.
type BrowserViewport = shared.BrowserViewport

// Initial browser window size in pixels with optional refresh rate. If omitted,
// image defaults apply (1920x1080@25). Only specific viewport configurations are
// supported. The server will reject unsupported combinations. Supported
// resolutions are: 2560x1440@10, 1920x1080@25, 1920x1200@25, 1440x900@25,
// 1024x768@60, 1200x800@60 If refresh_rate is not provided, it will be
// automatically determined from the width and height if they match a supported
// configuration exactly. Note: Higher resolutions may affect the responsiveness of
// live view browser
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
