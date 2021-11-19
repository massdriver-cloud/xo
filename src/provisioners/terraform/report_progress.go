package terraform

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"xo/src/massdriver"

	"github.com/rs/zerolog/log"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/zclconf/go-cty/cty"
	ctyconvert "github.com/zclconf/go-cty/cty/convert"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var ReportProgressSender func(*mdproto.ProvisionerProgressUpdateRequest) error

func sendToMassdriver(message *mdproto.ProvisionerProgressUpdateRequest) error {
	return massdriver.SendProvisionerProgressUpdate(message)
}

func init() {
	ReportProgressSender = sendToMassdriver
}

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

func ReportProgressFromLogs(deploymentId string, deploymentToken string, stream io.Reader) error {
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

		request, err := convertLogToProvisionerProgressUpdateRequest(&record)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while parsing status message")
		}

		if request != nil {
			request.DeploymentId = deploymentId
			request.DeploymentToken = deploymentToken

			err = ReportProgressSender(request)
			if err != nil {
				log.Error().Err(err).Msg("an error occurred while sending resource status to massdriver")
			}
		}
	}

	return nil
}

func convertLogToProvisionerProgressUpdateRequest(record *terraformLog) (*mdproto.ProvisionerProgressUpdateRequest, error) {
	var request mdproto.ProvisionerProgressUpdateRequest

	if record.Terraform != "" {
		terraformVersion = record.Terraform
		return nil, nil
	}

	switch record.Type {
	case "change_summary":
		err := parseChangeSummaryLog(record, &request)
		if err != nil {
			return nil, err
		}
	case "diagnostic":
		err := parseDiagnosticLog(record, &request)
		if err != nil {
			return nil, err
		}
	case "planned_change", "apply_start", "apply_complete", "apply_errored", "resource_drift":
		err := parseResourceUpdateLog(record, &request)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	request.Timestamp = record.Timestamp
	request.Metadata = &mdproto.ProvisionerMetadata{
		Provisioner:        mdproto.Provisioner_PROVISIONER_TERRAFORM,
		ProvisionerVersion: terraformVersion,
	}

	return &request, nil
}

func parseChangeSummaryLog(record *terraformLog, request *mdproto.ProvisionerProgressUpdateRequest) error {
	if record.Changes == nil {
		return errors.New("change summary without changes")
	}
	switch record.Changes.Operation {
	case "plan":
		request.Status = mdproto.ProvisionerStatus_PROVISIONER_STATUS_PLAN_COMPLETED
	case "apply":
		request.Status = mdproto.ProvisionerStatus_PROVISIONER_STATUS_APPLY_COMPLETED
	case "destroy":
		request.Status = mdproto.ProvisionerStatus_PROVISIONER_STATUS_DESTROY_COMPLETED
	default:
		return errors.New("unknown change_summary type: " + record.Changes.Operation)
	}
	return nil
}

func parseDiagnosticLog(record *terraformLog, request *mdproto.ProvisionerProgressUpdateRequest) error {
	var diagnostic mdproto.ProvisionerError

	if record.Diagnostic == nil {
		return errors.New("diagnostic struct missing")
	}

	switch record.Diagnostic.Severity {
	case "error":
		diagnostic.Level = mdproto.ProvisionerErrorLevel_PROVISIONER_ERROR_LEVEL_ERROR
	case "warning":
		diagnostic.Level = mdproto.ProvisionerErrorLevel_PROVISIONER_ERROR_LEVEL_WARNING
	default:
		return errors.New("unknown severity: " + record.Diagnostic.Severity)
	}

	request.Status = mdproto.ProvisionerStatus_PROVISIONER_STATUS_ERROR
	diagnostic.Message = record.Diagnostic.Summary

	request.Error = &diagnostic

	return nil
}

func parseResourceUpdateLog(record *terraformLog, request *mdproto.ProvisionerProgressUpdateRequest) error {

	var action *terraformAction
	if record.Hook != nil {
		action = record.Hook
	} else if record.Change != nil {
		action = record.Change
	} else {
		return errors.New("resource update without resource data")
	}

	var progress mdproto.ProvisionerResourceProgress
	switch action.Action {
	case "create":
		progress.Action = mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_CREATE
	case "update":
		progress.Action = mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_UPDATE
	case "delete":
		progress.Action = mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_DELETE
	case "replace":
		progress.Action = mdproto.ProvisionerResourceAction_PROVISIONER_RESOURCE_ACTION_RECREATE
	default:
		return errors.New("unknown action: " + action.Action)
	}

	switch record.Type {
	case "planned_change":
		progress.Status = mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_PENDING
	case "apply_start":
		progress.Status = mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_RUNNING
	case "apply_complete":
		progress.Status = mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_COMPLETED
	case "apply_errored":
		progress.Status = mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_FAILED
	case "resource_drift":
		progress.Status = mdproto.ProvisionerResourceStatus_PROVISIONER_RESOURCE_STATUS_DRIFT
	default:
		return errors.New("unknown type: " + record.Type)
	}

	progress.ResourceName = action.Resource.ResourceName
	progress.ResourceType = action.Resource.ResourceType
	progress.ResourceId = action.IDValue
	// Convert ResourceKey (which can be a string or an int) into a string if the value is non-null
	if !action.Resource.ResourceKey.IsNull() {
		key, err := ctyconvert.Convert(action.Resource.ResourceKey.Value, cty.String)
		if err != nil {
			return err
		}
		progress.ResourceKey = key.AsString()
	}

	request.Status = mdproto.ProvisionerStatus_PROVISIONER_STATUS_RESOURCE_UPDATE
	request.ResourceProgress = &progress

	return nil
}
