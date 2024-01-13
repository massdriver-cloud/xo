package massdriver

// func (c *MassdriverClient) UploadArtifactFile(file string) error {
// 	artifactHandle, err := os.Open(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer artifactHandle.Close()

// 	var artifact map[string]interface{}
// 	bytes, _ := io.ReadAll(artifactHandle)
// 	json.Unmarshal(bytes, &artifact)

// 	err = c.UploadArtifact(artifact)
// 	return err
// }

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
