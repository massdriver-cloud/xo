package resource

import (
	"context"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func Publish(ctx context.Context, svc ResourceService, resourceMap map[string]any, bun *bundle.Bundle, field, name string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ResourcePublish")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	resourceType, err := getResourceTypeFromBundle(bun, field)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resource := resources.Resource{}
	resource.Name = name
	resource.Field = field
	resource.Type = resourceType
	resource.Payload = resourceMap

	_, createErr := svc.CreateResource(ctx, &resource)
	if createErr != nil {
		span.RecordError(createErr)
		span.SetStatus(codes.Error, createErr.Error())
		return createErr
	}

	return nil
}
