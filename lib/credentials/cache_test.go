package credentials

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a manually advanced clock for a TokenCache.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2026, 11, 19, 12, 0, 0, 0, time.UTC)}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// countingProvider mints tokens "token-1", "token-2", ... that live for
// lifetime on clock.
type countingProvider struct {
	clock    *fakeClock
	lifetime time.Duration

	mu     sync.Mutex
	calls  int
	forced []bool
	err    error
	gate   chan struct{}
}

func (p *countingProvider) provide(ctx context.Context, force bool) (*Token, error) {
	p.mu.Lock()
	p.calls++
	n := p.calls
	p.forced = append(p.forced, force)
	err, gate := p.err, p.gate
	p.mu.Unlock()
	if gate != nil {
		<-gate
	}
	if err != nil {
		return nil, err
	}
	token := &Token{AccessToken: fmt.Sprintf("token-%d", n)}
	if p.lifetime > 0 {
		token.ExpiresAt = p.clock.Now().Add(p.lifetime)
	}
	return token, nil
}

func (p *countingProvider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

func (p *countingProvider) setErr(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.err = err
}

func newTestCache(lifetime time.Duration) (*TokenCache, *countingProvider, *fakeClock) {
	clock := newFakeClock()
	provider := &countingProvider{clock: clock, lifetime: lifetime}
	cache := NewTokenCache(provider.provide)
	cache.now = clock.Now
	return cache, provider, clock
}

func mustToken(t *testing.T, cache *TokenCache) string {
	t.Helper()
	token, err := cache.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	return token.AccessToken
}

// waitFor polls until cond holds; background refreshes finish on their own
// goroutine.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestTokenCacheRefreshThresholdsScaleWithLifetime(t *testing.T) {
	tests := []struct {
		lifetime          time.Duration
		advisory, mandate time.Duration
	}{
		{lifetime: 600 * time.Second, advisory: 120 * time.Second, mandate: 30 * time.Second},
		{lifetime: 300 * time.Second, advisory: 120 * time.Second, mandate: 30 * time.Second},
		{lifetime: 298 * time.Second, advisory: 119200 * time.Millisecond, mandate: 29800 * time.Millisecond},
		{lifetime: 60 * time.Second, advisory: 24 * time.Second, mandate: 6 * time.Second},
		{lifetime: 1 * time.Second, advisory: 400 * time.Millisecond, mandate: 100 * time.Millisecond},
	}
	for _, tt := range tests {
		now := time.Now()
		cached := newCachedToken(&Token{AccessToken: "t", ExpiresAt: now.Add(tt.lifetime)}, now)
		if cached.advisory != tt.advisory || cached.mandatory != tt.mandate {
			t.Errorf("lifetime %v: thresholds = %v advisory, %v mandatory; want %v, %v",
				tt.lifetime, cached.advisory, cached.mandatory, tt.advisory, tt.mandate)
		}
	}
}

func TestTokenCacheServesCachedTokenWhileFresh(t *testing.T) {
	cache, provider, clock := newTestCache(300 * time.Second)
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("first token = %q", got)
	}
	clock.Advance(179 * time.Second) // 121s left: just outside the advisory window
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("token = %q, want the cached one", got)
	}
	if calls := provider.callCount(); calls != 1 {
		t.Fatalf("provider called %d times, want 1", calls)
	}
}

func TestTokenCacheAdvisoryWindowRefreshesInBackground(t *testing.T) {
	cache, provider, clock := newTestCache(300 * time.Second)
	mustToken(t, cache)

	clock.Advance(200 * time.Second) // 100s left: advisory
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("token in the advisory window = %q, want the cached token-1", got)
	}
	waitFor(t, "background refresh", func() bool {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		return cache.cached != nil && cache.cached.token.AccessToken == "token-2"
	})
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("token after the background refresh = %q, want token-2", got)
	}
	if calls := provider.callCount(); calls != 2 {
		t.Fatalf("provider called %d times, want 2", calls)
	}
}

func TestTokenCacheAdvisoryFailureIsSwallowedAndBacksOff(t *testing.T) {
	cache, provider, clock := newTestCache(300 * time.Second)
	mustToken(t, cache)
	provider.setErr(errors.New("token endpoint down"))

	clock.Advance(200 * time.Second)
	for i := 0; i < 20; i++ {
		if got := mustToken(t, cache); got != "token-1" {
			t.Fatalf("token = %q, want the stale token-1", got)
		}
	}
	waitFor(t, "failed background refresh", func() bool {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		return !cache.lastAdvisoryFailure.IsZero()
	})
	// Still inside the backoff: no new attempt, however many requests arrive.
	clock.Advance(advisoryRefreshBackoff - time.Second)
	for i := 0; i < 20; i++ {
		mustToken(t, cache)
	}
	if calls := provider.callCount(); calls != 2 {
		t.Fatalf("provider called %d times during the backoff, want 2", calls)
	}
	// After the backoff the refresh is tried again.
	clock.Advance(2 * time.Second)
	mustToken(t, cache)
	waitFor(t, "second background refresh", func() bool { return provider.callCount() == 3 })
}

