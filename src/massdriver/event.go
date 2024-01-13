package massdriver

import (
	"os"
	"time"
)

const EVENT_TYPE_PROVISION_STARTED string = "provision_started"
const EVENT_TYPE_PROVISION_COMPLETED string = "provision_completed"
const EVENT_TYPE_PROVISION_FAILED string = "provision_failed"
const EVENT_TYPE_DECOMMISSION_STARTED string = "decommission_started"
const EVENT_TYPE_DECOMMISSION_COMPLETED string = "decommission_completed"
const EVENT_TYPE_DECOMMISSION_FAILED string = "decommission_failed"
const EVENT_TYPE_ARTIFACT_CREATED string = "artifact_created"
const EVENT_TYPE_ARTIFACT_UPDATED string = "artifact_updated"
const EVENT_TYPE_ARTIFACT_DELETED string = "artifact_deleted"

const EVENT_TYPE_PROVISIONER_VIOLATION string = "provisioner_violation"

const EVENT_TYPE_PROVISIONER_COMPLETE string = "provisioner_completed"
const EVENT_TYPE_PROVISIONER_ERROR string = "provisioner_error"
const EVENT_TYPE_RESOURCE_CREATE_PENDING string = "resource_create_pending"
const EVENT_TYPE_RESOURCE_CREATE_RUNNING string = "resource_create_running"
const EVENT_TYPE_RESOURCE_CREATE_COMPLETED string = "resource_create_completed"
const EVENT_TYPE_RESOURCE_CREATE_FAILED string = "resource_create_failed"
const EVENT_TYPE_RESOURCE_UPDATE_PENDING string = "resource_update_pending"
const EVENT_TYPE_RESOURCE_UPDATE_RUNNING string = "resource_update_running"
const EVENT_TYPE_RESOURCE_UPDATE_COMPLETED string = "resource_update_completed"
const EVENT_TYPE_RESOURCE_UPDATE_FAILED string = "resource_update_failed"
const EVENT_TYPE_RESOURCE_DELETE_PENDING string = "resource_delete_pending"
const EVENT_TYPE_RESOURCE_DELETE_RUNNING string = "resource_delete_running"
const EVENT_TYPE_RESOURCE_DELETE_COMPLETED string = "resource_delete_completed"
const EVENT_TYPE_RESOURCE_DELETE_FAILED string = "resource_delete_failed"
const EVENT_TYPE_RESOURCE_DRIFT_DETECTED string = "resource_drift_detected"

type Event struct {
	Metadata EventMetadata `json:"metadata"`
	Payload  interface{}   `json:"payload,omitempty"`
}

type EventMetadata struct {
	Timestamp   string `json:"timestamp"`
	Provisioner string `json:"provisioner"`
	Version     string `json:"version,omitempty"`
	EventType   string `json:"event_type"`
}

type EventPayloadProvisionerStatus struct {
	DeploymentId string `json:"deployment_id"`
}

type EventPayloadResourceProgress struct {
	DeploymentId string `json:"deployment_id"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	ResourceKey  string `json:"resource_key,omitempty"`
	ResourceId   string `json:"resource_id,omitempty"`
}

type EventPayloadArtifact struct {
	DeploymentId string                 `json:"deployment_id"`
	Artifact     map[string]interface{} `json:"artifact"`
}

type EventPayloadDiagnostic struct {
	DeploymentId string `json:"deployment_id"`
	Message      string `json:"error_message"`
	Details      string `json:"error_details"`
	Level        string `json:"error_level"`
}

type EventPayloadOPAViolation struct {
	DeploymentId string      `json:"deployment_id"`
	Rule         string      `json:"opa_rule"`
	Value        interface{} `json:"opa_value"`
}

var EventTimeString = time.Now().String

func NewEvent(eventType string) *Event {
	event := new(Event)
	event.Metadata.EventType = eventType
	event.Metadata.Timestamp = EventTimeString()
	event.Metadata.Provisioner = os.Getenv("MASSDRIVER_PROVISIONER")
	return event
}
