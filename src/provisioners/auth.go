package provisioners

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"xo/src/bundles"
	"xo/src/massdriver"
	"xo/src/provisioners/terraform"
	"xo/src/telemetry"
	"xo/src/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.opentelemetry.io/otel"
)

type AWSIAMPolicyDocument struct {
	Version   string                   `json:"Version,omitempty"`
	Statement []*AWSIAMPolicyStatement `json:"Statement,omitempty"`
}

type AWSIAMPolicyStatement struct {
	Sid       string                 `json:"Sid,omitempty"`
	Effect    string                 `json:"Effect,omitempty"`
	Action    []string               `json:"Action,omitempty"`
	Resource  []string               `json:"Resource,omitempty"`
	Condition map[string]interface{} `json:"Condition,omitempty"`
}

type STSAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

func GenerateProvisionerAWSCredentials(ctx context.Context, out io.Writer, stsClient STSAPI, spec *massdriver.Specification, roleARN string) error {
	_, span := otel.Tracer("xo").Start(ctx, "provisioners.terraform.ReportProgressFromLogs")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	policy := getProvisionerPolicy(spec)

	policyBytes, marshalErr := json.MarshalIndent(policy, "", "\t")
	if marshalErr != nil {
		util.LogError(marshalErr, span, "error while marshaling the generated AWS policy")
		return marshalErr
	}

	ari := sts.AssumeRoleInput{
		RoleArn:         &roleARN,
		Policy:          aws.String(string(policyBytes)),
		RoleSessionName: aws.String(spec.PackageID),
		DurationSeconds: aws.Int32(7200),
	}

	aro, assumeErr := stsClient.AssumeRole(ctx, &ari)
	if assumeErr != nil {
		util.LogError(marshalErr, span, "error while assume AWS role")
		return assumeErr
	}

	iniConfig := map[string]interface{}{
		"default": map[string]interface{}{
			"aws_access_key_id":     *aro.Credentials.AccessKeyId,
			"aws_secret_access_key": *aro.Credentials.SecretAccessKey,
			"aws_session_token":     *aro.Credentials.SessionToken,
		},
	}

	return renderINI(out, iniConfig)
}

func getProvisionerPolicy(spec *massdriver.Specification) *AWSIAMPolicyDocument {
	policy := AWSIAMPolicyDocument{
		Version:   "2012-10-17",
		Statement: []*AWSIAMPolicyStatement{},
	}

	policyFunctions := []func(*massdriver.Specification) []*AWSIAMPolicyStatement{
		getWorkflowProgressPolicies,
		getAssumeRolePolicies,
		getStateManagementPolicies,
		getBundleReadPolicies,
	}

	for _, policyFunction := range policyFunctions {
		statements := policyFunction(spec)
		policy.Statement = append(policy.Statement, statements...)
	}

	return &policy
}

func getWorkflowProgressPolicies(spec *massdriver.Specification) []*AWSIAMPolicyStatement {
	statements := make([]*AWSIAMPolicyStatement, 0, 1)

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "WorkflowProgressPublisher",
		Effect: "Allow",
		Action: []string{
			"sns:Publish",
		},
		Resource: []string{
			spec.EventTopicARN,
		},
	})

	return statements
}

func getAssumeRolePolicies(spec *massdriver.Specification) []*AWSIAMPolicyStatement {
	statements := make([]*AWSIAMPolicyStatement, 0, 1)

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "AssumeRole",
		Effect: "Allow",
		Action: []string{
			"sts:AssumeRole",
		},
		Resource: []string{
			"*",
		},
	})

	return statements
}

func getStateManagementPolicies(spec *massdriver.Specification) []*AWSIAMPolicyStatement {
	statements := make([]*AWSIAMPolicyStatement, 0, 3)

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "TerraformStateBucketList",
		Effect: "Allow",
		Action: []string{
			"s3:ListBucket",
		},
		Resource: []string{
			bucketNameToARN(spec.S3StateBucket),
		},
	})

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "TerraformStateBucketManage",
		Effect: "Allow",
		Action: []string{
			"s3:GetObject",
			"s3:PutObject",
		},
		Resource: []string{
			path.Join(bucketNameToARN(spec.S3StateBucket), path.Dir(terraform.GetS3StateKey(spec.OrganizationID, spec.PackageID, "remove")), "*"),
		},
	})

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "TerraformStateDynamoDBTableLock",
		Effect: "Allow",
		Action: []string{
			"dynamodb:PutItem",
			"dynamodb:GetItem",
			"dynamodb:DeleteItem",
		},
		Resource: []string{
			spec.DynamoDBStateLockTable,
		},
		// https://www.terraform.io/language/settings/backends/s3#protecting-access-to-workspace-state
		Condition: map[string]interface{}{
			"ForAllValues:StringLike": map[string][]string{
				"dynamodb:LeadingKeys": {
					path.Join(spec.S3StateBucket, spec.OrganizationID, spec.PackageID, "*"),
				},
			},
		},
	})

	return statements
}

func getBundleReadPolicies(spec *massdriver.Specification) []*AWSIAMPolicyStatement {
	statements := make([]*AWSIAMPolicyStatement, 0, 2)

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "BundleBucketList",
		Effect: "Allow",
		Action: []string{
			"s3:ListBucket",
		},
		Resource: []string{
			bucketNameToARN(spec.BundleBucket),
		},
	})

	statements = append(statements, &AWSIAMPolicyStatement{
		Sid:    "BundleBucketRead",
		Effect: "Allow",
		Action: []string{
			"s3:GetObject",
		},
		Resource: []string{
			path.Join(bucketNameToARN(spec.BundleBucket), bundles.GetBundleKey(spec.BundleOwnerOrganizationID, spec.BundleID)),
		},
	})

	return statements
}

func bucketNameToARN(name string) string {
	return fmt.Sprintf("arn:aws:s3:::%s", name)
}
