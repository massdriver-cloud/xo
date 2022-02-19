package massdriver

import (
	"context"

	"go.opentelemetry.io/otel"
)

func (c MassdriverClient) ReportDecommissionStarted(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportDecommissionStarted")
	defer span.End()
	event := NewEvent(EVENT_TYPE_DECOMMISSION_STARTED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportDecommissionCompleted(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportDecommissionCompleted")
	defer span.End()
	event := NewEvent(EVENT_TYPE_DECOMMISSION_COMPLETED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportDecommissionFailed(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportDecommissionFailed")
	defer span.End()
	event := NewEvent(EVENT_TYPE_DECOMMISSION_FAILED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}
