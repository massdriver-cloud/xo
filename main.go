package main

import (
	"os"
	"xo/cmd"

	"github.com/lightstep/otel-launcher-go/launcher"
)

func main() {
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	// Setup Tracing
	if os.Getenv("LS_ACCESS_TOKEN") != "" {
		otelLauncher := launcher.ConfigureOpentelemetry(
			launcher.WithServiceName("xo"),
			launcher.WithPropagators([]string{"tracecontext"}),
		)
		defer otelLauncher.Shutdown()
	}

	// Run application
	if err := cmd.Execute(); err != nil {
		exitCode = 1
		return
	}
}
