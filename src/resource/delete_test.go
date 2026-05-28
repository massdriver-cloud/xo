package resource_test

import (
	"context"
	"testing"
	"xo/src/resource"

	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	type testData struct {
		name    string
		service *fakeResourceService
		id      string
		field   string
		wantErr bool
	}

	tests := []testData{
		{
			name: "basic delete success",
			service: &fakeResourceService{
				ShouldError: false,
			},
			id:      "artId",
			field:   "foobar",
			wantErr: false,
		},
		{
			name: "delete failure",
			service: &fakeResourceService{
				ShouldError: true,
			},
			id:      "artId",
			field:   "foobar",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := resource.Delete(context.Background(), tc.service, tc.id)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.True(t, tc.service.DeleteCalled, "expected DeleteResource to be called")
		})
	}
}
