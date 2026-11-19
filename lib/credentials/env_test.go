package credentials_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kernel/kernel-go-sdk/lib/credentials"
)

func TestFederationFromEnvNotConfigured(t *testing.T) {
	cleanEnv(t)
	// An identity token alone does not turn federation on.
	t.Setenv(credentials.EnvIdentityToken, "a.b.c")
	source, opts, err := credentials.FederationFromEnv()
	if source != nil || err != nil || opts != (credentials.FederationOptions{}) {
		t.Fatalf("got source set = %v, opts %+v, err %v; want nothing", source != nil, opts, err)
	}
}

func TestFederationFromEnvSourceSelection(t *testing.T) {
	tokenFile := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenFile, []byte("from.file.jwt\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	server, _ := githubTokenService(t)
	github := map[string]string{
		"ACTIONS_ID_TOKEN_REQUEST_URL":   server.URL + "/token?api-version=2.0",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN": "runtime-token",
	}
	with := func(base map[string]string, extra ...string) map[string]string {
		env := map[string]string{}
		for k, v := range base {
			env[k] = v
		}
		for i := 0; i < len(extra); i += 2 {
			env[extra[i]] = extra[i+1]
		}
		return env
	}

	tests := []struct {
		name string
		env  map[string]string
		// wantToken is the token the source returns: an exact value, or
		// "github" for a token from the fake runner.
		wantToken string
		wantErr   string
	}{
		{name: "file", env: map[string]string{credentials.EnvIdentityTokenFile: tokenFile}, wantToken: "from.file.jwt"},
		{name: "value", env: map[string]string{credentials.EnvIdentityToken: "from.env.jwt"}, wantToken: "from.env.jwt"},
		{
			name:      "file wins over value",
			env:       map[string]string{credentials.EnvIdentityTokenFile: tokenFile, credentials.EnvIdentityToken: "from.env.jwt"},
			wantToken: "from.file.jwt",
		},
		{name: "GitHub Actions is detected", env: github, wantToken: "github"},
		{
			name:      "file wins over detected GitHub Actions",
			env:       with(github, credentials.EnvIdentityTokenFile, tokenFile),
			wantToken: "from.file.jwt",
		},
		{
			name:      "value wins over detected GitHub Actions",
			env:       with(github, credentials.EnvIdentityToken, "from.env.jwt"),
			wantToken: "from.env.jwt",
		},
		{
			name:      "explicit source wins over file and value",
			env:       with(github, credentials.EnvIdentityTokenSource, "github_actions", credentials.EnvIdentityTokenFile, tokenFile, credentials.EnvIdentityToken, "from.env.jwt"),
			wantToken: "github",
		},
		{
			name:      "explicit env source wins over file",
			env:       map[string]string{credentials.EnvIdentityTokenSource: "env", credentials.EnvIdentityTokenFile: tokenFile, credentials.EnvIdentityToken: "from.env.jwt"},
			wantToken: "from.env.jwt",
		},
		{
			name:      "explicit file source",
			env:       map[string]string{credentials.EnvIdentityTokenSource: "file", credentials.EnvIdentityTokenFile: tokenFile},
			wantToken: "from.file.jwt",
		},
		{name: "no source", env: map[string]string{}, wantErr: "KERNEL_IDENTITY_TOKEN_SOURCE"},
		{name: "file source without a path", env: map[string]string{credentials.EnvIdentityTokenSource: "file"}, wantErr: "requires KERNEL_IDENTITY_TOKEN_FILE"},
		{name: "env source without a value", env: map[string]string{credentials.EnvIdentityTokenSource: "env"}, wantErr: "requires KERNEL_IDENTITY_TOKEN"},
		{name: "unknown source", env: map[string]string{credentials.EnvIdentityTokenSource: "azure"}, wantErr: `unknown identity token source "azure"`},
		{
			name:    "aws_sts without the import",
			env:     map[string]string{credentials.EnvIdentityTokenSource: "aws_sts"},
			wantErr: `_ "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv(t)
			t.Setenv(credentials.EnvFederationRuleID, testRuleID)
			t.Setenv(credentials.EnvIdentityTokenAudience, "https://auth.onkernel.com")
			// A literal token works once per process, so each case gets its own.
			envToken := uniqueToken("from.env")
			for key, value := range tt.env {
				t.Setenv(key, strings.ReplaceAll(value, "from.env.jwt", envToken))
			}
			tt.wantToken = strings.ReplaceAll(tt.wantToken, "from.env.jwt", envToken)

			source, opts, err := credentials.FederationFromEnv()
			if tt.wantErr != "" {
				var wiErr *credentials.WorkloadIdentityError
				if !errors.As(err, &wiErr) || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want a WorkloadIdentityError containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("FederationFromEnv: %v", err)
			}
			if opts != (credentials.FederationOptions{FederationRuleID: testRuleID}) {
				t.Fatalf("options = %+v", opts)
			}
			token, err := source(context.Background())
			if err != nil {
				t.Fatalf("source: %v", err)
			}
			if tt.wantToken == "github" {
				if got := tokenAudiences(token); len(got) != 1 || got[0] != "https://auth.onkernel.com" {
					t.Fatalf("token %q is not from the GitHub runner", token)
				}
				return
			}
			if token != tt.wantToken {
				t.Fatalf("token = %q, want %q", token, tt.wantToken)
			}
		})
	}
}

func TestFederationFromEnvStaticTokenIsSingleUse(t *testing.T) {
	cleanEnv(t)
	t.Setenv(credentials.EnvFederationRuleID, testRuleID)
	t.Setenv(credentials.EnvIdentityToken, uniqueToken("from.env"))
	source, _, err := credentials.FederationFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := source(context.Background()); err != nil {
		t.Fatal(err)
	}
	// Reading the environment again, as every new client does, does not make
	// the token usable again.
	again, _, err := credentials.FederationFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range []credentials.IdentityTokenFunc{source, again} {
		if _, err := source(context.Background()); !errors.Is(err, credentials.ErrIdentityTokenConsumed) {
			t.Fatalf("second use error = %v, want ErrIdentityTokenConsumed", err)
		}
	}
	// A new value in the variable is a new token.
	t.Setenv(credentials.EnvIdentityToken, uniqueToken("replacement"))
	replaced, _, err := credentials.FederationFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := replaced(context.Background()); err != nil {
		t.Fatalf("a replaced token was refused: %v", err)
	}
}
