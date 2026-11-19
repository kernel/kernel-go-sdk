package option

import "github.com/kernel/kernel-go-sdk/lib/credentials"

// IdentityTokenFunc returns a fresh identity token (a JWT) from the
// workload's own identity provider, such as SPIFFE/SPIRE or a cloud SDK. It
// is called once per token exchange and must be safe for concurrent use.
type IdentityTokenFunc = credentials.IdentityTokenFunc

// FederationOptions configures [WithFederationTokenProvider].
type FederationOptions = credentials.FederationOptions

// IdentityTokenFile returns an [IdentityTokenFunc] that reads the token from
// path on every call, so a token that is rotated on disk, such as a
// Kubernetes projected service account token, is picked up.
func IdentityTokenFile(path string) IdentityTokenFunc {
	return credentials.File(path)
}

// WithFederationTokenProvider returns a RequestOption that authenticates with
// workload identity federation instead of an API key: provider's identity
// token is exchanged for a short-lived Kernel access token, which is cached
// and renewed before it expires. provider is called for every exchange and
// must return a fresh token each time.
//
// Credential options follow the usual rule that the last one wins, so this
// option replaces an API key from KERNEL_API_KEY or an earlier [WithAPIKey],
// and a later [WithAPIKey] replaces it. Build the option once and reuse it:
// the token cache belongs to the option value.
//
// See the lib/credentials package for the built-in identity-token sources.
func WithFederationTokenProvider(provider IdentityTokenFunc, opts FederationOptions) RequestOption {
	return credentials.ClientOption(provider, opts)
}
