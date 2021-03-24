package massdriver

import (
	context "context"
	http "net/http"
)

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func GetDeployment(id string) (string, error) {
	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// http://localhost:4000/rpc/deployment/twirp/mdtwirp.Deployments/Get
	md := NewDeploymentsProtobufClient("http://localhost:4000/rpc/deployment", Client)
	dep, _ := md.Get(context.Background(), &GetDeploymentRequest{Id: id})
	return dep.Token, nil
}
