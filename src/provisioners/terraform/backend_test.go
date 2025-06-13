package terraform

import (
	"testing"
	"xo/src/massdriver"

	"github.com/stretchr/testify/require"
)

func TestGenerateJSONBackendHTTPConfig(t *testing.T) {
	spec := massdriver.Specification{
		DeploymentID: "depId",
		Token:        "token",
		PackageName:  "pkg-id-long-0000",
		URL:          "https://foo.massdriver.cloud",
	}
	got, _ := GenerateJSONBackendHTTPConfig(&spec, "step")
	want := `
	{
		"terraform": {
			"backend": {
				"http": {
					"username": "depId",
					"password": "token",
					"address": "https://foo.massdriver.cloud/state/pkg-id-long/step",
					"lock_address": "https://foo.massdriver.cloud/state/pkg-id-long/step",
					"unlock_address": "https://foo.massdriver.cloud/state/pkg-id-long/step"
				}
			}
		}
	}
`

	require.JSONEq(t, string(got), want)
}
