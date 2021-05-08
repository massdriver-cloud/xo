package massdriver

import (
	"context"
	json "encoding/json"
	ioutil "io/ioutil"
	http "net/http"
	"os"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func UploadArtifactFile(file string, id string, token string) error {
	logger, _ := zap.NewProduction()
	artifactHandle, err := os.Open(file)
	if err != nil {
		logger.Error("Unable to open artifact file"+file, zap.Error(err))
		return err
	}
	defer artifactHandle.Close()

	bytes, _ := ioutil.ReadAll(artifactHandle)
	var artifacts []*Artifact
	err = json.Unmarshal(bytes, &artifacts)
	if err != nil {
		logger.Error("Failed to Unmarshall artifact file", zap.Error(err))
		return err
	}

	err = UploadArtifact(artifacts, id, token)
	return err
}

func UploadArtifact(artifacts []*Artifact, id string, token string) error {
	logger, _ := zap.NewProduction()
	md := NewWorkflowProtobufClient(s.URL, Client)

	header := make(http.Header)
	header.Set("Authorization", "Bearer "+token)

	ctx := context.Background()
	ctx, err := twirp.WithHTTPRequestHeaders(ctx, header)
	if err != nil {
		logger.Error("Error setting twirp headers", zap.Error(err))
		return err
	}

	_, err = md.UploadArtifacts(ctx, &UploadArtifactsRequest{DeploymentId: id, Artifacts: artifacts})
	if err != nil {
		logger.Error("Error sending artifacts to Massdriver", zap.Error(err))
		return err
	}
	return err
}
