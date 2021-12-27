package massdriver

func (c MassdriverClient) ReportDecommissionStarted(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_DECOMMISSION_STARTED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportDecommissionCompleted(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_DECOMMISSION_COMPLETED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}

func (c MassdriverClient) ReportDecommissionFailed(deploymentId string) error {
	event := NewEvent(EVENT_TYPE_DECOMMISSION_FAILED)
	event.Payload = EventPayloadProvisionerStatus{DeploymentId: deploymentId}
	return c.PublishEventToSNS(event)
}
