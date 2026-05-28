package resource_test

import (
	"context"
	"testing"
	"xo/src/bundle"
	"xo/src/resource"

	"github.com/stretchr/testify/require"
)

func TestPublish(t *testing.T) {
	type testData struct {
		name    string
		service *fakeResourceService
		wantErr bool
	}
	tests := []testData{
		{
			name:    "success",
			service: &fakeResourceService{ShouldError: false},
			wantErr: false,
		},
		{
			name:    "failure",
			service: &fakeResourceService{ShouldError: true},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resourceMap := map[string]any{"foo": "bar"}
			bun := &bundle.Bundle{Artifacts: map[string]interface{}{"properties": map[string]interface{}{"foobar": map[string]interface{}{"$ref": "massdriver/resource-type"}}}}
			err := resource.Publish(context.Background(), tc.service, resourceMap, bun, "foobar", "resourceName")

			require.True(t, tc.service.CreateCalled, "expected CreateResource to be called")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
