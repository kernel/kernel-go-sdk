// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package kernel_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/kernel/kernel-go-sdk"
	"github.com/kernel/kernel-go-sdk/internal/testutil"
	"github.com/kernel/kernel-go-sdk/option"
)

func TestWorkloadIdentityFederationRuleNewWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.WorkloadIdentityFederationRules.New(context.TODO(), kernel.WorkloadIdentityFederationRuleNewParams{
		WorkloadIdentityFederationRuleConfig: kernel.WorkloadIdentityFederationRuleConfigParam{
			Authorization: kernel.WorkloadIdentityAuthorizationParam{
				Policy: kernel.WorkloadIdentityAuthorizationPolicyParam{
					Statements: []kernel.WorkloadIdentityAuthorizationPolicyStatementParam{{
						Action:   "browsers:use",
						Resource: "projects/proj_abc123/browsers/br_abc123",
						Condition: kernel.WorkloadIdentityAuthorizationPolicyStatementConditionParam{
							Key:      "source_ip",
							Operator: "in",
							Values:   []string{"203.0.113.0/24"},
						},
					}},
				},
				ProjectID:          "x",
				TokenAudience:      "x",
				MaxTokenTtlSeconds: kernel.Int(1),
			},
			ExpectedAudience: "x",
			IssuerID:         "x",
			Match: kernel.WorkloadIdentityMatchParam{
				Type: kernel.WorkloadIdentityMatchTypeGitHubActions,
				Claims: map[string]string{
					"foo": "x",
				},
				Environment:       kernel.String("environment"),
				EventNames:        []string{"x"},
				JobWorkflowRef:    kernel.String("job_workflow_ref"),
				Namespace:         kernel.String("namespace"),
				Ref:               kernel.String("ref"),
				Repository:        kernel.String("repository"),
				RepositoryID:      kernel.String("269125115713"),
				RepositoryOwnerID: kernel.String("269125115713"),
				ServiceAccount:    kernel.String("service_account"),
				Subject:           kernel.String("x"),
				SubjectPrefix:     kernel.String("x"),
				Workflow:          kernel.String("workflow"),
				WorkflowRef:       kernel.String("workflow_ref"),
			},
			Name: "name",
			Target: kernel.WorkloadIdentityTargetParam{
				ServiceAccountID: "x",
				Type:             kernel.WorkloadIdentityTargetTypeServiceAccount,
			},
			Description: kernel.String("description"),
			Enabled:     kernel.Bool(true),
		},
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityFederationRuleGet(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.WorkloadIdentityFederationRules.Get(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityFederationRuleUpdateWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.WorkloadIdentityFederationRules.Update(
		context.TODO(),
		"id",
		kernel.WorkloadIdentityFederationRuleUpdateParams{
			ExpectedVersion: 1,
			Authorization: kernel.WorkloadIdentityAuthorizationParam{
				Policy: kernel.WorkloadIdentityAuthorizationPolicyParam{
					Statements: []kernel.WorkloadIdentityAuthorizationPolicyStatementParam{{
						Action:   "browsers:use",
						Resource: "projects/proj_abc123/browsers/br_abc123",
						Condition: kernel.WorkloadIdentityAuthorizationPolicyStatementConditionParam{
							Key:      "source_ip",
							Operator: "in",
							Values:   []string{"203.0.113.0/24"},
						},
					}},
				},
				ProjectID:          "x",
				TokenAudience:      "x",
				MaxTokenTtlSeconds: kernel.Int(1),
			},
			Description:      kernel.String("description"),
			Enabled:          kernel.Bool(true),
			ExpectedAudience: kernel.String("x"),
			IssuerID:         kernel.String("x"),
			Match: kernel.WorkloadIdentityMatchParam{
				Type: kernel.WorkloadIdentityMatchTypeGitHubActions,
				Claims: map[string]string{
					"foo": "x",
				},
				Environment:       kernel.String("environment"),
				EventNames:        []string{"x"},
				JobWorkflowRef:    kernel.String("job_workflow_ref"),
				Namespace:         kernel.String("namespace"),
				Ref:               kernel.String("ref"),
				Repository:        kernel.String("repository"),
				RepositoryID:      kernel.String("269125115713"),
				RepositoryOwnerID: kernel.String("269125115713"),
				ServiceAccount:    kernel.String("service_account"),
				Subject:           kernel.String("x"),
				SubjectPrefix:     kernel.String("x"),
				Workflow:          kernel.String("workflow"),
				WorkflowRef:       kernel.String("workflow_ref"),
			},
			Name: kernel.String("name"),
			Target: kernel.WorkloadIdentityTargetParam{
				ServiceAccountID: "x",
				Type:             kernel.WorkloadIdentityTargetTypeServiceAccount,
			},
		},
	)
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityFederationRuleListWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.WorkloadIdentityFederationRules.List(context.TODO(), kernel.WorkloadIdentityFederationRuleListParams{
		IssuerID:         kernel.String("issuer_id"),
		Limit:            kernel.Int(1),
		Offset:           kernel.Int(0),
		ProjectID:        kernel.String("project_id"),
		Query:            kernel.String("query"),
		ServiceAccountID: kernel.String("service_account_id"),
		Status:           kernel.WorkloadIdentityFederationRuleListParamsStatusActive,
	})
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWorkloadIdentityFederationRuleDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := kernel.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	err := client.WorkloadIdentityFederationRules.Delete(context.TODO(), "id")
	if err != nil {
		var apierr *kernel.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
