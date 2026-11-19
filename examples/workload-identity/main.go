// Authenticates with workload identity federation instead of an API key.
//
// With KERNEL_FEDERATION_RULE_ID set, and an identity-token source the SDK can
// find (a GitHub Actions job with "id-token: write", KERNEL_IDENTITY_TOKEN_FILE,
// or KERNEL_IDENTITY_TOKEN_SOURCE), no code is needed: kernel.NewClient() does
// the exchange. Pass -file to configure it in code instead.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	kernel "github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/option"
	// Uncomment for KERNEL_IDENTITY_TOKEN_SOURCE=aws_sts:
	// _ "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"
)

func main() {
	ruleID := flag.String("rule", "", "federation rule ID (wifr_...), with -file")
	tokenFile := flag.String("file", "", "path of a file holding the identity token, such as a Kubernetes projected token")
	flag.Parse()

	var opts []option.RequestOption
	if *tokenFile != "" {
		opts = append(opts, option.WithFederationTokenProvider(
			option.IdentityTokenFile(*tokenFile),
			option.FederationOptions{FederationRuleID: *ruleID},
		))
	}
	client := kernel.NewClient(opts...)

	authContext, err := client.Auth.Context.Get(context.Background())
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	fmt.Println(authContext.RawJSON())
}
