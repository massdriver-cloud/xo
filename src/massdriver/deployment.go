package massdriver

import (
	"context"
	json "encoding/json"
	"io"
	"path"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

var ParamsFileName = "params.auto.tfvars.json"
var ConnectionsFileName = "connections.auto.tfvars.json"
var OrganizationFileName = "organization.txt"
var ProjectFileName = "project.txt"
var TargetFileName = "target.txt"
var BundleFileName = "bundle.txt"

func StartDeployment(id string, token string, dest string) error {
	md := mdproto.NewWorkflowServiceJSONClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	resp, err := md.StartDeployment(context.Background(), &mdproto.StartDeploymentRequest{DeploymentId: id, DeploymentToken: token})
	if err != nil {
		return err
	}

	// Write out params
	paramsHandle, err := OutputGenerator(path.Join(dest, ParamsFileName))
	if err != nil {
		return err
	}
	err = writeSchema(resp.Deployment.Params, paramsHandle)
	if err != nil {
		return err
	}

	// Write out connections
	connectionsHandle, err := OutputGenerator(path.Join(dest, ConnectionsFileName))
	if err != nil {
		return err
	}
	err = writeSchema(resp.Deployment.ConnectionParams, connectionsHandle)
	if err != nil {
		return err
	}
	return nil
}

func writeSchema(schema *structpb.Struct, file io.Writer) error {
	jsonString, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	return err
}
