package terraform_test

import (
	"context"
	"os"
	"testing"
	"xo/src/massdriver"
	"xo/src/provisioners/terraform"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var testRequests []string

type SNSTestClient struct {
	Data *string
}

func (c *SNSTestClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	testRequests = append(testRequests, *params.Message)
	return &sns.PublishOutput{}, nil
}

func TestReportProgressFromLogs(t *testing.T) {
	type testData struct {
		name  string
		input string
		want  []string
	}
	tests := []testData{
		{
			name:  "empty schema",
			input: "testdata/terraform-output.ndjson",
			want: []string{
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Error creating S3 bucket: AccessDenied: Access Denied\n\tstatus code: 403, request id: 8ZJF3ZKYM9QE8Y5Y, host id: PE/mhk+dO5TDoPmLw/wCuKDRUcfuvP+LFx3cFl5EOhfYe0F9fKtmdIG+lAseO2QqufTN+69ihOw=","error_details":"","error_level":"error"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"create_pending"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"update_pending"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"delete_pending"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"recreate_pending"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"create_running"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"update_running"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket","resource_id":"242c983b-ff05-4b81-8dd4-afbac03ea364"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"delete_running"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket","resource_id":"242c983b-ff05-4b81-8dd4-afbac03ea364"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"create_completed"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket","resource_key":"key","resource_id":"242c983b-ff05-4b81-8dd4-afbac03ea364"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"update_completed"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket","resource_key":"0","resource_id":"242c983b-ff05-4b81-8dd4-afbac03ea364"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"delete_completed"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"create_failed"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"update_failed"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"delete_failed"},"payload":{"deployment_id":"id","resource_name":"two","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"drift_detected"},"payload":{"deployment_id":"id","resource_name":"one","resource_type":"aws_s3_bucket"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Argument is deprecated","error_details":"This field is being removed and instead the type is fetched from the massdriver.yaml file","error_level":"warning"}}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testRequests = make([]string, 0)
			testSNSClient := SNSTestClient{}
			testMassdriverClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}

			input, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			defer input.Close()

			err = terraform.ReportProgressFromLogs(context.Background(), &testMassdriverClient, "id", input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(tc.want) != len(testRequests) {
				t.Fatalf("want: %v, got: %v", len(tc.want), len(testRequests))
			}

			for i := 0; i < len(tc.want); i++ {
				if tc.want[i] != testRequests[i] {
					t.Fatalf("want: %v, got: %v", tc.want[i], testRequests[i])
				}
			}
		})
	}
}
