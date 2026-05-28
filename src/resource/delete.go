package resource

import (
	"context"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
)

func Delete(ctx context.Context, svc ResourceService, id, field string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ResourceDelete")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	return svc.DeleteResource(ctx, id)
}
