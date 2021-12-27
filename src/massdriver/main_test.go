package massdriver_test

import (
	"context"
	"testing"
	"xo/src/massdriver"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSTestClient struct {
	Input *sns.PublishInput
	Data  *string
}

func (c *SNSTestClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	c.Input = params
	c.Data = params.Message
	return &sns.PublishOutput{}, nil
}

func TestPublishEventToSNS(t *testing.T) {

	type testData struct {
		name  string
		input *massdriver.Event
		want  string
	}
	tests := []testData{
		{
			name: "Test Decommission Failed",
			input: &massdriver.Event{
				Metadata: massdriver.EventMetadata{
					Timestamp:   "2021-01-01 12:00:00.4321",
					Provisioner: "testaform",
					EventType:   "create_pending",
				},
				Payload: massdriver.EventPayloadProvisionerStatus{
					DeploymentId: "depId",
				}},
			want: `{"metadata":{"timestamp":"2021-01-01 12:00:00.4321","provisioner":"testaform","event_type":"create_pending"},"payload":{"deployment_id":"depId"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			testSNSClient := SNSTestClient{}
			testClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testClient.PublishEventToSNS(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Input
			if *got.Message != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, *got.Message)
			}
			if *got.MessageGroupId != "depId" {
				t.Fatalf("want: %v, got: %v", "depId", *got.MessageGroupId)
			}
		})
	}
}
