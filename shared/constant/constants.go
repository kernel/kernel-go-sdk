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

type AppVersionSummary string        // Always "app_version_summary"
type AwsUsEast1a string              // Always "aws.us-east-1a"
type Console string                  // Always "console"
type ConsoleError string             // Always "console_error"
type ConsoleLog string               // Always "console_log"
type DeploymentState string          // Always "deployment_state"
type Error string                    // Always "error"
type Interaction string              // Always "interaction"
type InteractionClick string         // Always "interaction_click"
type InteractionKey string           // Always "interaction_key"
type InteractionScrollSettled string // Always "interaction_scroll_settled"
type InvocationState string          // Always "invocation_state"
type Log string                      // Always "log"
type ManagedAuthState string         // Always "managed_auth_state"
type MonitorDisconnected string      // Always "monitor_disconnected"
type MonitorInitFailed string        // Always "monitor_init_failed"
type MonitorReconnectFailed string   // Always "monitor_reconnect_failed"
type MonitorReconnected string       // Always "monitor_reconnected"
type MonitorScreenshot string        // Always "monitor_screenshot"
type Network string                  // Always "network"
type NetworkIdle string              // Always "network_idle"
type NetworkLoadingFailed string     // Always "network_loading_failed"
type NetworkRequest string           // Always "network_request"
type NetworkResponse string          // Always "network_response"
type Page string                     // Always "page"
type PageDomContentLoaded string     // Always "page_dom_content_loaded"
type PageLayoutSettled string        // Always "page_layout_settled"
type PageLayoutShift string          // Always "page_layout_shift"
type PageLcp string                  // Always "page_lcp"
type PageLoad string                 // Always "page_load"
type PageNavigation string           // Always "page_navigation"
type PageNavigationSettled string    // Always "page_navigation_settled"
type PageTabOpened string            // Always "page_tab_opened"
type SseHeartbeat string             // Always "sse_heartbeat"
type System string                   // Always "system"

func (c AppVersionSummary) Default() AppVersionSummary { return "app_version_summary" }
func (c AwsUsEast1a) Default() AwsUsEast1a             { return "aws.us-east-1a" }
func (c Console) Default() Console                     { return "console" }
func (c ConsoleError) Default() ConsoleError           { return "console_error" }
func (c ConsoleLog) Default() ConsoleLog               { return "console_log" }
func (c DeploymentState) Default() DeploymentState     { return "deployment_state" }
func (c Error) Default() Error                         { return "error" }
func (c Interaction) Default() Interaction             { return "interaction" }
func (c InteractionClick) Default() InteractionClick   { return "interaction_click" }
func (c InteractionKey) Default() InteractionKey       { return "interaction_key" }
func (c InteractionScrollSettled) Default() InteractionScrollSettled {
	return "interaction_scroll_settled"
}
func (c InvocationState) Default() InvocationState               { return "invocation_state" }
func (c Log) Default() Log                                       { return "log" }
func (c ManagedAuthState) Default() ManagedAuthState             { return "managed_auth_state" }
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
func (c Page) Default() Page                                     { return "page" }
func (c PageDomContentLoaded) Default() PageDomContentLoaded     { return "page_dom_content_loaded" }
func (c PageLayoutSettled) Default() PageLayoutSettled           { return "page_layout_settled" }
func (c PageLayoutShift) Default() PageLayoutShift               { return "page_layout_shift" }
func (c PageLcp) Default() PageLcp                               { return "page_lcp" }
func (c PageLoad) Default() PageLoad                             { return "page_load" }
func (c PageNavigation) Default() PageNavigation                 { return "page_navigation" }
func (c PageNavigationSettled) Default() PageNavigationSettled   { return "page_navigation_settled" }
func (c PageTabOpened) Default() PageTabOpened                   { return "page_tab_opened" }
func (c SseHeartbeat) Default() SseHeartbeat                     { return "sse_heartbeat" }
func (c System) Default() System                                 { return "system" }

func (c AppVersionSummary) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c AwsUsEast1a) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Console) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c ConsoleError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c ConsoleLog) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c DeploymentState) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Interaction) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c InteractionClick) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c InteractionKey) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c InteractionScrollSettled) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c InvocationState) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Log) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c ManagedAuthState) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c MonitorDisconnected) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c MonitorInitFailed) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c MonitorReconnectFailed) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c MonitorReconnected) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c MonitorScreenshot) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c Network) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c NetworkIdle) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c NetworkLoadingFailed) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c NetworkRequest) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c NetworkResponse) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Page) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c PageDomContentLoaded) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c PageLayoutSettled) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c PageLayoutShift) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c PageLcp) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c PageLoad) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c PageNavigation) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c PageNavigationSettled) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c PageTabOpened) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c SseHeartbeat) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c System) MarshalJSON() ([]byte, error)                   { return marshalString(c) }

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
