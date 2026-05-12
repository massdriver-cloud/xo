package terraform

import (
	"testing"
	"xo/src/massdriver"

	"github.com/stretchr/testify/require"
)

func TestGenerateJSONBackendHTTPConfig(t *testing.T) {
	tests := []struct {
		name string
		spec massdriver.Specification
		step string
		want string
	}{
		{
			name: "generates backend from Package Name",
			spec: massdriver.Specification{
				DeploymentID: "depId",
				Token:        "token",
				PackageName:  "pkg-id-long-0000",
				URL:          "https://foo.massdriver.cloud",
			},
			step: "step",
			want: `{
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
			}`,
		}, {
			name: "generates backend config from Instance ID",
			spec: massdriver.Specification{
				DeploymentID: "depId",
				Token:        "token",
				InstanceID:   "inst-ance-id",
				URL:          "https://foo.massdriver.cloud",
			},
			step: "step",
			want: `{
				"terraform": {
					"backend": {
						"http": {
							"username": "depId",
							"password": "token",
							"address": "https://foo.massdriver.cloud/state/inst-ance-id/step",
							"lock_address": "https://foo.massdriver.cloud/state/inst-ance-id/step",
							"unlock_address": "https://foo.massdriver.cloud/state/inst-ance-id/step"
						}
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GenerateJSONBackendHTTPConfig(&tt.spec, tt.step)
			require.JSONEq(t, string(got), tt.want)
		})
	}
}
