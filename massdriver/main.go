package massdriver

import (
	context "context"
	"io"
	http "net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

func GetDeployment(id string, token string) (*Deployment, error) {
	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// curl http://localhost:4000/rpc/workflow/twirp/mdtwirp.Workflow/GetDeployment
	md := NewWorkflowProtobufClient(s.URL, Client)
	dep, err := md.GetDeployment(context.Background(), &GetDeploymentRequest{Id: id, Token: token})
	return dep, err
}

func WriteDeploymentToFile(dep *Deployment, dest string) error {
	inputHandle, err := os.OpenFile(dest+"/inputs.tfvars.json", os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	connHandle, err := os.OpenFile(dest+"/connections.tfvars.json", os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	writeSchema(dep.Inputs, inputHandle)
	writeSchema(dep.Connections, connHandle)

	return err
}

func writeSchema(schema *structpb.Struct, file io.Writer) error {
	jsonString, err := schema.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	return err
}
