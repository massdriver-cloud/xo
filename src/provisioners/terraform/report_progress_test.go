package terraform_test

import (
	"os"
	"testing"
	"xo/src/provisioners/terraform"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"google.golang.org/protobuf/proto"
)

var testRequests []*mdproto.ProvisionerProgressUpdateRequest

func storeInSlice(message *mdproto.ProvisionerProgressUpdateRequest) error {
	testRequests = append(testRequests, message)
	return nil
}

func TestReportProgressFromLogs(t *testing.T) {
	type testData struct {
		name     string
		input    string
		expected []mdproto.ProvisionerProgressUpdateRequest
	}
	tests := []testData{
		{
			name:  "empty schema",
			input: "testdata/terraform-output.ndjson",
			expected: []mdproto.ProvisionerProgressUpdateRequest{
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status:    mdproto.ProvisionerStatus_PROVISIONER_STATUS_PLAN_COMPLETED,
					Timestamp: "2021-10-22T09:49:55.916597-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status:    mdproto.ProvisionerStatus_PROVISIONER_STATUS_APPLY_COMPLETED,
					Timestamp: "2021-10-22T09:50:01.810341-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status:    mdproto.ProvisionerStatus_PROVISIONER_STATUS_DESTROY_COMPLETED,
					Timestamp: "2021-10-22T09:52:28.062482-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status:    mdproto.ProvisionerStatus_PROVISIONER_STATUS_ERROR,
					Timestamp: "2021-10-22T10:07:26.909862-07:00",
					Error: &mdproto.ProvisionerError{
						Message: "Error creating S3 bucket: AccessDenied: Access Denied\n\tstatus code: 403, request id: 8ZJF3ZKYM9QE8Y5Y, host id: PE/mhk+dO5TDoPmLw/wCuKDRUcfuvP+LFx3cFl5EOhfYe0F9fKtmdIG+lAseO2QqufTN+69ihOw=",
						Level:   mdproto.ProvisionerErrorLevel_PROVISIONER_ERROR_LEVEL_ERROR,
					},
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_PENDING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_CREATE,
					},
					Timestamp: "2021-10-22T10:07:25.323404-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_PENDING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_UPDATE,
					},
					Timestamp: "2021-10-22T10:07:25.323575-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_PENDING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE,
					},
					Timestamp: "2021-10-22T09:55:15.644830-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_PENDING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_RECREATE,
					},
					Timestamp: "2021-10-22T09:56:37.437385-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_RUNNING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_CREATE,
					},
					Timestamp: "2021-10-22T10:07:26.509289-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						ResourceId:   "242c983b-ff05-4b81-8dd4-afbac03ea364",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_RUNNING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_UPDATE,
					},
					Timestamp: "2021-10-22T09:54:38.720344-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						ResourceId:   "242c983b-ff05-4b81-8dd4-afbac03ea364",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_RUNNING,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE,
					},
					Timestamp: "2021-10-22T09:52:27.587225-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						ResourceKey:  "key",
						ResourceId:   "242c983b-ff05-4b81-8dd4-afbac03ea364",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_COMPLETED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_CREATE,
					},
					Timestamp: "2021-10-22T09:50:01.745859-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						ResourceKey:  "0",
						ResourceId:   "242c983b-ff05-4b81-8dd4-afbac03ea364",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_COMPLETED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_UPDATE,
					},
					Timestamp: "2021-10-22T09:54:41.915057-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_COMPLETED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE,
					},
					Timestamp: "2021-10-22T09:52:28.004763-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_FAILED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_CREATE,
					},
					Timestamp: "2021-10-22T10:01:42.821634-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_FAILED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_UPDATE,
					},
					Timestamp: "2021-10-22T10:01:42.821634-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "two",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_FAILED,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE,
					},
					Timestamp: "2021-10-22T10:01:42.821634-07:00",
				},
				{
					DeploymentId:    "id",
					DeploymentToken: "token",
					Metadata: &mdproto.ProvisionerMetadata{
						Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
						ProvisionerVersion: "1.0.7",
					},
					Status: mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE,
					ResourceProgress: &mdproto.ProvisionerResourceProgress{
						ResourceType: "aws_s3_bucket",
						ResourceName: "one",
						Status:       mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_DRIFT,
						Action:       mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE,
					},
					Timestamp: "2021-10-22T10:00:21.180582-07:00",
				},
			},
		},
	}

	terraform.ReportProgressSender = storeInSlice
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testRequests = make([]*mdproto.ProvisionerProgressUpdateRequest, 0)

			input, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			err = terraform.ReportProgressFromLogs("id", "token", input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(tc.expected) != len(testRequests) {
				t.Fatalf("expected: %v, got: %v", len(tc.expected), len(testRequests))
			}

			for i := 0; i < len(tc.expected); i++ {
				if !proto.Equal(&tc.expected[i], testRequests[i]) {
					t.Fatalf("expected: %v, got: %v", tc.expected[i], *testRequests[i])
				}
			}
		})
	}
}
