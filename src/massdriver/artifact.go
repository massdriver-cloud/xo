package massdriver

import (
	"encoding/json"
	ioutil "io/ioutil"
	"os"
)

func (c *MassdriverClient) UploadArtifactFile(file string, id string) error {
	artifactHandle, err := os.Open(file)
	if err != nil {
		return err
	}
	defer artifactHandle.Close()

	var artifacts []map[string]interface{}
	bytes, _ := ioutil.ReadAll(artifactHandle)
	json.Unmarshal(bytes, &artifacts)

	err = c.UploadArtifact(artifacts, id)
	return err
}

func (c *MassdriverClient) UploadArtifact(artifacts []map[string]interface{}, id string) error {
	event := NewEvent(EVENT_TYPE_ARTIFACT_UPDATE)
	event.Payload = EventPayloadArtifacts{DeploymentId: id, Artifacts: artifacts}
	return c.PublishEventToSNS(event)
}
