package massdriver

import (
	"context"
	"path"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"

	"github.com/twitchtv/twirp"
)

var ParamsFileName = "params.auto.tfvars.json"
var ConnectionsFileName = "connections.auto.tfvars.json"

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

func FailDeployment(id string, token string) error {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	_, err := md.FailDeployment(context.Background(), &mdproto.FailDeploymentRequest{DeploymentId: id, DeploymentToken: token})
	return err
}

func DestroyedDeployment(id string, token string) error {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	_, err := md.CompleteDestruction(context.Background(), &mdproto.CompleteDestructionRequest{DeploymentId: id, DeploymentToken: token})
	return err
}
