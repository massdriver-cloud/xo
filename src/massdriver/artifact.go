package massdriver

func PublishArtifact(c *MassdriverClient, artifact map[string]interface{}) error {
	event := NewEvent(EVENT_TYPE_ARTIFACT_UPDATED)
	event.Payload = EventPayloadArtifact{DeploymentId: c.Specification.DeploymentID, Artifact: artifact}
	return c.PublishEvent(event)
}

func DeleteArtifact(c *MassdriverClient, artifact map[string]interface{}) error {
	event := NewEvent(EVENT_TYPE_ARTIFACT_DELETED)
	event.Payload = EventPayloadArtifact{DeploymentId: c.Specification.DeploymentID, Artifact: artifact}
	return c.PublishEvent(event)
}
