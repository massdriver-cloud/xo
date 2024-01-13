package artifact_test

import (
	"bytes"
	"os"
	"testing"
	"xo/src/artifact"
)

func TestValidate(t *testing.T) {
	type testData struct {
		name       string
		field      string
		schemaPath string
		artifact   string
		want       bool
	}
	tests := []testData{
		{
			name:       "pass",
			field:      "one",
			schemaPath: "testdata/schema-artifacts.json",
			artifact:   `{"data":{"foo":{"bar":"baz"}},"specs":{"hello":"world"}}`,
			want:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// t.Setenv("MASSDRIVER_PROVISIONER", "testaform")
			// massdriver.EventTimeString = func() string { return "2021-01-01 12:00:00.1234" }
			// testClient := testmass.NewMassdriverTestClient(tc.deploymentId)

			schemasBytes, err := os.ReadFile(tc.schemaPath)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			schemasBuffer := bytes.NewBuffer(schemasBytes)

			artifactBuffer := bytes.NewBufferString(tc.artifact)

			//input := bytes.NewBuffer([]byte(tc.artifact))
			got, err := artifact.Validate(tc.field, artifactBuffer, schemasBuffer)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}

			// got := testClient.GetSNSMessages()
			// if got[0] != tc.want {
			// 	t.Fatalf("want: %v, got: %v", tc.want, got[0])
			// }
		})
	}
}
