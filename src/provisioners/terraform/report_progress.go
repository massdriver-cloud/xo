package terraform

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"xo/src/massdriver"
	"xo/src/telemetry"
	"xo/src/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/zclconf/go-cty/cty"
	ctyconvert "github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var terraformVersion string

type terraformLog struct {
	Level      string               `json:"@level"`
	Message    string               `json:"@message"`
	Module     string               `json:"@module"`
	Timestamp  string               `json:"@timestamp"`
	Changes    *terraformChanges    `json:"changes,omitempty"`
	Hook       *terraformAction     `json:"hook,omitempty"`
	Change     *terraformAction     `json:"change,omitempty"`
	Diagnostic *terraformDiagnostic `json:"diagnostic,omitempty"`
	Type       string               `json:"type"`
	Terraform  string               `json:"terraform,omitempty"`
}

type terraformDiagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Address  string `json:"address"`
}

type terraformChanges struct {
	Add       int    `json:"add"`
	Change    int    `json:"change"`
	Remove    int    `json:"remove"`
	Operation string `json:"operation"`
}

type terraformAction struct {
	Resource terraformResourceAddr `json:"resource"`
	Action   string                `json:"action"`
	IDKey    string                `json:"id_key,omitempty"`
	IDValue  string                `json:"id_value,omitempty"`
	Elapsed  float64               `json:"elapsed_seconds"`
}

type terraformResourceAddr struct {
	Addr            string                  `json:"addr"`
	Module          string                  `json:"module"`
	Resource        string                  `json:"resource"`
	ImpliedProvider string                  `json:"implied_provider"`
	ResourceType    string                  `json:"resource_type"`
	ResourceName    string                  `json:"resource_name"`
	ResourceKey     ctyjson.SimpleJSONValue `json:"resource_key"`
}

func ReportProgressFromLogs(ctx context.Context, client *massdriver.MassdriverClient, deploymentId string, stream io.Reader) error {
	_, span := otel.Tracer("xo").Start(ctx, "provisioners.terraform.ReportProgressFromLogs")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		// print the log to the console so its still visible
		fmt.Println(scanner.Text())

		var record terraformLog

		err := json.Unmarshal([]byte(scanner.Text()), &record)
		if err != nil {
			// maybe warn?
			continue
		}

		// normally we'd pass context instead of a span (probably an antipattern) but for readability
		// its much easier to annotate the same span, so we pass that instead
		event, err := convertLogToMassdriverEvent(span, &record, deploymentId)
		if err != nil {
			util.LogError(err, span, "an error occurred while parsing status message")
		}

		if event != nil {
			err = client.PublishEventToSNS(event)
			if err != nil {
				util.LogError(err, span, "an error occurred while sending resource status to massdriver")
			}
		}
	}

	return nil
}

func convertLogToMassdriverEvent(span trace.Span, record *terraformLog, deploymentId string) (*massdriver.Event, error) {
	if record.Terraform != "" {
		terraformVersion = record.Terraform
		return nil, nil
	}

	var event *massdriver.Event
	var err error

	switch record.Type {
	case "change_summary":
		// skipping change_summary events for now
		return nil, nil
	case "diagnostic":
		event, err = parseDiagnosticLog(span, record, deploymentId)
		if err != nil {
			return nil, err
		}
	case "planned_change", "apply_start", "apply_complete", "apply_errored", "resource_drift":
		event, err = parseResourceUpdateLog(span, record, deploymentId)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	if event != nil {
		event.Metadata.Version = terraformVersion
	}

	return event, nil
}

func parseDiagnosticLog(span trace.Span, record *terraformLog, deploymentId string) (*massdriver.Event, error) {
	if record.Diagnostic == nil {
		return nil, errors.New("diagnostic struct missing")
	}

	event := massdriver.NewEvent(massdriver.EVENT_TYPE_PROVISIONER_ERROR)
	diagnostic := new(massdriver.EventPayloadDiagnostic)
	diagnostic.DeploymentId = deploymentId
	diagnostic.Message = record.Diagnostic.Summary
	diagnostic.Details = record.Diagnostic.Detail

	switch record.Diagnostic.Severity {
	case "error":
		diagnostic.Level = "error"
		terraformError := errors.New(diagnostic.Message + " Details: " + diagnostic.Details)
		span.RecordError(terraformError)
		span.SetStatus(codes.Error, terraformError.Error())
	case "warning":
		diagnostic.Level = "warning"
		span.AddEvent("warning", trace.WithAttributes(
			attribute.String("message", diagnostic.Message+" Details: "+diagnostic.Details),
		))
	default:
		return nil, errors.New("unknown severity: " + record.Diagnostic.Severity)
	}

	event.Payload = diagnostic

	return event, nil
}

func parseResourceUpdateLog(span trace.Span, record *terraformLog, deploymentId string) (*massdriver.Event, error) {
	var action *terraformAction
	if record.Hook != nil {
		action = record.Hook
	} else if record.Change != nil {
		action = record.Change
	} else {
		return nil, errors.New("resource update without resource data")
	}

	var eventType string
	// we build the event type here by combining strings. Its cleaner than a massive case statement
	// and types are checked/enforce through tests
	switch action.Action {
	case "create":
		eventType = "create_"
	case "update":
		eventType = "update_"
	case "delete":
		eventType = "delete_"
	case "replace":
		eventType = "recreate_"
	case "read":
		return nil, nil // squelching these for now, I think these only happen on data lookups
	default:
		return nil, errors.New("unknown action: " + action.Action)
	}

	switch record.Type {
	case "planned_change":
		eventType += "pending"
	case "apply_start":
		eventType += "running"
	case "apply_complete":
		eventType += "completed"
	case "apply_errored":
		eventType += "failed"
	case "resource_drift":
		eventType = "drift_detected"
	default:
		return nil, errors.New("unknown type: " + record.Type)
	}

	event := massdriver.NewEvent(eventType)

	progress := new(massdriver.EventPayloadResourceProgress)
	progress.DeploymentId = deploymentId
	progress.ResourceName = action.Resource.ResourceName
	progress.ResourceType = action.Resource.ResourceType
	progress.ResourceId = action.IDValue
	// Convert ResourceKey (which can be a string or an int) into a string if the value is non-null
	if !action.Resource.ResourceKey.IsNull() {
		key, err := ctyconvert.Convert(action.Resource.ResourceKey.Value, cty.String)
		if err != nil {
			return nil, err
		}
		progress.ResourceKey = key.AsString()
	}

	span.AddEvent("resource-update", trace.WithAttributes(
		attribute.String("resource-status", eventType),
		attribute.String("resource-id", progress.ResourceId),
		attribute.String("resource-name", progress.ResourceName),
		attribute.String("resource-type", progress.ResourceType),
	))

	event.Payload = progress

	return event, nil
}
