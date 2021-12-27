package massdriver_test

import (
	"testing"
	"xo/src/massdriver"
)

func TestUploadArtifactFile(t *testing.T) {

	type testData struct {
		name      string
		id        string
		artifacts []map[string]interface{}
		want      string
	}
	tests := []testData{
		{
			name:      "Test Artifact Update",
			id:        "id",
			artifacts: []map[string]interface{}{{"foo": map[string]interface{}{"bar": "baz"}}, {"hello": "world"}},
			want:      `{"metadata":{"timestamp":"2021-01-01 12:00:00.1234","provisioner":"testaform","event_type":"artifact_update"},"payload":{"deployment_id":"id","artifacts":[{"foo":{"bar":"baz"}},{"hello":"world"}]}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			testSNSClient := SNSTestClient{}
			testMassdriverClient := massdriver.MassdriverClient{SNSClient: &testSNSClient}
			err := testMassdriverClient.UploadArtifact(tc.artifacts, tc.id)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			got := *testSNSClient.Data
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}
