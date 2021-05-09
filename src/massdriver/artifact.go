package massdriver

import (
	"context"
	json "encoding/json"
	ioutil "io/ioutil"
	http "net/http"
	"os"

	"github.com/twitchtv/twirp"
)

func UploadArtifactFile(file string, id string, token string) error {
	artifactHandle, err := os.Open(file)
	if err != nil {
		return err
	}
	defer artifactHandle.Close()

	bytes, _ := ioutil.ReadAll(artifactHandle)
	var artifacts []*Artifact
	err = json.Unmarshal(bytes, &artifacts)
	if err != nil {
		return err
	}

	err = UploadArtifact(artifacts, id, token)
	return err
}

func UploadArtifact(artifacts []*Artifact, id string, token string) error {
	md := NewWorkflowProtobufClient(s.URL, Client)

	header := make(http.Header)
	header.Set("Authorization", "Bearer "+token)

	ctx := context.Background()
	ctx, err := twirp.WithHTTPRequestHeaders(ctx, header)
	if err != nil {
		return err
	}

	_, err = md.UploadArtifacts(ctx, &UploadArtifactsRequest{DeploymentId: id, Artifacts: artifacts})
	if err != nil {
		return err
	}
	return err
}
