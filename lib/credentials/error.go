package credentials

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// maxErrorFieldLen caps each field copied out of a token endpoint error. The
// description is free-form text that ends up in logs.
const maxErrorFieldLen = 2000

const invalidTargetHint = "Set Resource to the Kernel API origin with no trailing slash."

// replayHint explains an "already used" refusal for the source the identity
// token came from. repeated is set when the same token was sent in the
// previous exchange as well.
func replayHint(source sourceInfo, repeated bool) string {
	switch source.kind {
	case sourceKindStatic:
		return "Every exchange needs a fresh identity token from the provider; " +
			"a static KERNEL_IDENTITY_TOKEN can be exchanged only once."
	case sourceKindFile:
		observed := fmt.Sprintf("The token in the identity token file %s was already exchanged", source.path)
		if repeated {
			observed = fmt.Sprintf("The identity token file %s held the same token as in the previous exchange", source.path)
		}
		return observed + ", so whatever writes the file must rotate it before each refresh. " +
			"For a Kubernetes projected token, the issuer can allow reuse instead."
	}
	if repeated {
		return "The identity token provider returned the same token as for the previous exchange; " +
			"it must return a fresh one every time it is called."
	}
	return "Every exchange needs a fresh identity token; an identity token can be exchanged only once."
}

// WorkloadIdentityError reports a failure to obtain a Kernel access token
// through workload identity federation: a misconfiguration, an identity-token
// source that failed, or a token endpoint that refused the exchange.
type WorkloadIdentityError struct {
	// StatusCode is the HTTP status of the token endpoint's response. It is 0
	// when the failure happened before a response arrived.
	StatusCode int
	// Code is the OAuth error code, such as "invalid_grant".
	Code string
	// Description is the OAuth error_description.
	Description string
	// ErrorURI is the OAuth error_uri.
	ErrorURI string
	// RequestID identifies the token request in Kernel's logs, when the
	// response carried one.
	RequestID string
	// Hint suggests a fix for a failure with a well-known cause.
	Hint string
	// Message describes a failure that has no OAuth error.
	Message string
	// Err is the underlying error, if any.
	Err error

	// retryable marks a failure that a new exchange, with a fresh identity
	// token, may get past.
	retryable bool
	// replay marks the refusal of an identity token that was already used.
	replay bool
}

func (e *WorkloadIdentityError) Error() string {
	var b strings.Builder
	b.WriteString("kernel: workload identity")
	if e.StatusCode != 0 {
		fmt.Fprintf(&b, ": token exchange failed (status %d", e.StatusCode)
		if e.RequestID != "" {
			fmt.Fprintf(&b, ", request id %s", e.RequestID)
		}
		b.WriteString(")")
	}
	for _, part := range []string{e.Message, e.Code, e.Description} {
		if part != "" {
			b.WriteString(": " + part)
		}
	}
	if e.Err != nil {
		b.WriteString(": " + e.Err.Error())
	}
	if e.ErrorURI != "" {
		fmt.Fprintf(&b, " (see %s)", e.ErrorURI)
	}
	if e.Hint != "" {
		b.WriteString(". " + e.Hint)
	}
	return b.String()
}

func (e *WorkloadIdentityError) Unwrap() error { return e.Err }

func configError(format string, args ...any) *WorkloadIdentityError {
	return &WorkloadIdentityError{Message: fmt.Sprintf(format, args...)}
}

// exchangeError builds the error for a token endpoint response that is not a
// success. Only the three OAuth error fields are kept: anything else in the
// body may echo the assertion or a partly minted token back.
func exchangeError(resp *http.Response, body []byte, assertion string) *WorkloadIdentityError {
	e := &WorkloadIdentityError{
		StatusCode: resp.StatusCode,
		RequestID:  requestID(resp, assertion),
		retryable: resp.StatusCode == http.StatusRequestTimeout ||
			resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= http.StatusInternalServerError,
	}
	var fields struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		ErrorURI         string `json:"error_uri"`
	}
	if err := json.Unmarshal(body, &fields); err != nil || fields.Error == "" {
		e.Message = "unexpected response from the token endpoint"
		return e
	}
	e.Code = sanitizeErrorField(fields.Error, assertion)
	e.Description = sanitizeErrorField(fields.ErrorDescription, assertion)
	e.ErrorURI = sanitizeErrorField(fields.ErrorURI, assertion)
	switch {
	case e.Code == "invalid_grant" && strings.Contains(strings.ToLower(e.Description), "already used"):
		e.replay = true
		e.Hint = replayHint(sourceInfo{}, false)
	case e.Code == "invalid_target":
		e.Hint = invalidTargetHint
	}
	return e
}

func sanitizeErrorField(value, assertion string) string {
	if assertion != "" {
		value = strings.ReplaceAll(value, assertion, "[redacted]")
	}
	if len(value) > maxErrorFieldLen {
		value = value[:maxErrorFieldLen] + "...[truncated]"
	}
	return value
}

func requestID(resp *http.Response, assertion string) string {
	id := resp.Header.Get("X-Request-Id")
	if id == "" {
		id = resp.Header.Get("Request-Id")
	}
	return sanitizeErrorField(id, assertion)
}

// redactedError is err with the assertion removed from its text. A transport
// error comes from the caller's HTTP client, which may describe the request
// it was given.
type redactedError struct {
	err  error
	text string
}

func redactError(err error, assertion string) redactedError {
	return redactedError{err: err, text: strings.ReplaceAll(err.Error(), assertion, "[redacted]")}
}

func (e redactedError) Error() string { return e.text }

func (e redactedError) Unwrap() error { return e.err }

// GoString keeps %#v from printing the wrapped error's fields.
func (e redactedError) GoString() string { return fmt.Sprintf("%q", e.text) }
