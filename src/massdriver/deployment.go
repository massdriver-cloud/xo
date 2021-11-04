package massdriver

import (
	"context"
	"path"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"

	"github.com/twitchtv/twirp"
)

var ParamsFileName = "params.auto.tfvars.json"
var ConnectionsFileName = "connections.auto.tfvars.json"
var OrganizationFileName = "organization.txt"
var ProjectFileName = "project.txt"
var TargetFileName = "target.txt"
var BundleFileName = "bundle.txt"

func StartDeployment(id string, token string, dest string) error {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	resp, err := md.StartDeployment(context.Background(), &mdproto.StartDeploymentRequest{DeploymentId: id, DeploymentToken: token})
	if err != nil {
		return err
	}

	// Write out params
	paramsHandle, err := OutputGenerator(path.Join(dest, ParamsFileName))
	if err != nil {
		return err
	}
	_, err = paramsHandle.Write([]byte(resp.Deployment.Params))
	if err != nil {
		return err
	}

	// Write out connections
	connectionsHandle, err := OutputGenerator(path.Join(dest, ConnectionsFileName))
	if err != nil {
		return err
	}
	_, err = connectionsHandle.Write([]byte(resp.Deployment.ConnectionParams))
	if err != nil {
		return err
	}
	return nil
}
