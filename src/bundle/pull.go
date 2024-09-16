package bundle

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"xo/src/api"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func Pull(ctx context.Context, client *massdriver.MassdriverClient, outFile io.Writer) error {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	bundleBytes, getErr := api.GetBundleSourceCode(client.GQLCLient, client.Specification.BundleID, client.Specification.OrganizationID)
	if getErr != nil {
		span.RecordError(getErr)
		span.SetStatus(codes.Error, getErr.Error())
		return fmt.Errorf("an error occurred while getting bundle source code: %w", getErr)
	}

	decodedBundleBytes, decodeErr := base64.StdEncoding.DecodeString(string(bundleBytes))
	if decodeErr != nil {
		span.RecordError(decodeErr)
		span.SetStatus(codes.Error, decodeErr.Error())
		return fmt.Errorf("an error occurred while decoding bundle source code: %w", decodeErr)
	}

	_, writeErr := outFile.Write(decodedBundleBytes)
	if writeErr != nil {
		span.RecordError(writeErr)
		span.SetStatus(codes.Error, writeErr.Error())
		return fmt.Errorf("an error occurred while writing bundle source code: %w", writeErr)
	}

	return nil
}
