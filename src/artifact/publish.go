package artifact

import (
	"context"
	"encoding/json"
	"fmt"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/artifacts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func Publish(ctx context.Context, svc ArtifactService, artifactBytes []byte, bun *bundle.Bundle, field, name string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ArtifactPublish")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	artifactType, err := getArtifactTypeFromBundle(bun, field)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	art := artifacts.Artifact{}
	err = json.Unmarshal(artifactBytes, &art)
	if err != nil {
		fmt.Println(string(artifactBytes))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	art.Name = name
	art.Field = field
	art.Type = artifactType

	_, createErr := svc.CreateArtifact(ctx, &art)
	if createErr != nil {
		span.RecordError(createErr)
		span.SetStatus(codes.Error, createErr.Error())
		return createErr
	}

	return nil
}
