package massdriver

import (
	"context"
	"errors"
	"xo/src/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func (c *MassdriverClient) ReportDeploymentStatus(ctx context.Context, deploymentId string, status string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportDeploymentStatus")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	var event *Event
	switch status {
	case "provision_start":
		event = NewEvent(EVENT_TYPE_PROVISION_STARTED)
	case "provision_complete":
		event = NewEvent(EVENT_TYPE_PROVISION_COMPLETED)
	case "provision_fail":
		event = NewEvent(EVENT_TYPE_PROVISION_FAILED)
	case "decommission_start":
		event = NewEvent(EVENT_TYPE_DECOMMISSION_STARTED)
	case "decommission_complete":
		event = NewEvent(EVENT_TYPE_DECOMMISSION_COMPLETED)
	case "decommission_fail":
		event = NewEvent(EVENT_TYPE_DECOMMISSION_FAILED)
	default:
		err := errors.New("Unknown deployment status: " + status)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}

	return c.PublishEvent(event)
}
