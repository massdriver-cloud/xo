package artifact

import (
	"context"
	"crypto/sha256"
	"fmt"
	"xo/src/bundle"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func Delete(ctx context.Context, c *massdriver.MassdriverClient, bun *bundle.Bundle, field, name string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ArtifactDelete")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	artifactType, err := getArtifactTypeFromBundle(bun, field)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Just hashing the field name for now. Eventually this will either be moved or repurposed.
	providerResourceId := fmt.Sprintf("%x", sha256.Sum256([]byte(field)))
	metadata := artifactMetadata{
		Field:              field,
		ProviderResourceID: providerResourceId,
		Type:               artifactType,
		Name:               name,
	}

	// the only thing the API cares about during a delete is the metadata block
	unmarshaledArtifact := map[string]interface{}{}
	unmarshaledArtifact["metadata"] = metadata
	massdriver.DeleteArtifact(c, unmarshaledArtifact)

	return nil
}
