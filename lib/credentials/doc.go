// Package credentials implements workload identity federation for the Kernel
// API: a workload trades a token from its own identity provider (GitHub
// Actions, AWS, Google Cloud, Kubernetes, any OIDC issuer) for a short-lived
// Kernel access token instead of storing an API key.
//
// Most callers never use this package directly. A client picks federation up
// from the environment (KERNEL_FEDERATION_RULE_ID plus an identity-token
// source), or from [github.com/kernel/kernel-go-sdk/option.WithFederationTokenProvider].
// The building blocks are exported for callers that manage tokens themselves:
//
//   - [ExchangeFederationAssertion] performs one exchange.
//   - [NewFederationProvider] performs an exchange with a fresh identity
//     token on every call, and [TokenCache] caches and refreshes the result.
//   - [File], [Static], [GitHubActions] and [GCPMetadata] are identity-token
//     sources. AWS lives in the awssts subpackage so that this package does
//     not depend on the AWS SDK.
//
// Identity tokens are usually single-use, so nothing here ever sends the same
// one twice: every exchange asks the source for a fresh token.
package credentials
