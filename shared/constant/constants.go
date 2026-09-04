// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package constant

import (
	shimjson "github.com/kernel/kernel-go-sdk/internal/encoding/json"
)

type Constant[T any] interface {
	Default() T
}

// ValueOf gives the default value of a constant from its type. It's helpful when
// constructing constants as variants in a one-of. Note that empty structs are
// marshalled by default. Usage: constant.ValueOf[constant.Foo]()
func ValueOf[T Constant[T]]() T {
	var t T
	return t.Default()
}

type Agentcard string                       // Always "agentcard"
type APICall string                         // Always "api_call"
type AppVersionSummary string               // Always "app_version_summary"
type AutofillTrigger string                 // Always "Autofill.trigger"
type AwsUsEast1a string                     // Always "aws.us-east-1a"
type BrowserCancelDownload string           // Always "Browser.cancelDownload"
type BrowserClose string                    // Always "Browser.close"
type BrowserSetContentsSize string          // Always "Browser.setContentsSize"
type BrowserSetWindowBounds string          // Always "Browser.setWindowBounds"
type Captcha string                         // Always "captcha"
type CaptchaChallengeResult string          // Always "captcha_challenge_result"
type CaptchaSolveResult string              // Always "captcha_solve_result"
type CaptchaSolveStarted string             // Always "captcha_solve_started"
type Card string                            // Always "card"
type CardEnrollment string                  // Always "card_enrollment"
type CdpCommand string                      // Always "cdp_command"
type CdpConnect string                      // Always "cdp_connect"
type CdpDisconnect string                   // Always "cdp_disconnect"
type Collect string                         // Always "collect"
type Connection string                      // Always "connection"
type Console string                         // Always "console"
type ConsoleError string                    // Always "console_error"
type ConsoleLog string                      // Always "console_log"
type Control string                         // Always "control"
type DeploymentState string                 // Always "deployment_state"
type Direct string                          // Always "direct"
type DomFocus string                        // Always "DOM.focus"
type DomScrollIntoViewIfNeeded string       // Always "DOM.scrollIntoViewIfNeeded"
type DomSetFileInputFiles string            // Always "DOM.setFileInputFiles"
type EmbeddedCeremony string                // Always "embedded_ceremony"
type Error string                           // Always "error"
type InputCancelDragging string             // Always "Input.cancelDragging"
type InputDispatchDragEvent string          // Always "Input.dispatchDragEvent"
type InputDispatchKeyEvent string           // Always "Input.dispatchKeyEvent"
type InputDispatchMouseEvent string         // Always "Input.dispatchMouseEvent"
type InputDispatchTouchEvent string         // Always "Input.dispatchTouchEvent"
type InputEmulateTouchFromMouseEvent string // Always "Input.emulateTouchFromMouseEvent"
type InputImeSetComposition string          // Always "Input.imeSetComposition"
type InputInsertText string                 // Always "Input.insertText"
type InputSynthesizePinchGesture string     // Always "Input.synthesizePinchGesture"
type InputSynthesizeScrollGesture string    // Always "Input.synthesizeScrollGesture"
type InputSynthesizeTapGesture string       // Always "Input.synthesizeTapGesture"
type Interaction string                     // Always "interaction"
type InteractionClick string                // Always "interaction_click"
type InteractionKey string                  // Always "interaction_key"
type InteractionScrollSettled string        // Always "interaction_scroll_settled"
type InvocationState string                 // Always "invocation_state"
type Link string                            // Always "link"
type LinkOAuth string                       // Always "link_oauth"
type LiveViewConnect string                 // Always "live_view_connect"
type LiveViewDisconnect string              // Always "live_view_disconnect"
type Log string                             // Always "log"
type Managed string                         // Always "managed"
type ManagedAuthState string                // Always "managed_auth_state"
type Mfa string                             // Always "mfa"
type Monitor string                         // Always "monitor"
type MonitorDisconnected string             // Always "monitor_disconnected"
type MonitorInitFailed string               // Always "monitor_init_failed"
type MonitorReconnectFailed string          // Always "monitor_reconnect_failed"
type MonitorReconnected string              // Always "monitor_reconnected"
type MonitorScreenshot string               // Always "monitor_screenshot"
type Network string                         // Always "network"
type NetworkIdle string                     // Always "network_idle"
type NetworkLoadingFailed string            // Always "network_loading_failed"
type NetworkRequest string                  // Always "network_request"
type NetworkResponse string                 // Always "network_response"
type NoRecommendation string                // Always "no_recommendation"
type Page string                            // Always "page"
type PageCrashed string                     // Always "page_crashed"
type PageDomContentLoaded string            // Always "page_dom_content_loaded"
type PageLayoutSettled string               // Always "page_layout_settled"
type PageLayoutShift string                 // Always "page_layout_shift"
type PageLcp string                         // Always "page_lcp"
type PageLoad string                        // Always "page_load"
type PageNavigation string                  // Always "page_navigation"
type PageNavigationSettled string           // Always "page_navigation_settled"
type PageTabOpened string                   // Always "page_tab_opened"
type PageBringToFront string                // Always "Page.bringToFront"
type PageCaptureScreenshot string           // Always "Page.captureScreenshot"
type PageCaptureSnapshot string             // Always "Page.captureSnapshot"
type PageClose string                       // Always "Page.close"
type PageHandleJavaScriptDialog string      // Always "Page.handleJavaScriptDialog"
type PageNavigate string                    // Always "Page.navigate"
type PageNavigateToHistoryEntry string      // Always "Page.navigateToHistoryEntry"
type PagePrintToPdf string                  // Always "Page.printToPDF"
type PageReload string                      // Always "Page.reload"
type PageSetWebLifecycleState string        // Always "Page.setWebLifecycleState"
type PageStartScreencast string             // Always "Page.startScreencast"
type PageStopLoading string                 // Always "Page.stopLoading"
type PageStopScreencast string              // Always "Page.stopScreencast"
type Platform string                        // Always "platform"
type PlatformAPICall string                 // Always "platform_api_call"
type ProxyError string                      // Always "proxy_error"
type PushApproval string                    // Always "push_approval"
type Recommendation string                  // Always "recommendation"
type Screenshot string                      // Always "screenshot"
type ServiceCrashed string                  // Always "service_crashed"
type SpendApproval string                   // Always "spend_approval"
type SseHeartbeat string                    // Always "sse_heartbeat"
type System string                          // Always "system"
type SystemOomKill string                   // Always "system_oom_kill"
type TargetActivateTarget string            // Always "Target.activateTarget"
type TargetCloseTarget string               // Always "Target.closeTarget"
type TargetCreateBrowserContext string      // Always "Target.createBrowserContext"
type TargetCreateTarget string              // Always "Target.createTarget"
type TargetDisposeBrowserContext string     // Always "Target.disposeBrowserContext"
type TargetOpenDevTools string              // Always "Target.openDevTools"
type Wallet string                          // Always "wallet"

