package provisioners_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"xo/src/massdriver"
	"xo/src/provisioners"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type stsMock struct {
	AssumeRoleOutput sts.AssumeRoleOutput
	AssumeRoleInput  *sts.AssumeRoleInput
}

func (m *stsMock) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	m.AssumeRoleInput = params
	return &m.AssumeRoleOutput, nil
}

func TestGenerateProvisionerAWSCredentials(t *testing.T) {
	buf := bytes.Buffer{}
	stsMock := stsMock{
		AssumeRoleOutput: sts.AssumeRoleOutput{
			Credentials: &types.Credentials{
				AccessKeyId:     aws.String("FAKEACCESSKEYID"),
				SecretAccessKey: aws.String("FAKESECRETACCESSKEY"),
				SessionToken:    aws.String("FakeSessionToken=="),
			},
		},
	}
	spec := massdriver.Specification{
		EventTopicARN:             "arn:aws:sns:eventTopicArn",
		S3StateBucket:             "stateBucket",
		OrganizationID:            "orgId",
		PackageID:                 "packageId",
		DeploymentID:              "deploymentId",
		DynamoDBStateLockTableArn: "aws:aws:dynamodb:table/some-table",
		BundleBucket:              "bundleBucket",
		BundleOwnerOrganizationID: "bundleOrgId",
		BundleID:                  "bundleId",
	}

	roleName := "arn:aws:iam:::role/foo"
	externalId := "foobar"

	err := provisioners.GenerateProvisionerAWSCredentials(context.Background(), &buf, &stsMock, &spec, roleName, externalId)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	gotRole := *stsMock.AssumeRoleInput.RoleArn
	if gotRole != roleName {
		t.Fatalf("want: %v, got: %v", roleName, gotRole)
	}

	gotExternalId := *stsMock.AssumeRoleInput.ExternalId
	if gotExternalId != externalId {
		t.Fatalf("want: %v, got: %v", externalId, gotExternalId)
	}

	gotSessionName := *stsMock.AssumeRoleInput.RoleSessionName
	wantSessionName := spec.DeploymentID
	if gotSessionName != wantSessionName {
		t.Fatalf("want: %v, got: %v", wantSessionName, gotSessionName)
	}

	gotPolicy := *stsMock.AssumeRoleInput.Policy
	wantPolicy := `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "WorkflowProgressPublisher",
			"Effect": "Allow",
			"Action": [
				"sns:Publish"
			],
			"Resource": [
				"arn:aws:sns:eventTopicArn"
			]
		},
		{
			"Sid": "AssumeRole",
			"Effect": "Allow",
			"Action": [
				"sts:AssumeRole"
			],
			"Resource": [
				"*"
			]
		},
		{
			"Sid": "TerraformStateBucketManage",
			"Effect": "Allow",
			"Action": [
				"s3:GetObject",
				"s3:PutObject"
			],
			"Resource": [
				"arn:aws:s3:::stateBucket/orgId/packageId/*"
			]
		},
		{
			"Sid": "TerraformStateDynamoDBTableLock",
			"Effect": "Allow",
			"Action": [
				"dynamodb:PutItem",
				"dynamodb:GetItem",
				"dynamodb:DeleteItem"
			],
			"Resource": [
				"aws:aws:dynamodb:table/some-table"
			],
			"Condition": {
				"ForAllValues:StringLike": {
					"dynamodb:LeadingKeys": [
						"stateBucket/orgId/packageId/*"
					]
				}
			}
		},
		{
			"Sid": "BundleBucketRead",
			"Effect": "Allow",
			"Action": [
				"s3:GetObject"
			],
			"Resource": [
				"arn:aws:s3:::bundleBucket/bundles/bundleOrgId/bundleId/bundle.tar.gz"
			]
		},
		{
			"Sid": "BucketList",
			"Effect": "Allow",
			"Action": [
				"s3:ListBucket"
			],
			"Resource": [
				"arn:aws:s3:::bundleBucket",
				"arn:aws:s3:::stateBucket"
			]
		}
	]
}`

	if gotPolicy != strings.TrimSpace(wantPolicy) {
		t.Fatalf("want: %v, got: %v", strings.TrimSpace(wantPolicy), gotPolicy)
	}

	wantIni := `[default]
aws_access_key_id=FAKEACCESSKEYID
aws_secret_access_key=FAKESECRETACCESSKEY
aws_session_token=FakeSessionToken==
`

	gotIniLines := strings.Split(buf.String(), "\n")
	wantIniLines := strings.Split(wantIni, "\n")

	if len(wantIniLines) != len(gotIniLines) {
		t.Fatalf("want: %v, got: %v", len(wantIniLines), len(gotIniLines))
	}
	if gotIniLines[0] != wantIniLines[0] {
		t.Fatalf("want: %v, got: %v", wantIniLines[0], gotIniLines[0])
	}
	for i := 1; i < len(wantIniLines); i++ {
		if !contains(gotIniLines, wantIniLines[i]) {
			t.Fatalf("Line missing: want: %v, got: %v", wantIniLines[i], gotIniLines)
		}
	}
}

func contains(str []string, eq string) bool {
	for _, val := range str {
		if val == eq {
			return true
		}
	}
	return false
}
