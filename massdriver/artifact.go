package massdriver

import (
	"context"
	json "encoding/json"
	ioutil "io/ioutil"
	"os"
)

func UploadArtifactFile(file string, id string, token string) error {
	artifactHandle, err := os.Open(file)
	if err != nil {
		return err
	}
	defer artifactHandle.Close()

	bytes, _ := ioutil.ReadAll(artifactHandle)
	artifacts, err := createArtifactsFromJsonBytes(bytes)
	if err != nil {
		return err
	}

	err = UploadArtifact(artifacts, id, token)
	return err
}

func UploadArtifact(artifacts []*Artifact, id string, token string) error {
	md := NewWorkflowProtobufClient(s.URL, Client)
	_, err := md.UploadArtifacts(context.Background(), &UploadArtifactsRequest{DeploymentId: id, Token: token, Artifacts: artifacts})
	return err
}

func createArtifactsFromJsonBytes(bytes []byte) ([]*Artifact, error) {
	var artifacts []*Artifact
	err := json.Unmarshal(bytes, &artifacts)
	if err != nil {
		return artifacts, err
	}

	return artifacts, err
}
