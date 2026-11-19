package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kernel/kernel-go-sdk/internal"
)

const (
	// Values of KERNEL_IDENTITY_TOKEN_SOURCE.
	SourceGitHubActions = "github_actions"
	SourceGCPMetadata   = "gcp_metadata"
	SourceAWSSTS        = "aws_sts"
	SourceFile          = "file"
	SourceEnv           = "env"

	awsSTSImportPath = "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"

	sourceRequestTimeout = 30 * time.Second
)

var (
	sourcesMu sync.RWMutex
	sources   = map[string]func(audience string) IdentityTokenFunc{
		SourceGitHubActions: GitHubActions,
		SourceGCPMetadata:   GCPMetadata,
	}
)

// RegisterIdentityTokenSource makes an identity-token source selectable with
// KERNEL_IDENTITY_TOKEN_SOURCE=name. It is meant to be called from the init
// function of a package that the application imports for that purpose, as
// database/sql drivers are:
//
//	import _ "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"
//
// newSource receives the configured audience, which is empty when the source
// should fall back to [ResolveAudience]. RegisterIdentityTokenSource panics if
// newSource is nil or name is already taken.
func RegisterIdentityTokenSource(name string, newSource func(audience string) IdentityTokenFunc) {
	sourcesMu.Lock()
	defer sourcesMu.Unlock()
	if newSource == nil {
		panic("credentials: RegisterIdentityTokenSource: source is nil")
	}
	if _, taken := sources[name]; taken || name == SourceFile || name == SourceEnv {
		panic("credentials: RegisterIdentityTokenSource: source name " + name + " is already taken")
	}
	sources[name] = newSource
}

func identityTokenSource(name string) (func(audience string) IdentityTokenFunc, error) {
	sourcesMu.RLock()
	newSource, ok := sources[name]
	sourcesMu.RUnlock()
	if ok {
		return newSource, nil
	}
	if name == SourceAWSSTS {
		return nil, fmt.Errorf("identity token source %q is not compiled in: add the import _ %q to your program", name, awsSTSImportPath)
	}
	return nil, fmt.Errorf("unknown identity token source %q (want %s, %s, %s, %s or %s)",
		name, SourceGitHubActions, SourceGCPMetadata, SourceAWSSTS, SourceFile, SourceEnv)
}

// sourceHTTPClient talks to GitHub's token service. Redirects are refused
// because the request carries the job's bearer token.
var sourceHTTPClient = &http.Client{
	Timeout:       sourceRequestTimeout,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

// metadataHTTPClient talks to the Google Cloud metadata server, which is
// link-local and must not go through a proxy.
var metadataHTTPClient = &http.Client{
	Timeout:       sourceRequestTimeout,
	Transport:     &http.Transport{Proxy: nil},
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

// GitHubActions returns an identity-token source for a GitHub Actions job. It
// asks the runner for a new OIDC token on every call, which works for the
// whole job. The workflow must grant "permissions: id-token: write".
//
// audience is the aud claim to request; when empty, [ResolveAudience] decides.
func GitHubActions(audience string) IdentityTokenFunc {
	return func(ctx context.Context) (string, error) {
		requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
		requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
		if requestURL == "" || requestToken == "" {
			return "", errors.New("ACTIONS_ID_TOKEN_REQUEST_URL and ACTIONS_ID_TOKEN_REQUEST_TOKEN are not set: " +
				"the GitHub Actions identity token is only available in a workflow job that grants `permissions: id-token: write`")
		}
		audience, err := ResolveAudience(ctx, audience)
		if err != nil {
			return "", err
		}
		separator := "&"
		if !strings.Contains(requestURL, "?") {
			separator = "?"
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL+separator+"audience="+url.QueryEscape(audience), nil)
		if err != nil {
			return "", errors.New("ACTIONS_ID_TOKEN_REQUEST_URL is not a valid URL")
		}
		req.Header.Set("Authorization", "Bearer "+requestToken)
		req.Header.Set("Accept", "application/json")
		body, err := fetchIdentityToken(sourceHTTPClient, req, "GitHub Actions")
		if err != nil {
			return "", err
		}
		var parsed struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil || parsed.Value == "" {
			return "", errors.New("the GitHub Actions identity token response has no token")
		}
		return parsed.Value, nil
	}
}

// GCPMetadata returns an identity-token source for Google Cloud (Compute
// Engine, GKE, Cloud Run, Cloud Functions). It asks the metadata server for
// an ID token of the default service account on every call.
//
// audience is the aud claim to request; when empty, [ResolveAudience] decides.
func GCPMetadata(audience string) IdentityTokenFunc {
	return func(ctx context.Context) (string, error) {
		audience, err := ResolveAudience(ctx, audience)
		if err != nil {
			return "", err
		}
		host := os.Getenv("GCE_METADATA_HOST")
		if host == "" {
			host = "metadata.google.internal"
		}
		query := url.Values{"audience": {audience}, "format": {"full"}}
		endpoint := "http://" + host + "/computeMetadata/v1/instance/service-accounts/default/identity?" + query.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return "", errors.New("GCE_METADATA_HOST is not a valid host")
		}
		req.Header.Set("Metadata-Flavor", "Google")
		body, err := fetchIdentityToken(metadataHTTPClient, req, "Google Cloud metadata")
		if err != nil {
			return "", err
		}
		token := strings.TrimSpace(string(body))
		if token == "" {
			return "", errors.New("the Google Cloud metadata server returned an empty identity token")
		}
		return token, nil
	}
}

// transientError marks a source failure that another attempt may get past:
// the provider's endpoint was unreachable, overloaded or rate limited.
type transientError struct{ error }

func (e transientError) Unwrap() error { return e.error }

// fetchIdentityToken sends req and returns the body of a 200 response. The
// body of any other response is dropped: it is not needed to diagnose the
// failure and is not known to be free of secrets.
func fetchIdentityToken(client *http.Client, req *http.Request, provider string) ([]byte, error) {
	req.Header.Set("User-Agent", fmt.Sprintf("Kernel/Go %s", internal.PackageVersion))
	resp, err := client.Do(req)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			// The request URL can carry the job's runtime token parameters.
			err = urlErr.Err
		}
		return nil, transientError{fmt.Errorf("%s identity token request failed: %w", provider, err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%s identity token request failed with status %d", provider, resp.StatusCode)
		if resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= http.StatusInternalServerError {
			return nil, transientError{err}
		}
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxTokenResponseSize))
	if err != nil {
		return nil, transientError{fmt.Errorf("%s identity token request failed: %w", provider, err)}
	}
	return body, nil
}