func (c Agentcard) Default() Agentcard                           { return "agentcard" }
func (c APICall) Default() APICall                               { return "api_call" }
func (c AppVersionSummary) Default() AppVersionSummary           { return "app_version_summary" }
func (c AutofillTrigger) Default() AutofillTrigger               { return "Autofill.trigger" }
func (c AwsUsEast1a) Default() AwsUsEast1a                       { return "aws.us-east-1a" }
func (c BrowserCancelDownload) Default() BrowserCancelDownload   { return "Browser.cancelDownload" }
func (c BrowserClose) Default() BrowserClose                     { return "Browser.close" }
func (c BrowserSetContentsSize) Default() BrowserSetContentsSize { return "Browser.setContentsSize" }
func (c BrowserSetWindowBounds) Default() BrowserSetWindowBounds { return "Browser.setWindowBounds" }
func (c Captcha) Default() Captcha                               { return "captcha" }
func (c CaptchaChallengeResult) Default() CaptchaChallengeResult { return "captcha_challenge_result" }
func (c CaptchaSolveResult) Default() CaptchaSolveResult         { return "captcha_solve_result" }
func (c CaptchaSolveStarted) Default() CaptchaSolveStarted       { return "captcha_solve_started" }
func (c Card) Default() Card                                     { return "card" }
func (c CardEnrollment) Default() CardEnrollment                 { return "card_enrollment" }
func (c CdpCommand) Default() CdpCommand                         { return "cdp_command" }
func (c CdpConnect) Default() CdpConnect                         { return "cdp_connect" }
func (c CdpDisconnect) Default() CdpDisconnect                   { return "cdp_disconnect" }
func (c Collect) Default() Collect                               { return "collect" }
func (c Connection) Default() Connection                         { return "connection" }
func (c Console) Default() Console                               { return "console" }
func (c ConsoleError) Default() ConsoleError                     { return "console_error" }
func (c ConsoleLog) Default() ConsoleLog                         { return "console_log" }
func (c Control) Default() Control                               { return "control" }
func (c DeploymentState) Default() DeploymentState               { return "deployment_state" }
func (c Direct) Default() Direct                                 { return "direct" }
func (c DomFocus) Default() DomFocus                             { return "DOM.focus" }
func (c DomScrollIntoViewIfNeeded) Default() DomScrollIntoViewIfNeeded {
	return "DOM.scrollIntoViewIfNeeded"
}
func (c DomSetFileInputFiles) Default() DomSetFileInputFiles       { return "DOM.setFileInputFiles" }
func (c EmbeddedCeremony) Default() EmbeddedCeremony               { return "embedded_ceremony" }
func (c Error) Default() Error                                     { return "error" }
func (c InputCancelDragging) Default() InputCancelDragging         { return "Input.cancelDragging" }
func (c InputDispatchDragEvent) Default() InputDispatchDragEvent   { return "Input.dispatchDragEvent" }
func (c InputDispatchKeyEvent) Default() InputDispatchKeyEvent     { return "Input.dispatchKeyEvent" }
func (c InputDispatchMouseEvent) Default() InputDispatchMouseEvent { return "Input.dispatchMouseEvent" }
func (c InputDispatchTouchEvent) Default() InputDispatchTouchEvent { return "Input.dispatchTouchEvent" }
func (c InputEmulateTouchFromMouseEvent) Default() InputEmulateTouchFromMouseEvent {
	return "Input.emulateTouchFromMouseEvent"
}
func (c InputImeSetComposition) Default() InputImeSetComposition { return "Input.imeSetComposition" }
func (c InputInsertText) Default() InputInsertText               { return "Input.insertText" }
func (c InputSynthesizePinchGesture) Default() InputSynthesizePinchGesture {
	return "Input.synthesizePinchGesture"
}
func (c InputSynthesizeScrollGesture) Default() InputSynthesizeScrollGesture {
	return "Input.synthesizeScrollGesture"
}
func (c InputSynthesizeTapGesture) Default() InputSynthesizeTapGesture {
	return "Input.synthesizeTapGesture"
}
func (c Interaction) Default() Interaction           { return "interaction" }
func (c InteractionClick) Default() InteractionClick { return "interaction_click" }
func (c InteractionKey) Default() InteractionKey     { return "interaction_key" }
func (c InteractionScrollSettled) Default() InteractionScrollSettled {
	return "interaction_scroll_settled"
}
func (c InvocationState) Default() InvocationState               { return "invocation_state" }
func (c Link) Default() Link                                     { return "link" }
func (c LinkOAuth) Default() LinkOAuth                           { return "link_oauth" }
func (c LiveViewConnect) Default() LiveViewConnect               { return "live_view_connect" }
func (c LiveViewDisconnect) Default() LiveViewDisconnect         { return "live_view_disconnect" }
func (c Log) Default() Log                                       { return "log" }
func (c Managed) Default() Managed                               { return "managed" }
func (c ManagedAuthState) Default() ManagedAuthState             { return "managed_auth_state" }
func (c Mfa) Default() Mfa                                       { return "mfa" }
func (c Monitor) Default() Monitor                               { return "monitor" }
func (c MonitorDisconnected) Default() MonitorDisconnected       { return "monitor_disconnected" }
func (c MonitorInitFailed) Default() MonitorInitFailed           { return "monitor_init_failed" }
func (c MonitorReconnectFailed) Default() MonitorReconnectFailed { return "monitor_reconnect_failed" }
func (c MonitorReconnected) Default() MonitorReconnected         { return "monitor_reconnected" }
func (c MonitorScreenshot) Default() MonitorScreenshot           { return "monitor_screenshot" }
func (c Network) Default() Network                               { return "network" }
func (c NetworkIdle) Default() NetworkIdle                       { return "network_idle" }
func (c NetworkLoadingFailed) Default() NetworkLoadingFailed     { return "network_loading_failed" }
func (c NetworkRequest) Default() NetworkRequest                 { return "network_request" }
func (c NetworkResponse) Default() NetworkResponse               { return "network_response" }
func (c NoRecommendation) Default() NoRecommendation             { return "no_recommendation" }
func (c Page) Default() Page                                     { return "page" }
func (c PageCrashed) Default() PageCrashed                       { return "page_crashed" }
func (c PageDomContentLoaded) Default() PageDomContentLoaded     { return "page_dom_content_loaded" }
func (c PageLayoutSettled) Default() PageLayoutSettled           { return "page_layout_settled" }
func (c PageLayoutShift) Default() PageLayoutShift               { return "page_layout_shift" }
func (c PageLcp) Default() PageLcp                               { return "page_lcp" }
func (c PageLoad) Default() PageLoad                             { return "page_load" }
func (c PageNavigation) Default() PageNavigation                 { return "page_navigation" }
func (c PageNavigationSettled) Default() PageNavigationSettled   { return "page_navigation_settled" }
func (c PageTabOpened) Default() PageTabOpened                   { return "page_tab_opened" }
func (c PageBringToFront) Default() PageBringToFront             { return "Page.bringToFront" }
func (c PageCaptureScreenshot) Default() PageCaptureScreenshot   { return "Page.captureScreenshot" }
func (c PageCaptureSnapshot) Default() PageCaptureSnapshot       { return "Page.captureSnapshot" }
func (c PageClose) Default() PageClose                           { return "Page.close" }
func (c PageHandleJavaScriptDialog) Default() PageHandleJavaScriptDialog {
	return "Page.handleJavaScriptDialog"
}
func (c PageNavigate) Default() PageNavigate { return "Page.navigate" }
func (c PageNavigateToHistoryEntry) Default() PageNavigateToHistoryEntry {
	return "Page.navigateToHistoryEntry"
}
func (c PagePrintToPdf) Default() PagePrintToPdf { return "Page.printToPDF" }
func (c PageReload) Default() PageReload         { return "Page.reload" }
func (c PageSetWebLifecycleState) Default() PageSetWebLifecycleState {
	return "Page.setWebLifecycleState"
}
func (c PageStartScreencast) Default() PageStartScreencast   { return "Page.startScreencast" }
func (c PageStopLoading) Default() PageStopLoading           { return "Page.stopLoading" }
func (c PageStopScreencast) Default() PageStopScreencast     { return "Page.stopScreencast" }
func (c Platform) Default() Platform                         { return "platform" }
func (c PlatformAPICall) Default() PlatformAPICall           { return "platform_api_call" }
func (c ProxyError) Default() ProxyError                     { return "proxy_error" }
func (c PushApproval) Default() PushApproval                 { return "push_approval" }
func (c Recommendation) Default() Recommendation             { return "recommendation" }
func (c Screenshot) Default() Screenshot                     { return "screenshot" }
func (c ServiceCrashed) Default() ServiceCrashed             { return "service_crashed" }
func (c SpendApproval) Default() SpendApproval               { return "spend_approval" }
func (c SseHeartbeat) Default() SseHeartbeat                 { return "sse_heartbeat" }
func (c System) Default() System                             { return "system" }
func (c SystemOomKill) Default() SystemOomKill               { return "system_oom_kill" }
func (c TargetActivateTarget) Default() TargetActivateTarget { return "Target.activateTarget" }
func (c TargetCloseTarget) Default() TargetCloseTarget       { return "Target.closeTarget" }
func (c TargetCreateBrowserContext) Default() TargetCreateBrowserContext {
	return "Target.createBrowserContext"
}
func (c TargetCreateTarget) Default() TargetCreateTarget { return "Target.createTarget" }
func (c TargetDisposeBrowserContext) Default() TargetDisposeBrowserContext {
	return "Target.disposeBrowserContext"
}
func (c TargetOpenDevTools) Default() TargetOpenDevTools { return "Target.openDevTools" }
func (c Wallet) Default() Wallet                         { return "wallet" }

