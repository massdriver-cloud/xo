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

func StartDeployment(id string, token string, dest string) (*mdproto.Deployment, error) {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	resp, err := md.StartDeployment(context.Background(), &mdproto.StartDeploymentRequest{DeploymentId: id, DeploymentToken: token})
	if err != nil {
		return nil, err
	}

	// Write out params
	paramsHandle, err := OutputGenerator(path.Join(dest, ParamsFileName))
	if err != nil {
		return nil, err
	}
	err = writeSchema(resp.Deployment.Params, paramsHandle)
	if err != nil {
		return nil, err
	}

	// Write out connections
	connectionsHandle, err := OutputGenerator(path.Join(dest, ConnectionsFileName))
	if err != nil {
		return nil, err
	}
	err = writeSchema(resp.Deployment.Connections, connectionsHandle)
	if err != nil {
		return nil, err
	}

	// Write out Organization ID
	organizationHandle, err := OutputGenerator(path.Join(dest, OrganizationFileName))
	if err != nil {
		return nil, err
	}
	_, err = organizationHandle.Write([]byte(resp.Deployment.Organization.Id))
	if err != nil {
		return nil, err
	}

	// Write out Project ID
	projectHandle, err := OutputGenerator(path.Join(dest, ProjectFileName))
	if err != nil {
		return nil, err
	}
	_, err = projectHandle.Write([]byte(resp.Deployment.Project.Id))
	if err != nil {
		return nil, err
	}

	// Write out Target ID
	targetHandle, err := OutputGenerator(path.Join(dest, TargetFileName))
	if err != nil {
		return nil, err
	}
	_, err = targetHandle.Write([]byte(resp.Deployment.Target.Id))
	if err != nil {
		return nil, err
	}

	// Write out Bundle type
	bundleHandle, err := OutputGenerator(path.Join(dest, BundleFileName))
	if err != nil {
		return nil, err
	}
	_, err = bundleHandle.Write([]byte(resp.Deployment.Bundle.Type))
	return resp.Deployment, err
}

func writeSchema(schema *structpb.Struct, file io.Writer) error {
	jsonString, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	return err
}
