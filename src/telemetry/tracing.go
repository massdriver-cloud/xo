package telemetry

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func GetContextWithTraceParentFromEnv() context.Context {
	return GetContextWithTraceParent(os.Getenv("TRACEPARENT"))
}

func GetContextWithTraceParent(traceparent string) context.Context {
	carrier := propagation.MapCarrier{}
	carrier.Set("traceparent", traceparent)
	return otel.GetTextMapPropagator().Extract(context.Background(), carrier)
}

func GetTraceParentFromContext(ctx context.Context) string {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier.Get("traceparent")
}

func SetSpanAttributes(span trace.Span) {
	span.SetAttributes(
		attribute.String("massdriver.deployment_id", os.Getenv("MASSDRIVER_DEPLOYMENT_ID")),
		attribute.String("massdriver.bundle_id", os.Getenv("MASSDRIVER_BUNDLE_ID")),
		attribute.String("massdriver.organization_id", os.Getenv("MASSDRIVER_ORGANIZATION_ID")),
		attribute.String("massdriver.package_id", os.Getenv("MASSDRIVER_PACKAGE_ID")),
		attribute.String("massdriver.package_name", os.Getenv("MASSDRIVER_PACKAGE_NAME")),
		attribute.String("massdriver.bundle_access", os.Getenv("MASSDRIVER_BUNDLE_ACCESS")),
		attribute.String("massdriver.bundle_name", os.Getenv("MASSDRIVER_BUNDLE_NAME")),
	)
}

func LogError(span trace.Span, err error, message string) error {
	log.Error().Err(err).Msg(message)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	return err
}
