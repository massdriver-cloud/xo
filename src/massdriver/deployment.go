package massdriver

import (
	"context"
	json "encoding/json"
	"io"
	"os"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

var ParamsFileName = "params.auto.tfvars.json"
var ConnectionsFileName = "connections.auto.tfvars.json"

func StartDeployment(id string, token string) (*mdproto.Deployment, error) {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	resp, err := md.StartDeployment(context.Background(), &mdproto.StartDeploymentRequest{DeploymentId: id, DeploymentToken: token})
	if err != nil {
		return nil, err
	}
	return resp.Deployment, err
}

func WriteDeploymentToFile(dep *mdproto.Deployment, dest string) error {
	inputHandle, err := os.OpenFile(dest+"/"+ParamsFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	connHandle, err := os.OpenFile(dest+"/"+ConnectionsFileName, os.O_CREATE|os.O_WRONLY, 0644)
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
