package credentials

import (
	"context"
	"net/http"
)

// AccessTokenProvider mints an access token, usually over the network. Only
// [Token.AccessToken] has to be set; a token without [Token.ExpiresAt] is kept
// by a [TokenCache] until it is invalidated.
//
// forceRefresh is true on the first call after [TokenCache.Invalidate]: the
// server has just rejected the previous token, so a provider that keeps its
// own cache must bypass it. Providers without one can ignore the flag.
//
// A provider must be safe for concurrent use.
type AccessTokenProvider func(ctx context.Context, forceRefresh bool) (*Token, error)

// IdentityTokenFunc returns an identity token (a JWT) from the workload's own
// identity provider. It is called once per exchange and must return a fresh
// token every time: most providers issue single-use tokens, and Kernel rejects
// a token it has already seen.
//
// An IdentityTokenFunc must be safe for concurrent use and must return when
// ctx is done.
type IdentityTokenFunc func(ctx context.Context) (string, error)

// HTTPDoer sends an HTTP request. [*http.Client] implements it.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}
