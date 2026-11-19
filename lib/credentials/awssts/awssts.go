// Package awssts provides an identity-token source for workloads running on
// AWS (ECS, Lambda, EC2, EKS): it asks AWS STS for a web identity token of the
// workload's IAM role with GetWebIdentityToken.
//
// Importing the package registers the source, so that
// KERNEL_IDENTITY_TOKEN_SOURCE=aws_sts works with a client built from the
// environment:
//
//	import _ "github.com/kernel/kernel-go-sdk/lib/credentials/awssts"
//
// It is a separate package so that only programs that import it depend on the
// AWS SDK.
package awssts

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/kernel/kernel-go-sdk/lib/credentials"
)

const (
	// defaultRegion is used when the AWS configuration names no region.
	// GetWebIdentityToken is served by regional STS endpoints only.
	defaultRegion = "us-east-1"
	// tokenDuration is the lifetime of the identity token. It is used once,
	// straight away.
	tokenDuration = 300
)

func init() {
	credentials.RegisterIdentityTokenSource(credentials.SourceAWSSTS, func(audience string) credentials.IdentityTokenFunc {
		return IdentityToken(audience)
	})
}

// IdentityToken returns an identity-token source that calls STS
// GetWebIdentityToken on every call. It uses the ambient AWS credentials and
// the region from AWS_REGION or AWS_DEFAULT_REGION, falling back to us-east-1.
// The credentials refresh themselves, so the source works for as long as the
// workload runs. optFns adjust how the AWS configuration is loaded.
//
// audience is the aud claim to request; when empty,
// [credentials.ResolveAudience] decides.
func IdentityToken(audience string, optFns ...func(*config.LoadOptions) error) credentials.IdentityTokenFunc {
	var (
		mu     sync.Mutex
		client *sts.Client
	)
	stsClient := func(ctx context.Context) (*sts.Client, error) {
		mu.Lock()
		defer mu.Unlock()
		if client != nil {
			return client, nil
		}
		cfg, err := config.LoadDefaultConfig(ctx, optFns...)
		if err != nil {
			return nil, fmt.Errorf("could not load the AWS configuration: %w", err)
		}
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}
		client = sts.NewFromConfig(cfg)
		return client, nil
	}

	return func(ctx context.Context) (string, error) {
		audience, err := credentials.ResolveAudience(ctx, audience)
		if err != nil {
			return "", err
		}
		client, err := stsClient(ctx)
		if err != nil {
			return "", err
		}
		out, err := client.GetWebIdentityToken(ctx, &sts.GetWebIdentityTokenInput{
			Audience:         []string{audience},
			SigningAlgorithm: aws.String("RS256"),
			DurationSeconds:  aws.Int32(tokenDuration),
		})
		if err != nil {
			return "", fmt.Errorf("AWS STS GetWebIdentityToken failed: %w", err)
		}
		if out.WebIdentityToken == nil || *out.WebIdentityToken == "" {
			return "", errors.New("AWS STS GetWebIdentityToken returned no token")
		}
		return *out.WebIdentityToken, nil
	}
}
