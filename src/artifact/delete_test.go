package artifact_test

import (
	"context"
	"testing"
	"xo/src/artifact"
	"xo/src/bundle"
	"xo/src/massdriver"
	testmass "xo/test"
)

func TestDelete(t *testing.T) {
	type testData struct {
		name         string
		deploymentId string
		bun          bundle.Bundle
		field        string
		artifactName string
		want         string
	}
	tests := []testData{
		{
			name:         "basic",
			deploymentId: "depId",
			bun:          bundle.Bundle{Artifacts: map[string]interface{}{"properties": map[string]interface{}{"foobar": map[string]interface{}{"$ref": "massdriver/artifact-type"}}}},
			field:        "foobar",
			artifactName: "artName",
			want:         `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"artifact_deleted"},"payload":{"deployment_id":"depId","artifact":{"metadata":{"field":"foobar","provider_resource_id":"c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2","type":"massdriver/artifact-type","name":"artName"}}}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testClient := testmass.NewMassdriverTestClient(tc.deploymentId)

			err := artifact.Delete(context.TODO(), &testClient.MassClient, &tc.bun, tc.field, tc.artifactName)
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
