package main

import (
	"context"
	"log"
	"os"
	"xo/cmd"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/credentials"
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
		otlptracegrpc.WithHeaders(map[string]string{
			"x-honeycomb-team":    os.Getenv("HONEYCOMB_API_KEY"),
			"x-honeycomb-dataset": os.Getenv("HONEYCOMB_DATASET"),
		}),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}

	client := otlptracegrpc.NewClient(opts...)
	return otlptrace.New(ctx, client)
}

func newTraceProvider(exp *otlptrace.Exporter) *sdktrace.TracerProvider {
	// The service.name attribute is required.
	resource :=
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("xo"),
			attribute.String("MASSDRIVER_DEPLOYMENT_ID", os.Getenv("MASSDRIVER_DEPLOYMENT_ID")),
			attribute.String("MASSDRIVER_ORGANIZATION_ID", os.Getenv("MASSDRIVER_ORGANIZATION_ID")),
			attribute.String("MASSDRIVER_PACKAGE_ID", os.Getenv("MASSDRIVER_PACKAGE_ID")),
			attribute.String("MASSDRIVER_BUNDLE_ACCESS", os.Getenv("MASSDRIVER_BUNDLE_ACCESS")),
			attribute.String("MASSDRIVER_BUNDLE_NAME", os.Getenv("MASSDRIVER_BUNDLE_NAME")),
		)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
	)
}

func main() {
	// Setup Tracing
	// TODO: this is a poor check but it works for now
	if os.Getenv("HONEYCOMB_API_KEY") != "" && os.Getenv("HONEYCOMB_DATASET") != "" {
		ctx := context.Background()

		// Configure a new exporter using environment variables for sending data to Honeycomb over gRPC.
		exp, err := newExporter(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new tracer provider with a batch span processor and the otlp exporter.
		tp := newTraceProvider(exp)

		// Handle this error in a sensible manner where possible
		defer func() { _ = tp.Shutdown(ctx) }()

		// Set the Tracer Provider and the W3C Trace Context propagator as globals
		otel.SetTracerProvider(tp)

		// Register the trace context and baggage propagators so data is propagated across services/processes.
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
	}

	// Run application
	cmd.Execute()
}