func (c Agentcard) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c APICall) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c AppVersionSummary) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c AutofillTrigger) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c AwsUsEast1a) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c BrowserCancelDownload) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c BrowserClose) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c BrowserSetContentsSize) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c BrowserSetWindowBounds) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Captcha) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c CaptchaChallengeResult) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c CaptchaSolveResult) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c CaptchaSolveStarted) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c Card) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c CardEnrollment) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c CdpCommand) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c CdpConnect) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c CdpDisconnect) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Collect) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Connection) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Console) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c ConsoleError) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ConsoleLog) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Control) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c DeploymentState) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Direct) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c DomFocus) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c DomScrollIntoViewIfNeeded) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c DomSetFileInputFiles) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c EmbeddedCeremony) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c InputCancelDragging) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c InputDispatchDragEvent) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c InputDispatchKeyEvent) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c InputDispatchMouseEvent) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c InputDispatchTouchEvent) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c InputEmulateTouchFromMouseEvent) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c InputImeSetComposition) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c InputInsertText) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c InputSynthesizePinchGesture) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c InputSynthesizeScrollGesture) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c InputSynthesizeTapGesture) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c Interaction) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c InteractionClick) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c InteractionKey) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c InteractionScrollSettled) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c InvocationState) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Link) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c LinkOAuth) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c LiveViewConnect) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c LiveViewDisconnect) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Log) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Managed) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c ManagedAuthState) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Mfa) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Monitor) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c MonitorDisconnected) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c MonitorInitFailed) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c MonitorReconnectFailed) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c MonitorReconnected) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c MonitorScreenshot) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Network) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c NetworkIdle) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c NetworkLoadingFailed) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c NetworkRequest) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c NetworkResponse) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c NoRecommendation) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Page) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c PageCrashed) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c PageDomContentLoaded) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c PageLayoutSettled) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c PageLayoutShift) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c PageLcp) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c PageLoad) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c PageNavigation) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c PageNavigationSettled) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c PageTabOpened) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c PageBringToFront) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c PageCaptureScreenshot) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c PageCaptureSnapshot) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c PageClose) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c PageHandleJavaScriptDialog) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c PageNavigate) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c PageNavigateToHistoryEntry) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c PagePrintToPdf) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c PageReload) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c PageSetWebLifecycleState) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c PageStartScreencast) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c PageStopLoading) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c PageStopScreencast) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Platform) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c PlatformAPICall) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c ProxyError) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c PushApproval) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Recommendation) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Screenshot) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c ServiceCrashed) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c SpendApproval) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c SseHeartbeat) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c System) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c SystemOomKill) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c TargetActivateTarget) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c TargetCloseTarget) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c TargetCreateBrowserContext) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c TargetCreateTarget) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c TargetDisposeBrowserContext) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c TargetOpenDevTools) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Wallet) MarshalJSON() ([]byte, error)                          { return marshalString(c) }

type constant[T any] interface {
	Constant[T]
	*T
}

func marshalString[T ~string, PT constant[T]](v T) ([]byte, error) {
	var zero T
	if v == zero {
		v = PT(&v).Default()
	}
	return shimjson.Marshal(string(v))
}
