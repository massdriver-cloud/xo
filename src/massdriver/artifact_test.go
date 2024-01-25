package massdriver_test

import (
	"testing"
	"xo/src/massdriver"
	testmass "xo/test"
)

func TestPublishArtifact(t *testing.T) {
	type testData struct {
		name         string
		deploymentId string
		artifact     map[string]interface{}
		want         string
	}
	tests := []testData{
		{
			name:         "Test Artifact Update",
			deploymentId: "id",
			artifact:     map[string]interface{}{"foo": map[string]interface{}{"bar": "baz"}},
			want:         `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"artifact_updated"},"payload":{"deployment_id":"id","artifact":{"foo":{"bar":"baz"}}}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testClient := testmass.NewMassdriverTestClient("")
			testClient.MassClient.Specification.DeploymentID = tc.deploymentId
			err := massdriver.PublishArtifact(&testClient.MassClient, tc.artifact)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := testClient.GetSNSMessages()
			if got[0] != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got[0])
			}
		})
	}
}

func TestDeleteArtifact(t *testing.T) {
	type testData struct {
		name         string
		deploymentId string
		artifact     map[string]interface{}
		want         string
	}
	tests := []testData{
		{
			name:         "Test Artifact Delete",
			deploymentId: "id",
			artifact:     map[string]interface{}{"foo": map[string]interface{}{"bar": "baz"}},
			want:         `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"artifact_deleted"},"payload":{"deployment_id":"id","artifact":{"foo":{"bar":"baz"}}}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testClient := testmass.NewMassdriverTestClient("")
			testClient.MassClient.Specification.DeploymentID = tc.deploymentId
			err := massdriver.DeleteArtifact(&testClient.MassClient, tc.artifact)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := testClient.GetSNSMessages()
			if got[0] != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got[0])
			}
		})
	}
}