func TestTokenCacheMandatoryWindowBlocksAndReturnsErrors(t *testing.T) {
	cache, provider, clock := newTestCache(300 * time.Second)
	mustToken(t, cache)

	clock.Advance(271 * time.Second) // 29s left: mandatory
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("token in the mandatory window = %q, want a fresh token-2", got)
	}

	want := errors.New("token endpoint down")
	provider.setErr(want)
	clock.Advance(280 * time.Second)
	if _, err := cache.Token(context.Background()); !errors.Is(err, want) {
		t.Fatalf("Token error = %v, want %v", err, want)
	}
	// A failed mandatory refresh is not cached.
	provider.setErr(nil)
	if got := mustToken(t, cache); got != "token-4" {
		t.Fatalf("token after recovery = %q, want token-4", got)
	}
}

func TestTokenCacheShortLivedToken(t *testing.T) {
	cache, provider, clock := newTestCache(time.Second)
	mustToken(t, cache)
	clock.Advance(500 * time.Millisecond) // 500ms left, advisory starts at 400ms
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("token = %q, want the cached token-1", got)
	}
	if calls := provider.callCount(); calls != 1 {
		t.Fatalf("provider called %d times, want 1: a 1s token must not be refreshed on every request", calls)
	}
	clock.Advance(450 * time.Millisecond) // 50ms left: mandatory
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("token = %q, want a fresh token-2", got)
	}
}

func TestTokenCacheWithoutExpiryCachesUntilInvalidated(t *testing.T) {
	cache, provider, clock := newTestCache(0)
	mustToken(t, cache)
	clock.Advance(1000 * time.Hour)
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("token = %q, want token-1", got)
	}
	cache.Invalidate()
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("token after Invalidate = %q, want token-2", got)
	}
	if calls := provider.callCount(); calls != 2 {
		t.Fatalf("provider called %d times, want 2", calls)
	}
}

func TestTokenCacheSingleFlight(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	provider.gate = make(chan struct{})

	const callers = 50
	var wg sync.WaitGroup
	tokens := make([]string, callers)
	for i := range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if token, err := cache.Token(context.Background()); err == nil {
				tokens[i] = token.AccessToken
			}
		}()
	}
	waitFor(t, "the provider call", func() bool { return provider.callCount() == 1 })
	close(provider.gate)
	wg.Wait()

	for i, token := range tokens {
		if token != "token-1" {
			t.Fatalf("caller %d got %q, want token-1", i, token)
		}
	}
	if calls := provider.callCount(); calls != 1 {
		t.Fatalf("provider called %d times for %d concurrent callers, want 1", calls, callers)
	}
}

func TestTokenCacheInvalidateForcesRefresh(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	mustToken(t, cache)
	cache.Invalidate()
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("token after Invalidate = %q, want token-2", got)
	}
	mustToken(t, cache)

	provider.mu.Lock()
	defer provider.mu.Unlock()
	if len(provider.forced) != 2 || provider.forced[0] || !provider.forced[1] {
		t.Fatalf("forceRefresh flags = %v, want [false true]", provider.forced)
	}
}

func TestTokenCacheKeepsForcingUntilRefreshSucceeds(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	mustToken(t, cache)
	cache.Invalidate()
	provider.setErr(errors.New("token endpoint down"))
	if _, err := cache.Token(context.Background()); err == nil {
		t.Fatal("Token succeeded, want the provider's error")
	}
	provider.setErr(nil)
	mustToken(t, cache)
	mustToken(t, cache)

	provider.mu.Lock()
	defer provider.mu.Unlock()
	if got := fmt.Sprint(provider.forced); got != "[false true true]" {
		t.Fatalf("forceRefresh flags = %s, want [false true true]", got)
	}
}

func TestTokenCacheInvalidateTokenIgnoresReplacedToken(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	first := mustToken(t, cache)
	cache.invalidateToken(first)
	second := mustToken(t, cache)
	// A slower request that was also rejected with the first token must not
	// throw away its replacement.
	cache.invalidateToken(first)
	if got := mustToken(t, cache); got != second {
		t.Fatalf("token = %q, want %q", got, second)
	}
	if calls := provider.callCount(); calls != 2 {
		t.Fatalf("provider called %d times, want 2", calls)
	}
}

