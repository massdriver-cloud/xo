package helm

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"xo/src/massdriver"
	"xo/src/telemetry"
	"xo/src/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	creatingRegex  = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Creating ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)\.{3}$`)
	readingRegex   = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Looks like there are no changes for ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`)
	patchingRegex  = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Patching ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`)
	deletingRegex  = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Deleting ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`)
	completedRegex = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Status ([a-zA-Z0-9]+) is ready: ([a-zA-Z0-9._\/-]+)$`)
	deletedRegex   = regexp.MustCompile(`^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Deleted ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`)
	errorRegex     = regexp.MustCompile(`^Error: (.*)$`)
)

func ReportProgressFromLogs(ctx context.Context, client *massdriver.MassdriverClient, deploymentId string, stream io.Reader) error {
	_, span := otel.Tracer("xo").Start(ctx, "provisioners.helm.ReportProgressFromLogs")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	helmVersion := "unknown"
	if version, ok := os.LookupEnv("HELM_VERSION"); ok {
		helmVersion = version
	}

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		// print the log to the console so its still visible
		fmt.Println(scanner.Text())

		// normally we'd pass context instead of a span (probably an antipattern) but for readability
		// its much easier to annotate the same span, so we pass that instead
		event, err := convertLogToMassdriverEvent(span, scanner.Text(), deploymentId)
		if err != nil {
			util.LogError(err, span, "an error occurred while parsing status message")
		}

		if event != nil {
			event.Metadata.Version = helmVersion
			err = client.PublishEvent(event)
			if err != nil {
				util.LogError(err, span, "an error occurred while sending resource status to massdriver")
			}
		}
	}

	return nil
}

func convertLogToMassdriverEvent(span trace.Span, logLine, deploymentId string) (*massdriver.Event, error) {

	if errorRegex.Match([]byte(logLine)) {
		return parseErrorLog(span, logLine, deploymentId)
	} else {
		return parseResourceUpdateLog(span, logLine, deploymentId)
	}
}

func parseErrorLog(span trace.Span, logLine, deploymentId string) (*massdriver.Event, error) {

	result := errorRegex.FindStringSubmatch(logLine)

	helmError := result[1]

	event := massdriver.NewEvent(massdriver.EVENT_TYPE_PROVISIONER_ERROR)
	diagnostic := new(massdriver.EventPayloadDiagnostic)
	diagnostic.DeploymentId = deploymentId
	diagnostic.Message = helmError
	diagnostic.Level = "error"
	span.RecordError(errors.New(helmError))
	span.SetStatus(codes.Error, helmError)

	event.Payload = diagnostic

	return event, nil
}

func parseResourceUpdateLog(span trace.Span, logLine, deploymentId string) (*massdriver.Event, error) {

	eventToRegex := map[string]*regexp.Regexp{
		massdriver.EVENT_TYPE_RESOURCE_CREATE_RUNNING:   creatingRegex,
		massdriver.EVENT_TYPE_RESOURCE_UPDATE_RUNNING:   patchingRegex,
		massdriver.EVENT_TYPE_RESOURCE_DELETE_RUNNING:   deletingRegex,
		massdriver.EVENT_TYPE_RESOURCE_CREATE_COMPLETED: completedRegex,
		massdriver.EVENT_TYPE_RESOURCE_DELETE_COMPLETED: deletedRegex,
	}

	for eventType, pattern := range eventToRegex {
		if !pattern.Match([]byte(logLine)) {
			continue
		}

		results := pattern.FindStringSubmatch(logLine)

		resourceType := results[1]
		resourceName := strings.TrimPrefix(results[2], "/")

		// TODO remove this trim-prefix when massdriver server can accept it
		event := massdriver.NewEvent(strings.TrimPrefix(eventType, "resource_"))

		progress := new(massdriver.EventPayloadResourceProgress)
		progress.DeploymentId = deploymentId
		progress.ResourceName = resourceName
		progress.ResourceType = resourceType
		progress.ResourceId = fmt.Sprintf("%s:%s", resourceType, resourceName)

		span.AddEvent("resource-update", trace.WithAttributes(
			attribute.String("resource-status", eventType),
			attribute.String("resource-id", progress.ResourceId),
			attribute.String("resource-name", progress.ResourceName),
			attribute.String("resource-type", progress.ResourceType),
		))

		event.Payload = progress

		return event, nil
	}

	return nil, nil
}
