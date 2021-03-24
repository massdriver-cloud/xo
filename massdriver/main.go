package massdriver

import (
	context "context"
	http "net/http"

	"github.com/kelseyhightower/envconfig"
)

var (
	Client HTTPClient
	s      Specification
)

func init() {
	Client = &http.Client{}
	envconfig.Process("massdriver", &s)
}

type Specification struct {
	URL string `default:"http://localhost:4000/rpc/deployment"`
}

func GetDeployment(id string) (string, error) {
	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// http://localhost:4000/rpc/deployment/twirp/mdtwirp.Deployments/Get
	md := NewDeploymentsProtobufClient(s.URL, Client)
	dep, _ := md.Get(context.Background(), &GetDeploymentRequest{Id: id})
	return dep.Token, nil
}
