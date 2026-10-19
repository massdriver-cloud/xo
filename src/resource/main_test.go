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

func (f *fakeResourceService) CreateResource(ctx context.Context, input *resources.ResourceInput) (*resources.Resource, error) {
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

func (f *fakeResourceService) ListResources(ctx context.Context) ([]resources.Resource, error) {
	return nil, nil
}

// recordingResourceService records creates and deletes. Existing is what
// ListResources returns.
type recordingResourceService struct {
	Existing   []resources.Resource
	ListErr    error
	ListCalled bool
	Created    []resources.ResourceInput
	Deleted    []string
	CreateErr  error
	DeleteErr  map[string]error
}

func (r *recordingResourceService) CreateResource(ctx context.Context, input *resources.ResourceInput) (*resources.Resource, error) {
	if r.CreateErr != nil {
		return nil, r.CreateErr
	}
	r.Created = append(r.Created, *input)
	return &resources.Resource{ID: "inst." + input.Field}, nil
}

func (r *recordingResourceService) DeleteResource(ctx context.Context, id string) error {
	if err := r.DeleteErr[id]; err != nil {
		return err
	}
	r.Deleted = append(r.Deleted, id)
	return nil
}

func (r *recordingResourceService) ListResources(ctx context.Context) ([]resources.Resource, error) {
	r.ListCalled = true
	return r.Existing, r.ListErr
}
