package massdriver

import (
	"context"
	ioutil "io/ioutil"
	"os"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
)

func UploadArtifactFile(file string, id string, token string) error {
	artifactHandle, err := os.Open(file)
	if err != nil {
		return err
	}
	defer artifactHandle.Close()

	artifacts, _ := ioutil.ReadAll(artifactHandle)
	err = UploadArtifact(string(artifacts), id, token)
	return err
}

func UploadArtifact(artifacts string, id string, token string) error {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))

	_, err := md.CompleteDeployment(context.Background(), &mdproto.CompleteDeploymentRequest{DeploymentId: id, DeploymentToken: token, Artifacts: artifacts})
	return err
}
