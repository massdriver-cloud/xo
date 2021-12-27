package terraform

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"xo/src/massdriver"

	"github.com/rs/zerolog/log"

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

func ReportProgressFromLogs(client *massdriver.MassdriverClient, deploymentId string, stream io.Reader) error {
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

		event, err := convertLogToMassdriverEvent(&record, deploymentId)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while parsing status message")
		}

		if event != nil {
			err = client.PublishEventToSNS(event)
			if err != nil {
				log.Error().Err(err).Msg("an error occurred while sending resource status to massdriver")
			}
		}
	}

	return nil
}

func convertLogToMassdriverEvent(record *terraformLog, deploymentId string) (*massdriver.Event, error) {
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
		event, err = parseDiagnosticLog(record, deploymentId)
		if err != nil {
			return nil, err
		}
	case "planned_change", "apply_start", "apply_complete", "apply_errored", "resource_drift":
		event, err = parseResourceUpdateLog(record, deploymentId)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	event.Metadata.Version = terraformVersion

	return event, nil
}

func parseDiagnosticLog(record *terraformLog, deploymentId string) (*massdriver.Event, error) {
	if record.Diagnostic == nil {
		return nil, errors.New("diagnostic struct missing")
	}

	event := massdriver.NewEvent(massdriver.EVENT_TYPE_PROVISIONER_ERROR)
	diagnostic := new(massdriver.EventPayloadDiagnostic)
	diagnostic.DeploymentId = deploymentId
	diagnostic.Message = record.Diagnostic.Summary

	switch record.Diagnostic.Severity {
	case "error":
		diagnostic.Level = "error"
	case "warning":
		diagnostic.Level = "warning"
	default:
		return nil, errors.New("unknown severity: " + record.Diagnostic.Severity)
	}

	event.Payload = diagnostic

	return event, nil
}

func parseResourceUpdateLog(record *terraformLog, deploymentId string) (*massdriver.Event, error) {

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

	event.Payload = progress

	return event, nil
}
