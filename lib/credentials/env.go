package credentials

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Environment variables that configure workload identity federation.
const (
	// EnvFederationRuleID turns federation on. It holds the federation rule
	// ID ("wifr_...").
	EnvFederationRuleID = "KERNEL_FEDERATION_RULE_ID"
	// EnvIdentityTokenSource selects where identity tokens come from:
	// github_actions, gcp_metadata, aws_sts, file or env.
	EnvIdentityTokenSource = "KERNEL_IDENTITY_TOKEN_SOURCE"
	// EnvIdentityTokenFile is the path of a file holding the identity token.
	// It is read again for every exchange.
	EnvIdentityTokenFile = "KERNEL_IDENTITY_TOKEN_FILE"
	// EnvIdentityToken is a literal identity token. It can be exchanged once.
	EnvIdentityToken = "KERNEL_IDENTITY_TOKEN"
	// EnvIdentityTokenAudience is the audience the built-in sources request.
	// It defaults to the auth base URL.
	EnvIdentityTokenAudience = "KERNEL_IDENTITY_TOKEN_AUDIENCE"
	// EnvAuthBaseURL is the Kernel authorization server. It is required with
	// a custom API base URL.
	EnvAuthBaseURL = "KERNEL_AUTH_BASE_URL"

	envAPIKey = "KERNEL_API_KEY"
)

// FederationFromEnv reads the federation settings from the environment:
// KERNEL_FEDERATION_RULE_ID plus an identity-token source. The source is the
// one KERNEL_IDENTITY_TOKEN_SOURCE names. Without that variable it is
// KERNEL_IDENTITY_TOKEN_FILE, then KERNEL_IDENTITY_TOKEN, then GitHub Actions
// when the job exposes an OIDC token.
//
// It returns a nil source and no error when KERNEL_FEDERATION_RULE_ID is not
// set. It does not look at KERNEL_API_KEY.
func FederationFromEnv() (IdentityTokenFunc, FederationOptions, error) {
	ruleID := os.Getenv(EnvFederationRuleID)
	if ruleID == "" {
		return nil, FederationOptions{}, nil
	}
	source, err := identityTokenSourceFromEnv()
	if err != nil {
		return nil, FederationOptions{}, &WorkloadIdentityError{
			Message: EnvFederationRuleID + " is set but no identity token source could be configured",
			Err:     err,
		}
	}
	return source, FederationOptions{FederationRuleID: ruleID}, nil
}

func identityTokenSourceFromEnv() (IdentityTokenFunc, error) {
	file := os.Getenv(EnvIdentityTokenFile)
	value := os.Getenv(EnvIdentityToken)

	name := os.Getenv(EnvIdentityTokenSource)
	if name == "" {
		switch {
		case file != "":
			name = SourceFile
		case value != "":
			name = SourceEnv
		case os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") != "":
			name = SourceGitHubActions
		default:
			return nil, fmt.Errorf("set %s (%s, %s, %s, %s or %s), %s or %s",
				EnvIdentityTokenSource, SourceGitHubActions, SourceGCPMetadata, SourceAWSSTS, SourceFile, SourceEnv,
				EnvIdentityTokenFile, EnvIdentityToken)
		}
	}

	switch name {
	case SourceFile:
		if file == "" {
			return nil, fmt.Errorf("%s=%s requires %s", EnvIdentityTokenSource, SourceFile, EnvIdentityTokenFile)
		}
		return File(file), nil
	case SourceEnv:
		if value == "" {
			return nil, fmt.Errorf("%s=%s requires %s", EnvIdentityTokenSource, SourceEnv, EnvIdentityToken)
		}
		return Static(value), nil
	}
	newSource, err := identityTokenSource(name)
	if err != nil {
		return nil, err
	}
	return newSource(""), nil
}

var warnAPIKeyWins sync.Once

// apiKeyWinsInEnv reports whether KERNEL_API_KEY is set, which takes
// precedence over the federation variables. It says so once when both are
// present, because the federation variables are then silently unused. This
// is about the environment only: an explicit option still beats both.
func apiKeyWinsInEnv() bool {
	if os.Getenv(envAPIKey) == "" {
		return false
	}
	if os.Getenv(EnvFederationRuleID) != "" {
		warnAPIKeyWins.Do(func() {
			log.Printf("kernel: %s and %s are both set; %s takes precedence over the federation variables",
				envAPIKey, EnvFederationRuleID, envAPIKey)
		})
	}
	return true
}
