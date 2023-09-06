package helm_test

import (
	"context"
	"os"
	"testing"
	"xo/src/massdriver"
	"xo/src/provisioners/helm"

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
			name:  "standard",
			input: "testdata/helm-output.log",
			want: []string{
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"create_running"},"payload":{"deployment_id":"id","resource_name":"foo","resource_type":"Namespace","resource_id":"Namespace:foo"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"create_running"},"payload":{"deployment_id":"id","resource_name":"foo/foo-provisioner","resource_type":"Deployment","resource_id":"Deployment:foo/foo-provisioner"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"update_running"},"payload":{"deployment_id":"id","resource_name":"foo/foo-provisioner","resource_type":"Deployment","resource_id":"Deployment:foo/foo-provisioner"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"create_completed"},"payload":{"deployment_id":"id","resource_name":"foo/foo-provisioner","resource_type":"Deployment","resource_id":"Deployment:foo/foo-provisioner"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"delete_running"},"payload":{"deployment_id":"id","resource_name":"foo/foo-provisioner","resource_type":"ServiceAccount","resource_id":"ServiceAccount:foo/foo-provisioner"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"delete_completed"},"payload":{"deployment_id":"id","resource_name":"foo/foo-provisioner","resource_type":"Deployment","resource_id":"Deployment:foo/foo-provisioner"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"Chart.yaml file is missing","error_details":"","error_level":"error"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"UPGRADE FAILED: context deadline exceeded","error_details":"","error_level":"error"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"UPGRADE FAILED: error validating \"\": error validating data: ValidationError(Deployment.spec.template.spec.containers[0]): unknown field \"foo\" in io.k8s.api.core.v1.Container","error_details":"","error_level":"error"}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"helm","version":"1.2.3","event_type":"provisioner_error"},"payload":{"deployment_id":"id","error_message":"uninstall: Release not loaded: foo: release: not found","error_details":"","error_level":"error"}}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "helm")
			t.Setenv("HELM_VERSION", "1.2.3")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testRequests = make([]string, 0)
			testSNSClient := SNSTestClient{}
			testMassdriverClient := massdriver.MassdriverClient{SNSClient: &testSNSClient, Specification: &massdriver.Specification{}}

			input, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			defer input.Close()

			err = helm.ReportProgressFromLogs(context.Background(), &testMassdriverClient, "id", input)
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
