package resource

import (
	"context"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Publish creates a resource for the current deployment. The resource type is
// resolved server-side from the deployment's release pin and field, so it is
// not derived from the bundle here.
func Publish(ctx context.Context, svc ResourceService, resourceMap map[string]any, field, name string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ResourcePublish")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	input := resources.ResourceInput{
		Field:   field,
		Name:    name,
		Payload: resourceMap,
	}

	_, createErr := svc.CreateResource(ctx, &input)
	if createErr != nil {
		span.RecordError(createErr)
		span.SetStatus(codes.Error, createErr.Error())
		return createErr
	}

	return nil
}
