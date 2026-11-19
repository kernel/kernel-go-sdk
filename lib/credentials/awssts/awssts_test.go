package awssts_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/lib/credentials"
	"github.com/kernel/kernel-go-sdk/lib/credentials/awssts"
)

// awsEnv gives the AWS SDK static credentials and keeps it away from the
// machine's own configuration and from the instance metadata service.
func awsEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"AWS_REGION", "AWS_DEFAULT_REGION", "AWS_PROFILE", "AWS_SESSION_TOKEN",
		"AWS_ENDPOINT_URL", "AWS_ENDPOINT_URL_STS", "AWS_WEB_IDENTITY_TOKEN_FILE", "AWS_ROLE_ARN",
		"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "AWS_CONTAINER_CREDENTIALS_FULL_URI",
		"KERNEL_API_KEY", "KERNEL_BASE_URL", "KERNEL_CUSTOM_HEADERS",
		credentials.EnvFederationRuleID, credentials.EnvIdentityTokenSource, credentials.EnvIdentityTokenFile,
		credentials.EnvIdentityToken, credentials.EnvIdentityTokenAudience, credentials.EnvAuthBaseURL,
		"ACTIONS_ID_TOKEN_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_TOKEN",
	} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
	missing := filepath.Join(t.TempDir(), "missing")
	t.Setenv("AWS_CONFIG_FILE", missing)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", missing)
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIDEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
}

// fakeSTS answers GetWebIdentityToken the way STS does and records what it
// was asked. It serves both as an HTTP client for the AWS SDK and as an
// HTTP handler.
type fakeSTS struct {
	mu       sync.Mutex
	requests []stsRequest
}

type stsRequest struct {
	Host          string
	Authorization string
	Form          url.Values
}

func (s *fakeSTS) record(req *http.Request) string {
	body, _ := io.ReadAll(req.Body)
	form, _ := url.ParseQuery(string(body))
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, stsRequest{Host: host, Authorization: req.Header.Get("Authorization"), Form: form})
	token := fmt.Sprintf("aws.identity.token-%d-for-%s", len(s.requests), form.Get("Audience.member.1"))
	return `<GetWebIdentityTokenResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetWebIdentityTokenResult>
    <WebIdentityToken>` + token + `</WebIdentityToken>
    <Expiration>2026-11-19T12:05:00Z</Expiration>
  </GetWebIdentityTokenResult>
  <ResponseMetadata><RequestId>c6104cbe-af31-11e0-8154-cbc7ccf896c7</RequestId></ResponseMetadata>
</GetWebIdentityTokenResponse>`
}

func (s *fakeSTS) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"text/xml"}},
		Body:       io.NopCloser(strings.NewReader(s.record(req))),
		Request:    req,
	}, nil
}

func (s *fakeSTS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	_, _ = io.WriteString(w, s.record(req))
}

func (s *fakeSTS) seen() []stsRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]stsRequest(nil), s.requests...)
}

func TestIdentityTokenRequest(t *testing.T) {
	awsEnv(t)
	t.Setenv("AWS_REGION", "eu-west-1")
	sts := &fakeSTS{}
	source := awssts.IdentityToken("https://auth.onkernel.com", config.WithHTTPClient(sts))

	first, err := source(context.Background())
	if err != nil {
		t.Fatalf("IdentityToken: %v", err)
	}
	second, err := source(context.Background())
	if err != nil {
		t.Fatalf("IdentityToken: %v", err)
	}
	if first != "aws.identity.token-1-for-https://auth.onkernel.com" || second == first {
		t.Fatalf("tokens = %q, %q; want a new STS token per call", first, second)
	}

	req := sts.seen()[0]
	if req.Host != "sts.eu-west-1.amazonaws.com" {
		t.Errorf("host = %q, want the regional endpoint sts.eu-west-1.amazonaws.com", req.Host)
	}
	want := url.Values{
		"Action":            {"GetWebIdentityToken"},
		"Version":           {"2011-06-15"},
		"Audience.member.1": {"https://auth.onkernel.com"},
		"SigningAlgorithm":  {"RS256"},
		"DurationSeconds":   {"300"},
	}
	if req.Form.Encode() != want.Encode() {
		t.Errorf("request = %q\nwant %q", req.Form.Encode(), want.Encode())
	}
}

