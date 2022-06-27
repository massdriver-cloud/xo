package opa_test

import (
	"context"
	"os"
	"testing"
	"xo/src/massdriver"
	"xo/src/provisioners/opa"

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
			name:  "empty",
			input: "testdata/opa-output-empty.ndjson",
			want:  []string{},
		},
		{
			name:  "2 values",
			input: "testdata/opa-output-multiple.ndjson",
			want: []string{
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provisioner_violation"},"payload":{"deployment_id":"id","opa_rule":"Deletion Violation","opa_value":["random_pet.one","random_pet.two"]}}`,
				`{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"provisioner_violation"},"payload":{"deployment_id":"id","opa_rule":"Deletion Violation","opa_value":["random_pet.three"]}}`,
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

			err = opa.ReportResults(context.Background(), &testMassdriverClient, "id", input)
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
