package deployment_test

import (
	"context"
	"fmt"
	"testing"
	"xo/src/deployment"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/deployments"
)

type fakeDeploymentService struct {
	Status deployments.Status
	ID     string
	Err    error
}

func (f *fakeDeploymentService) UpdateDeploymentStatus(ctx context.Context, id string, status deployments.Status) (*deployments.Deployment, error) {
	f.ID = id
	f.Status = status
	if f.Err != nil {
		return nil, f.Err
	}
	return &deployments.Deployment{}, nil
}

func TestUpdateDeploymentStatus(t *testing.T) {
	type testData struct {
		name        string
		id          string
		err         error
		shouldError bool
		status      deployments.Status
	}

	tests := []testData{
		{
			name:        "basic update",
			id:          "test-id",
			status:      deployments.StatusRunning,
			shouldError: false,
		},
		{
			name:        "invalid transition",
			id:          "test-id",
			status:      deployments.StatusFailed,
			err:         &deployments.InvalidTransitionError{},
			shouldError: false,
		},
		{
			name:        "update with error",
			id:          "test-id",
			status:      deployments.StatusRunning,
			err:         fmt.Errorf("simulated error"),
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeService := &fakeDeploymentService{Err: tc.err}
			err := deployment.UpdateDeploymentStatus(context.Background(), fakeService, tc.id, tc.status)

			if err != nil && !tc.shouldError {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tc.shouldError {
				t.Errorf("expected error, got nil")
			}

			if fakeService.ID != tc.id {
				t.Errorf("expected ID %s, got %s", tc.id, fakeService.ID)
			}

			if fakeService.Status != tc.status {
				t.Errorf("expected status %s, got %s", tc.status, fakeService.Status)
			}
		})
	}

}
