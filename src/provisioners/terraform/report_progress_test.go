package terraform_test

import (
	"context"
	"os"
	"testing"
	"xo/src/massdriver"
	"xo/src/provisioners/terraform"
	testmass "xo/test"
)

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
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Error creating S3 bucket: AccessDenied: Access Denied\n\tstatus code: 403, request id: 8ZJF3ZKYM9QE8Y5Y, host id: PE/mhk+dO5TDoPmLw/wCuKDRUcfuvP+LFx3cFl5EOhfYe0F9fKtmdIG+lAseO2QqufTN+69ihOw=","error_details":"{\"severity\":\"error\",\"summary\":\"Error creating S3 bucket: AccessDenied: Access Denied\\n\\tstatus code: 403, request id: 8ZJF3ZKYM9QE8Y5Y, host id: PE/mhk+dO5TDoPmLw/wCuKDRUcfuvP+LFx3cFl5EOhfYe0F9fKtmdIG+lAseO2QqufTN+69ihOw=\",\"detail\":\"\",\"address\":\"aws_s3_bucket.one\",\"range\":{\"filename\":\"main.tf\",\"start\":{\"line\":6,\"column\":32,\"byte\":109},\"end\":{\"line\":6,\"column\":33,\"byte\":110}},\"snippet\":{\"context\":\"resource \\\"aws_s3_bucket\\\" \\\"one\\\"\",\"code\":\"resource \\\"aws_s3_bucket\\\" \\\"one\\\" {\",\"start_line\":6}}","error_level":"error"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Reference to undeclared input variable","error_details":"{\"severity\":\"error\",\"summary\":\"Reference to undeclared input variable\",\"detail\":\"An input variable with the name \\\"api_gateway\\\" has not been declared. This variable can be declared with a variable \\\"api_gateway\\\" {} block.\",\"address\":\"\",\"range\":{\"filename\":\"main.tf\",\"start\":{\"line\":2,\"column\":23,\"byte\":31},\"end\":{\"line\":2,\"column\":38,\"byte\":46}},\"snippet\":{\"context\":\"locals\",\"code\":\"  api_id = split(\\\"/\\\", var.api_gateway.data.infrastructure.arn)[2]\",\"start_line\":2}}","error_level":"error"}}`,
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
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","version":"1.0.7","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Argument is deprecated","error_details":"{\"severity\":\"warning\",\"summary\":\"Argument is deprecated\",\"detail\":\"This field is being removed and instead the type is fetched from the massdriver.yaml file\",\"address\":\"massdriver_artifact.subnetwork\",\"range\":{\"filename\":\"_artifacts.tf\",\"start\":{\"line\":5,\"column\":26,\"byte\":169},\"end\":{\"line\":5,\"column\":42,\"byte\":185}},\"snippet\":{\"context\":\"resource \\\"massdriver_artifact\\\" \\\"subnetwork\\\"\",\"code\":\"type = \\\"gcp-subnetwork\\\"\",\"start_line\":5}}","error_level":"warning"}}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testMassdriverClient := testmass.NewMassdriverTestClient("")

			input, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			defer input.Close()

			err = terraform.ReportProgressFromLogs(context.Background(), &testMassdriverClient.MassClient, "id", input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			testRequests := testMassdriverClient.GetSNSMessages()

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
