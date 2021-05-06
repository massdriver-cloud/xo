package massdriver

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func GetDeployment(id string, token string) (*Deployment, error) {
	logger, _ := zap.NewProduction()
	// The URL below sucks. We need to fix something in elixir so its not so redundant...
	// curl http://localhost:4000/rpc/workflow/twirp/mdtwirp.Workflow/GetDeployment
	md := NewWorkflowProtobufClient(s.URL, Client)
	dep, err := md.GetDeployment(context.Background(), &GetDeploymentRequest{Id: id, Token: token})
	if err != nil {
		logger.Error("Error fetching Deployment object from Massdriver", zap.Error(err))
	}
	return dep, err
}

func WriteDeploymentToFile(dep *Deployment, dest string) error {
	logger, _ := zap.NewProduction()
	inputHandle, err := os.OpenFile(dest+"/inputs.tfvars.json", os.O_CREATE, 0644)
	if err != nil {
		logger.Error("Error opening inputs.tfvars.json", zap.Error(err))
		return err
	}
	connHandle, err := os.OpenFile(dest+"/connections.tfvars.json", os.O_CREATE, 0644)
	if err != nil {
		logger.Error("Error opening connections.tfvars.json", zap.Error(err))
		return err
	}

	writeSchema(dep.Inputs, inputHandle)
	writeSchema(dep.Connections, connHandle)

	return err
}

func writeSchema(schema *structpb.Struct, file io.Writer) error {
	logger, _ := zap.NewProduction()
	jsonString, err := schema.MarshalJSON()
	if err != nil {
		logger.Error("Error marshaling JSON file", zap.Error(err))
		return err
	}
	_, err = file.Write(jsonString)
	return err
}
