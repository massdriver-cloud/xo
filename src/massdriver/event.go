package massdriver

import (
	"os"
	"time"
)

const EVENT_TYPE_PLAN_STARTED string = "plan_started"
const EVENT_TYPE_PLAN_COMPLETED string = "plan_completed"
const EVENT_TYPE_PLAN_FAILED string = "plan_failed"
const EVENT_TYPE_PROVISION_STARTED string = "provision_started"
const EVENT_TYPE_PROVISION_COMPLETED string = "provision_completed"
const EVENT_TYPE_PROVISION_FAILED string = "provision_failed"
const EVENT_TYPE_DECOMMISSION_STARTED string = "decommission_started"
const EVENT_TYPE_DECOMMISSION_COMPLETED string = "decommission_completed"
const EVENT_TYPE_DECOMMISSION_FAILED string = "decommission_failed"

const EVENT_TYPE_ARTIFACT_CREATED string = "artifact_created"
const EVENT_TYPE_ARTIFACT_UPDATED string = "artifact_updated"
const EVENT_TYPE_ARTIFACT_DELETED string = "artifact_deleted"

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

type EventPayloadArtifact struct {
	DeploymentId string                 `json:"deployment_id"`
	Artifact     map[string]interface{} `json:"artifact"`
}

var EventTimeString = time.Now().String

func NewEvent(eventType string) *Event {
	event := new(Event)
	event.Metadata.EventType = eventType
	event.Metadata.Timestamp = EventTimeString()
	event.Metadata.Provisioner = os.Getenv("MASSDRIVER_PROVISIONER")
	return event
}
