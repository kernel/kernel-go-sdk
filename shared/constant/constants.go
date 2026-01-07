// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package constant

import (
	shimjson "github.com/onkernel/kernel-go-sdk/internal/encoding/json"
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

type AppVersionSummary string // Always "app_version_summary"
type AwsUsEast1a string       // Always "aws.us-east-1a"
type DeploymentState string   // Always "deployment_state"
type Error string             // Always "error"
type InvocationState string   // Always "invocation_state"
type Log string               // Always "log"
type SseHeartbeat string      // Always "sse_heartbeat"

func (c AppVersionSummary) Default() AppVersionSummary { return "app_version_summary" }
func (c AwsUsEast1a) Default() AwsUsEast1a             { return "aws.us-east-1a" }
func (c DeploymentState) Default() DeploymentState     { return "deployment_state" }
func (c Error) Default() Error                         { return "error" }
func (c InvocationState) Default() InvocationState     { return "invocation_state" }
func (c Log) Default() Log                             { return "log" }
func (c SseHeartbeat) Default() SseHeartbeat           { return "sse_heartbeat" }

func (c AppVersionSummary) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c AwsUsEast1a) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c DeploymentState) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c InvocationState) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c Log) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c SseHeartbeat) MarshalJSON() ([]byte, error)      { return marshalString(c) }

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