func TestIdentityTokenRegion(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"AWS_REGION", map[string]string{"AWS_REGION": "ap-southeast-2", "AWS_DEFAULT_REGION": "eu-west-1"}, "ap-southeast-2"},
		{"AWS_DEFAULT_REGION", map[string]string{"AWS_DEFAULT_REGION": "eu-west-1"}, "eu-west-1"},
		{"default", nil, "us-east-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			awsEnv(t)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			sts := &fakeSTS{}
			if _, err := awssts.IdentityToken("aud", config.WithHTTPClient(sts))(context.Background()); err != nil {
				t.Fatalf("IdentityToken: %v", err)
			}
			req := sts.seen()[0]
			if want := "sts." + tt.want + ".amazonaws.com"; req.Host != want {
				t.Errorf("host = %q, want %q", req.Host, want)
			}
			if !strings.Contains(req.Authorization, "/"+tt.want+"/sts/aws4_request") {
				t.Errorf("request signed for another region: %q", req.Authorization)
			}
		})
	}
}

func TestIdentityTokenErrors(t *testing.T) {
	awsEnv(t)
	t.Run("no audience", func(t *testing.T) {
		sts := &fakeSTS{}
		_, err := awssts.IdentityToken("", config.WithHTTPClient(sts))(context.Background())
		if err == nil || !strings.Contains(err.Error(), "KERNEL_IDENTITY_TOKEN_AUDIENCE") {
			t.Fatalf("error = %v, want one about the missing audience", err)
		}
		if n := len(sts.seen()); n != 0 {
			t.Fatalf("%d STS requests, want 0", n)
		}
	})
	t.Run("STS refuses", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/xml")
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `<ErrorResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/"><Error><Type>Sender</Type>`+
				`<Code>AccessDenied</Code><Message>not authorized to perform sts:GetWebIdentityToken</Message></Error>`+
				`<RequestId>1</RequestId></ErrorResponse>`)
		}))
		defer server.Close()
		t.Setenv("AWS_ENDPOINT_URL_STS", server.URL)
		_, err := awssts.IdentityToken("aud")(context.Background())
		if err == nil || !strings.Contains(err.Error(), "GetWebIdentityToken failed") || !strings.Contains(err.Error(), "AccessDenied") {
			t.Fatalf("error = %v, want the STS refusal", err)
		}
	})
}

// TestClientWithAWSSTSFromEnvironment covers the whole path a workload on AWS
// takes: KERNEL_IDENTITY_TOKEN_SOURCE=aws_sts, which works because this test
// binary imports the awssts package.
func TestClientWithAWSSTSFromEnvironment(t *testing.T) {
	awsEnv(t)
	sts := &fakeSTS{}
	stsServer := httptest.NewServer(sts)
	defer stsServer.Close()

	var assertions []string
	auth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		assertions = append(assertions, r.PostForm.Get("assertion"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"token_type": "Bearer", "access_token": "kernel-access-token", "expires_in": 298})
	}))
	defer auth.Close()

	var authorizations []string
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizations = append(authorizations, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer api.Close()

	t.Setenv("AWS_REGION", "us-west-2")
	t.Setenv("AWS_ENDPOINT_URL_STS", stsServer.URL)
	t.Setenv("KERNEL_BASE_URL", api.URL)
	t.Setenv(credentials.EnvAuthBaseURL, auth.URL)
	t.Setenv(credentials.EnvFederationRuleID, "wifr_test123")
	t.Setenv(credentials.EnvIdentityTokenSource, "aws_sts")

	client := kernel.NewClient()
	for i := 0; i < 2; i++ {
		var res map[string]any
		if err := client.Get(context.Background(), "auth/context", nil, &res); err != nil {
			t.Fatalf("request failed: %v", err)
		}
	}

	if got := strings.Join(authorizations, ","); got != "Bearer kernel-access-token,Bearer kernel-access-token" {
		t.Errorf("API requests carried %q", got)
	}
	// The audience defaults to the auth base URL, and the STS token is what
	// gets exchanged.
	if want := "aws.identity.token-1-for-" + auth.URL; len(assertions) != 1 || assertions[0] != want {
		t.Errorf("assertions = %q, want [%q]", assertions, want)
	}
	if requests := sts.seen(); len(requests) != 1 || !strings.Contains(requests[0].Authorization, "/us-west-2/sts/aws4_request") {
		t.Errorf("STS requests = %+v, want one signed for us-west-2", requests)
	}
}
