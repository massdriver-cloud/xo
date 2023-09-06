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

var creatingRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Creating ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)\.{3}$`
var readingRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Looks like there are no changes for ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`
var patchingRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Patching ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`
var deletingRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Deleting ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`
var completedRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Status ([a-zA-Z0-9]+) is ready: ([a-zA-Z0-9._\/-]+)$`
var deletedRegex = `^[a-zA-Z0-9]+\.go:[0-9]+: \[debug\] Deleted ([a-zA-Z0-9]+) ([a-zA-Z0-9._\/-]+)$`
var errorRegex = `^Error: (.*)$`

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
			err = client.PublishEventToSNS(event)
			if err != nil {
				util.LogError(err, span, "an error occurred while sending resource status to massdriver")
			}
		}
	}

	return nil
}

func convertLogToMassdriverEvent(span trace.Span, logLine, deploymentId string) (*massdriver.Event, error) {

	matchErrorRegex, err := regexp.Match(errorRegex, []byte(logLine))
	if matchErrorRegex && err == nil {
		return parseErrorLog(span, logLine, deploymentId)
	} else {
		return parseResourceUpdateLog(span, logLine, deploymentId)
	}
}

func parseErrorLog(span trace.Span, logLine, deploymentId string) (*massdriver.Event, error) {

	re, err := regexp.Compile(errorRegex)
	if err != nil {
		return nil, err
	}
	result := re.FindStringSubmatch(logLine)

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

	eventToRegex := map[string]string{
		creatingRegex:  massdriver.EVENT_TYPE_RESOURCE_CREATE_RUNNING,
		patchingRegex:  massdriver.EVENT_TYPE_RESOURCE_UPDATE_RUNNING,
		deletingRegex:  massdriver.EVENT_TYPE_RESOURCE_DELETE_RUNNING,
		completedRegex: massdriver.EVENT_TYPE_RESOURCE_CREATE_COMPLETED,
		deletedRegex:   massdriver.EVENT_TYPE_RESOURCE_DELETE_COMPLETED,
	}

	for pattern, eventType := range eventToRegex {
		matchRegex, _ := regexp.Match(pattern, []byte(logLine))
		if !matchRegex {
			continue
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		results := re.FindStringSubmatch(logLine)

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
