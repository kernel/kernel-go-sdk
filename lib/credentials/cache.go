package credentials

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// A token is refreshed in the background once less than
	// min(advisoryRefreshThreshold, advisoryRefreshFraction of its lifetime)
	// remains, and synchronously once less than
	// min(mandatoryRefreshThreshold, mandatoryRefreshFraction of its lifetime)
	// remains. Scaling by lifetime keeps both windows meaningful for tokens
	// that live for seconds rather than minutes.
	advisoryRefreshThreshold  = 120 * time.Second
	mandatoryRefreshThreshold = 30 * time.Second
	advisoryRefreshFraction   = 0.4
	mandatoryRefreshFraction  = 0.1

	// advisoryRefreshBackoff is how long a failed background refresh waits
	// before the next one, so an outage is not retried at request rate.
	advisoryRefreshBackoff = 5 * time.Second

	// refreshTimeout bounds one call to the provider. A refresh is shared by
	// every caller waiting on it, so it runs detached from any one caller's
	// context and needs its own deadline.
	refreshTimeout = 60 * time.Second
)

// TokenCache caches the token of an [AccessTokenProvider] and refreshes it
// before it expires.
//
//   - With plenty of lifetime left, [TokenCache.Token] returns the cached token.
//   - In the advisory window it still returns the cached token and starts a
//     refresh in the background. A failed background refresh is dropped and
//     not repeated for a few seconds.
//   - In the mandatory window, or with no token, it waits for a refresh and
//     returns its error if it fails.
//
// Concurrent callers share one refresh, and the provider is never called
// twice at the same time. A TokenCache is safe for concurrent use.
type TokenCache struct {
	provider AccessTokenProvider
	now      func() time.Time

	mu                  sync.Mutex
	cached              *cachedToken
	refresh             *refreshCall
	generation          uint64
	forceNext           bool
	lastAdvisoryFailure time.Time
}

type cachedToken struct {
	token     Token
	advisory  time.Duration
	mandatory time.Duration
}

type refreshCall struct {
	done chan struct{}
	// started is set once the provider has been called. Until then the call
	// is queued behind an older one.
	started bool
	// generation is the cache's generation when the provider was called.
	generation uint64
	// force is whether the provider was told to bypass its own cache.
	force bool
	token *Token
	err   error
}

// NewTokenCache returns a [TokenCache] that gets its tokens from provider.
func NewTokenCache(provider AccessTokenProvider) *TokenCache {
	return &TokenCache{provider: provider, now: time.Now}
}

// Token returns a valid access token, refreshing it first when needed.
//
// Cancelling ctx stops the wait, not the refresh: other callers may be
// waiting on the same refresh, and the identity token it spent cannot be
// used again.
func (c *TokenCache) Token(ctx context.Context) (*Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	if cached := c.cached; cached != nil {
		token := cached.token
		remaining := token.ExpiresAt.Sub(c.now())
		if token.ExpiresAt.IsZero() || remaining > cached.advisory {
			c.mu.Unlock()
			return &token, nil
		}
		if remaining > cached.mandatory {
			if c.refresh == nil && c.now().Sub(c.lastAdvisoryFailure) >= advisoryRefreshBackoff {
				c.startRefresh(ctx, true)
			}
			c.mu.Unlock()
			return &token, nil
		}
	}
	call := c.refresh
	// A refresh that called the provider before the last Invalidate may
	// return the token that was just rejected, so it is not joined: a new
	// one is queued behind it. A queued refresh has not called the provider
	// yet and serves everyone who arrives before it does.
	if call == nil || (call.started && call.generation != c.generation) {
		call = c.startRefresh(ctx, false)
	}
	c.mu.Unlock()

	select {
	case <-call.done:
		if call.err != nil {
			return nil, call.err
		}
		token := *call.token
		return &token, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Invalidate drops the cached token, so the next [TokenCache.Token] call
// fetches a new one. Call it when the server rejects the token.
func (c *TokenCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.invalidateLocked()
}

// invalidateToken invalidates the cache only while it still holds
// accessToken. When several requests are rejected with the same token, the
// first one invalidates it and the rest reuse the replacement.
func (c *TokenCache) invalidateToken(accessToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cached != nil && c.cached.token.AccessToken == accessToken {
		c.invalidateLocked()
	}
}

func (c *TokenCache) invalidateLocked() {
	c.cached = nil
	c.generation++
	c.forceNext = true
}

// startRefresh starts a provider call that every waiter shares. It runs
// after the refresh already in flight, if any. c.mu must be held.
func (c *TokenCache) startRefresh(ctx context.Context, advisory bool) *refreshCall {
	previous := c.refresh
	call := &refreshCall{done: make(chan struct{})}
	c.refresh = call

	ctx = context.WithoutCancel(ctx)
	go func() {
		if previous != nil {
			<-previous.done
		}
		c.mu.Lock()
		call.started = true
		call.generation = c.generation
		call.force = c.forceNext
		c.forceNext = false
		c.mu.Unlock()

		refreshCtx, cancel := context.WithTimeout(ctx, refreshTimeout)
		defer cancel()
		token, err := c.callProvider(refreshCtx, call.force)
		c.finishRefresh(call, token, err, advisory)
	}()
	return call
}

func (c *TokenCache) callProvider(ctx context.Context, force bool) (token *Token, err error) {
	// The provider runs caller code (an IdentityTokenFunc) on a goroutine the
	// caller cannot recover on.
	defer func() {
		if r := recover(); r != nil {
			token, err = nil, fmt.Errorf("credentials: token provider panicked: %v", r)
		}
	}()
	token, err = c.provider(ctx, force)
	if err == nil && (token == nil || token.AccessToken == "") {
		return nil, errors.New("credentials: token provider returned no token")
	}
	return token, err
}

func (c *TokenCache) finishRefresh(call *refreshCall, token *Token, err error, advisory bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.refresh == call {
		c.refresh = nil
	}
	call.err = err
	if err != nil {
		// The rejected token has not been replaced yet, so the next attempt
		// must bypass the provider's cache too.
		if call.force {
			c.forceNext = true
		}
		if advisory {
			c.lastAdvisoryFailure = c.now()
		}
		close(call.done)
		return
	}
	call.token = token
	if call.generation == c.generation {
		c.cached = newCachedToken(token, c.now())
		c.lastAdvisoryFailure = time.Time{}
	}
	close(call.done)
}

func newCachedToken(token *Token, now time.Time) *cachedToken {
	cached := &cachedToken{token: *token}
	if lifetime := token.ExpiresAt.Sub(now); !token.ExpiresAt.IsZero() && lifetime > 0 {
		cached.advisory = min(advisoryRefreshThreshold, time.Duration(float64(lifetime)*advisoryRefreshFraction))
		cached.mandatory = min(mandatoryRefreshThreshold, time.Duration(float64(lifetime)*mandatoryRefreshFraction))
	}
	return cached
}
