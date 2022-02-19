package massdriver

import (
	"context"

	"go.opentelemetry.io/otel"
)

func (c MassdriverClient) ReportProvisionStarted(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportProvisionStarted")
	defer span.End()
	event := NewEvent(EVENT_TYPE_PROVISION_STARTED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportProvisionCompleted(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportProvisionCompleted")
	defer span.End()
	event := NewEvent(EVENT_TYPE_PROVISION_COMPLETED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportProvisionFailed(ctx context.Context, deploymentId string) error {
	_, span := otel.Tracer("xo").Start(ctx, "ReportProvisionFailed")
	defer span.End()
	event := NewEvent(EVENT_TYPE_PROVISION_FAILED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}
