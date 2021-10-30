package massdriver

import (
	"context"
	json "encoding/json"
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

	bytes, _ := ioutil.ReadAll(artifactHandle)
	var artifacts []*mdproto.Artifact
	err = json.Unmarshal(bytes, &artifacts)
	if err != nil {
		return err
	}

	err = UploadArtifact(artifacts, id, token)
	return err
}

func UploadArtifact(artifacts []*mdproto.Artifact, id string, token string) error {
	md := mdproto.NewWorkflowServiceJSONClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))

	_, err := md.CompleteDeployment(context.Background(), &mdproto.CompleteDeploymentRequest{DeploymentId: id, DeploymentToken: token, Artifacts: artifacts})
	return err
}
