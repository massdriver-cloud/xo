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
	URL string `default:"http://localhost:4000/rpc/workflow"`
}

func GetDeployment(id string, token string) (string, error) {
	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// curl http://localhost:4000/rpc/workflow/twirp/mdtwirp.Workflow/GetDeployment
	md := NewWorkflowProtobufClient(s.URL, Client)
        dep, _ := md.GetDeployment(context.Background(), &GetDeploymentRequest{Id: id, Token: token})
        _ = dep
	return "nothing", nil
}
