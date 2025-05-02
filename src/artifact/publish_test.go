package artifact_test

import (
	"context"
	"fmt"
	"testing"
	"xo/src/artifact"
	"xo/src/bundle"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/artifacts"
	"github.com/stretchr/testify/require"
)

func (f *fakeArtifactService) CreateArtifact(ctx context.Context, a *artifacts.Artifact) (*artifacts.Artifact, error) {
	f.CreateCalled = true
	if f.ShouldError {
		return nil, fmt.Errorf("simulated failure")
	}
	return &artifacts.Artifact{ID: "test-id"}, nil
}

func TestPublish(t *testing.T) {
	type testData struct {
		name    string
		service *fakeArtifactService
		wantErr bool
	}
	tests := []testData{
		{
			name:    "success",
			service: &fakeArtifactService{ShouldError: false},
			wantErr: false,
		},
		{
			name:    "failure",
			service: &fakeArtifactService{ShouldError: true},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			artifactBytes := []byte(`{"foo":"bar"}`)
			bun := &bundle.Bundle{Artifacts: map[string]interface{}{"properties": map[string]interface{}{"foobar": map[string]interface{}{"$ref": "massdriver/artifact-type"}}}}
			err := artifact.Publish(context.Background(), tc.service, artifactBytes, bun, "foobar", "artName")

			require.True(t, tc.service.CreateCalled, "expected DeleteArtifact to be called")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
