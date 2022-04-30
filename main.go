package main

import (
	"os"
	"xo/cmd"

	"github.com/lightstep/otel-launcher-go/launcher"
)

func main() {
	// Setup Tracing
	if os.Getenv("LS_ACCESS_TOKEN") != "" {
		otelLauncher := launcher.ConfigureOpentelemetry(
			launcher.WithServiceName("xo"),
			launcher.WithAccessToken(os.Getenv("LS_ACCESS_TOKEN")),
			launcher.WithResourceAttributes(map[string]string{
				"massdriver.deployment_id":   os.Getenv("MASSDRIVER_DEPLOYMENT_ID"),
				"massdriver.organization_id": os.Getenv("MASSDRIVER_ORGANIZATION_ID"),
				"massdriver.package_id":      os.Getenv("MASSDRIVER_PACKAGE_ID"),
				"massdriver.bundle_access":   os.Getenv("MASSDRIVER_BUNDLE_ACCESS"),
				"massdriver.bundle_name":     os.Getenv("MASSDRIVER_BUNDLE_NAME"),
			}),
		)
		defer otelLauncher.Shutdown()
	}

	// Run application
	cmd.Execute()
}
