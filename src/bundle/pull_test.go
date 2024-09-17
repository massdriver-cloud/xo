package bundle_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"
	"xo/src/bundle"
	"xo/src/massdriver"

	"github.com/massdriver-cloud/mass/pkg/gqlmock"
)

func TestPull(t *testing.T) {
	type testData struct {
		name string
		data []byte
	}
	tests := []testData{
		{
			name: "basic",
			data: []byte(`data`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gqlClient := gqlmock.NewClientWithSingleJSONResponse(map[string]interface{}{
				"data": map[string]interface{}{
					"bundleSourceCode": map[string]interface{}{
						"source": base64.StdEncoding.EncodeToString(tc.data),
					},
				},
			})

			client := massdriver.MassdriverClient{
				GQLCLient: gqlClient,
				Specification: &massdriver.Specification{
					BundleID:       "bundleuuid1",
					OrganizationID: "orguuid1",
				},
			}

			buf := new(bytes.Buffer)

			pullErr := bundle.Pull(context.Background(), &client, buf)
			if pullErr != nil {
				t.Errorf("pull failed: %v", pullErr)
			}

			got := buf.String()
			want := string(tc.data)
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}

}
