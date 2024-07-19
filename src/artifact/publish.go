package artifact

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"xo/src/bundle"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func Publish(ctx context.Context, c *massdriver.MassdriverClient, artifact []byte, bun *bundle.Bundle, field, name string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ArtifactPublish")
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

	// this here is a bit clunky. We're nesting the metadata object WITHIN the artifact. However, the schemas don't expect
	// the metadata block. So after validation we need to unmarshal the artifact to a map so we can add the metadata in
	var unmarshaledArtifact map[string]interface{}
	err = json.Unmarshal(artifact, &unmarshaledArtifact)
	if err != nil {
		fmt.Println(string(artifact))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	unmarshaledArtifact["metadata"] = metadata

	massdriver.PublishArtifact(c, unmarshaledArtifact)

	return nil
}
