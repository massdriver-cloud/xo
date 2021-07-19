package massdriver

import (
	"context"
	json "encoding/json"
	"io"
	"os"

	"github.com/twitchtv/twirp"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func StartDeployment(id string, token string) (*Deployment, error) {
	md := NewWorkflowProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	return md.StartDeployment(context.Background(), &StartDeploymentRequest{Id: id, Token: token})
}

func WriteDeploymentToFile(dep *Deployment, dest string) error {
	inputHandle, err := os.OpenFile(dest+"/params.auto.tfvars.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	connHandle, err := os.OpenFile(dest+"/connections.auto.tfvars.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = writeSchema(dep.Params, inputHandle)
	if err != nil {
		return err
	}
	err = writeSchema(dep.Connections, connHandle)
	return err
}

func writeSchema(schema *structpb.Struct, file io.Writer) error {
	jsonString, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	return err
}