func TestTokenCacheQueuesOneRefreshBehindInvalidatedOne(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	gate := make(chan struct{})
	provider.gate = gate

	tokenOf := func(out chan<- string) {
		token, err := cache.Token(context.Background())
		if err != nil {
			out <- err.Error()
			return
		}
		out <- token.AccessToken
	}
	stale := make(chan string, 1)
	go tokenOf(stale)
	waitFor(t, "the first provider call", func() bool { return provider.callCount() == 1 })

	// The token being fetched is rejected three times over while the fetch is
	// still running. Everyone who asks afterwards shares one more fetch, and
	// it waits for the first to finish: the provider is never called twice
	// at once.
	fresh := make(chan string, 3)
	for range 3 {
		cache.Invalidate()
		go tokenOf(fresh)
	}
	waitFor(t, "the callers to queue", func() bool {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		return cache.refresh != nil && !cache.refresh.started
	})
	// One more rejection that nobody follows up on before the first fetch
	// returns. The queued fetch starts after it, so its token is still good.
	cache.Invalidate()
	time.Sleep(20 * time.Millisecond)
	if calls := provider.callCount(); calls != 1 {
		t.Fatalf("provider called %d times while the first call was running, want 1", calls)
	}
	close(gate)

	if got := <-stale; got != "token-1" {
		t.Fatalf("first caller got %q, want token-1", got)
	}
	for range 3 {
		if got := <-fresh; got != "token-2" {
			t.Fatalf("caller after Invalidate got %q, want token-2", got)
		}
	}
	if got := mustToken(t, cache); got != "token-2" {
		t.Fatalf("cached token = %q, want token-2", got)
	}
	if calls := provider.callCount(); calls != 2 {
		t.Fatalf("provider called %d times, want 2", calls)
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	if !provider.forced[1] {
		t.Fatal("the refresh after Invalidate was not forced")
	}
}

func TestTokenCacheRespectsContextCancellation(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	provider.gate = make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() {
		_, err := cache.Token(ctx)
		errc <- err
	}()
	waitFor(t, "the provider call", func() bool { return provider.callCount() == 1 })
	cancel()
	if err := <-errc; !errors.Is(err, context.Canceled) {
		t.Fatalf("Token error = %v, want context.Canceled", err)
	}

	// The refresh outlives the caller that gave up, so the token it paid for
	// is not wasted.
	close(provider.gate)
	if got := mustToken(t, cache); got != "token-1" {
		t.Fatalf("token = %q, want token-1", got)
	}
	if calls := provider.callCount(); calls != 1 {
		t.Fatalf("provider called %d times, want 1", calls)
	}
}

func TestTokenCacheDoesNothingForDoneContext(t *testing.T) {
	cache, provider, _ := newTestCache(300 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := cache.Token(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Token error = %v, want context.Canceled", err)
	}
	time.Sleep(20 * time.Millisecond)
	if calls := provider.callCount(); calls != 0 {
		t.Fatalf("provider called %d times for a cancelled caller, want 0", calls)
	}
	mustToken(t, cache)
	if _, err := cache.Token(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Token error with a cached token = %v, want context.Canceled", err)
	}
}

func TestTokenCacheProviderMisbehaviour(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		cache := NewTokenCache(func(context.Context, bool) (*Token, error) { panic("boom") })
		if _, err := cache.Token(context.Background()); err == nil {
			t.Fatal("Token succeeded, want an error")
		}
	})
	t.Run("no token", func(t *testing.T) {
		cache := NewTokenCache(func(context.Context, bool) (*Token, error) { return nil, nil })
		if _, err := cache.Token(context.Background()); err == nil {
			t.Fatal("Token succeeded, want an error")
		}
	})
}

func TestTokenCacheConcurrentUse(t *testing.T) {
	cache, _, clock := newTestCache(300 * time.Second)
	var failures atomic.Int32
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 200 {
				if _, err := cache.Token(context.Background()); err != nil {
					failures.Add(1)
				}
				switch {
				case j%50 == 0 && i%4 == 0:
					cache.Invalidate()
				case j%20 == 0:
					clock.Advance(40 * time.Second)
				}
			}
		}()
	}
	wg.Wait()
	if n := failures.Load(); n != 0 {
		t.Fatalf("%d Token calls failed", n)
	}
}
