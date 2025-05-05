package artifact_test

import (
	"context"
	"fmt"
	"testing"
	"xo/src/artifact"

	"github.com/stretchr/testify/require"
)

func (f *fakeArtifactService) DeleteArtifact(ctx context.Context, id, field string) error {
	f.DeleteCalled = true
	if f.ShouldError {
		return fmt.Errorf("simulated delete failure")
	}
	return nil
}

func TestDelete(t *testing.T) {
	type testData struct {
		name    string
		service *fakeArtifactService
		id      string
		field   string
		wantErr bool
	}

	tests := []testData{
		{
			name: "basic delete success",
			service: &fakeArtifactService{
				ShouldError: false,
			},
			id:      "artId",
			field:   "foobar",
			wantErr: false,
		},
		{
			name: "delete failure",
			service: &fakeArtifactService{
				ShouldError: true,
			},
			id:      "artId",
			field:   "foobar",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := artifact.Delete(context.Background(), tc.service, tc.id, tc.field)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.True(t, tc.service.DeleteCalled, "expected DeleteArtifact to be called")
		})
	}
}
