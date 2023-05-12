package main

import (
	"log"
	"os"
	"xo/cmd"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	// Setup Tracing
	if os.Getenv("HONEYCOMB_API_KEY") != "" {
		bsp := honeycomb.NewBaggageSpanProcessor()

		// use honeycomb distro to setup OpenTelemetry SDK
		otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
			otelconfig.WithSpanProcessor(bsp),
		)
		if err != nil {
			log.Fatalf("error setting up OTel SDK - %e", err)
		}
		defer otelShutdown()
	}

	// Run application
	if err := cmd.Execute(); err != nil {
		exitCode = 1
		return
	}
}
