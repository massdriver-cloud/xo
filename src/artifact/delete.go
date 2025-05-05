package artifact

import (
	"context"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
)

func Delete(ctx context.Context, svc ArtifactService, id, field string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ArtifactDelete")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	return svc.DeleteArtifact(ctx, id, field)
}
