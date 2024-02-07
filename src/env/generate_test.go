package env_test

import (
	"reflect"
	"testing"
	"xo/src/env"
)

func TestGenerateEnvs(t *testing.T) {
	type testData struct {
		name        string
		envs        map[string]string
		params      map[string]interface{}
		connections map[string]interface{}
		want        map[string]string
	}
	tests := []testData{
		{
			name: "basic",
			envs: map[string]string{
				"foo":              `@text "bar"`,
				"params_test":      `.params.nested.value`,
				"connections_test": `.connections.artifact.data.infrastructure.some`,
			},
			params: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": "something",
				},
			},
			connections: map[string]interface{}{
				"artifact": map[string]interface{}{
					"data": map[string]interface{}{
						"infrastructure": map[string]interface{}{
							"some": "field",
						},
					},
				},
			},
			want: map[string]string{
				"foo":              "bar",
				"params_test":      "something",
				"connections_test": "field",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := env.GenerateEnvs(tc.envs, tc.params, tc.connections)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got: %+v, want %+v", got, tc.want)
			}
		})
	}
}
