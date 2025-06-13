package bundle

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"xo/src/api"
	"xo/src/massdriver"
	"xo/src/telemetry"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	oras "oras.land/oras-go/v2"
)

func PullV0(ctx context.Context, client *massdriver.MassdriverClient, outFile io.Writer) error {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	bundleBytes, getErr := api.GetBundleSourceCode(client.GQLCLient, client.Specification.BundleID, client.Specification.OrganizationUUID)
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

func PullV1(ctx context.Context, repo oras.Target, target oras.Target, tag string) (v1.Descriptor, error) {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePullV1")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	return oras.Copy(ctx, repo, tag, target, tag, oras.DefaultCopyOptions)
}
