package resource

import (
	"context"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
)

type ResourceService interface {
	CreateResource(ctx context.Context, input *resources.ResourceInput) (*resources.Resource, error)
	DeleteResource(ctx context.Context, id string) error
}
