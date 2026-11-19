package credentials

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// ErrIdentityTokenConsumed is returned by a [Static] identity-token source
// when it is asked for its token a second time.
var ErrIdentityTokenConsumed = errors.New(
	"the static identity token was already used: an identity token can be exchanged only once, " +
		"so a fixed value (such as KERNEL_IDENTITY_TOKEN) cannot renew the Kernel access token; " +
		"use KERNEL_IDENTITY_TOKEN_FILE or another source that returns a fresh token for every exchange")

type audienceKey struct{}

// sourceInfo is what a built-in source tells the federation provider about
// itself, so that a refused token can be explained in terms of where it came
// from. A source the provider knows nothing about leaves it empty.
type sourceInfo struct {
	kind string
	path string
}

const (
	sourceKindFile   = "file"
	sourceKindStatic = "static"
)

type sourceInfoKey struct{}

func withSourceInfo(ctx context.Context, info *sourceInfo) context.Context {
	return context.WithValue(ctx, sourceInfoKey{}, info)
}

func describeSource(ctx context.Context, info sourceInfo) {
	if target, _ := ctx.Value(sourceInfoKey{}).(*sourceInfo); target != nil {
		*target = info
	}
}

func withAudience(ctx context.Context, audience string) context.Context {
	return context.WithValue(ctx, audienceKey{}, audience)
}

// ResolveAudience returns the audience an identity token should be issued
// for: explicit when it is set, else KERNEL_IDENTITY_TOKEN_AUDIENCE, else the
// URL of the authorization server the token is about to be sent to, which a
// federation provider passes through ctx (for example
// "https://auth.onkernel.com"). The federation rule's expected_audience must
// be set to the same value.
//
// The built-in sources use it, and a custom [IdentityTokenFunc] can too.
func ResolveAudience(ctx context.Context, explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if audience := os.Getenv(EnvIdentityTokenAudience); audience != "" {
		return audience, nil
	}
	if audience, _ := ctx.Value(audienceKey{}).(string); audience != "" {
		return audience, nil
	}
	return "", fmt.Errorf("no identity token audience: pass one or set %s", EnvIdentityTokenAudience)
}

// File returns an identity-token source that reads the token from path on
// every call, so a token that is rotated on disk, such as a Kubernetes
// projected service account token, is picked up. Surrounding whitespace is
// trimmed and an empty file is an error.
//
// Between rotations every call returns the same token, so the issuer must
// allow a token to be exchanged more than once, as Kernel does for Kubernetes
// by default. With a single-use issuer, something has to write a new token to
// the file before each exchange.
func File(path string) IdentityTokenFunc {
	return func(ctx context.Context) (string, error) {
		describeSource(ctx, sourceInfo{kind: sourceKindFile, path: path})
		if path == "" {
			return "", errors.New("identity token file path is empty")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("could not read the identity token file: %w", err)
		}
		token := strings.TrimSpace(string(data))
		if token == "" {
			return "", fmt.Errorf("identity token file %q is empty", path)
		}
		return token, nil
	}
}

// usedStaticTokens holds the SHA-256 of every token a [Static] source has
// handed out. It is shared by all of them: each client built from the
// environment makes its own source for the same KERNEL_IDENTITY_TOKEN.
var usedStaticTokens sync.Map

// Static returns an identity-token source that returns token once per
// process. Any later call, by this source or by another one made for the same
// token, fails with [ErrIdentityTokenConsumed] instead of replaying a token
// the server would reject. A client built on it works until its first access
// token expires.
func Static(token string) IdentityTokenFunc {
	token = strings.TrimSpace(token)
	digest := sha256.Sum256([]byte(token))
	return func(ctx context.Context) (string, error) {
		describeSource(ctx, sourceInfo{kind: sourceKindStatic})
		if token == "" {
			return "", errors.New("identity token value is empty")
		}
		if _, used := usedStaticTokens.LoadOrStore(digest, struct{}{}); used {
			return "", ErrIdentityTokenConsumed
		}
		return token, nil
	}
}
