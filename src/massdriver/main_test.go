package massdriver_test

import (
	"testing"

	"xo/src/massdriver"
	testmass "xo/test"
)

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
			deploymentId := "depId"
			testClient := testmass.NewMassdriverTestClient(deploymentId)
			testClient.MassClient.Specification.DeploymentID = deploymentId
			err := testClient.MassClient.PublishEvent(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			testClient.GetSNS()
			got := testClient.GetSNS().Input
			if *got.Message != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, *got.Message)
			}
			if *got.MessageGroupId != deploymentId {
				t.Fatalf("want: %v, got: %v", deploymentId, *got.MessageGroupId)
			}
		})
	}
}
