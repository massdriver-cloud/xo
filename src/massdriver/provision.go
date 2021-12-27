package massdriver

func (c MassdriverClient) ReportProvisionStarted(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_PROVISION_STARTED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportProvisionCompleted(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_PROVISION_COMPLETED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportProvisionFailed(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_PROVISION_FAILED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}
