package bundle

import (
	"context"
	"xo/src/api"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
)

func Pull(ctx context.Context, client *massdriver.MassdriverClient) error {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePull")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	bundleBytes, err := api.GetBundleSourceCode(client.GQLCLient, client.Specification.BundleID, client.Specification.OrganizationID)

	// defer outFile.Close()
	// _, err = io.Copy(outFile, obj.Body)
	// if err != nil {
	// 	span.RecordError(err)
	// 	span.SetStatus(codes.Error, err.Error())
	// 	return err
	// }

	return nil
}
