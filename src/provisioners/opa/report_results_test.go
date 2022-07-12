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
			input: "testdata/opa-output-empty.json",
			want:  []string{},
		},
		{
			name:  "6 violations to 12 events (error + diagnostic for each violation)",
			input: "testdata/opa-output-multiple.json",
			want: []string{
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"four\",\"resource_type\":\"random_pet\",\"resource_key\":\"foo\",\"resource_id\":\"ace-mutt\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"ace-mutt\\\",\\\"resource_key\\\":\\\"foo\\\",\\\"resource_name\\\":\\\"four\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"two\",\"resource_type\":\"random_pet\",\"resource_id\":\"loving-zebra\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"loving-zebra\\\",\\\"resource_key\\\":\\\"\\\",\\\"resource_name\\\":\\\"two\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"three\",\"resource_type\":\"random_pet\",\"resource_key\":\"0\",\"resource_id\":\"moving-newt\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"moving-newt\\\",\\\"resource_key\\\":\\\"0\\\",\\\"resource_name\\\":\\\"three\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"three\",\"resource_type\":\"random_pet\",\"resource_key\":\"1\",\"resource_id\":\"steady-pegasus\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"steady-pegasus\\\",\\\"resource_key\\\":\\\"1\\\",\\\"resource_name\\\":\\\"three\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"one\",\"resource_type\":\"random_pet\",\"resource_id\":\"striking-skunk\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"striking-skunk\\\",\\\"resource_key\\\":\\\"\\\",\\\"resource_name\\\":\\\"one\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"resource_update_failed\"},\"payload\":{\"deployment_id\":\"id\",\"resource_name\":\"four\",\"resource_type\":\"random_pet\",\"resource_key\":\"bar\",\"resource_id\":\"upward-ox\"}}",
                "{\"metadata\":{\"timestamp\":\"2021-01-01 12:00:00.1234\",\"provisioner\":\"testaform\",\"event_type\":\"provisioner_error\"},\"payload\":{\"deployment_id\":\"id\",\"error_message\":\"data.terraform.deletion_violations[x]\",\"error_details\":\"{\\\"resource_id\\\":\\\"upward-ox\\\",\\\"resource_key\\\":\\\"bar\\\",\\\"resource_name\\\":\\\"four\\\",\\\"resource_type\\\":\\\"random_pet\\\"}\",\"error_level\":\"error\"}}",
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
