package resource_test

import (
	"context"
	"fmt"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
)

type fakeResourceService struct {
	ShouldError  bool
	CreateCalled bool
	DeleteCalled bool
}

func (f *fakeResourceService) CreateResource(ctx context.Context, r *resources.Resource) (*resources.Resource, error) {
	f.CreateCalled = true
	if f.ShouldError {
		return nil, fmt.Errorf("simulated failure")
	}
	return &resources.Resource{ID: "test-id"}, nil
}

func (f *fakeResourceService) DeleteResource(ctx context.Context, id string) error {
	f.DeleteCalled = true
	if f.ShouldError {
		return fmt.Errorf("simulated delete failure")
	}
	return nil
}
