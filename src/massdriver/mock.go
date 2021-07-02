package massdriver

import (
	"context"
	"net/http"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

type massdriverMockServer struct{}

func (s *massdriverMockServer) GetDeployment(context.Context, *GetDeploymentRequest) (*Deployment, error) {
	mockParams, _ := structpb.NewStruct(map[string]interface{}{
		"some_key": "some_value",
	})
	mockConnections, _ := structpb.NewStruct(map[string]interface{}{
		"default": map[string]interface{}{
			"aws_access_key_id":     "ACOVIBUOISKLWJEFKJL",
			"aws_secret_access_key": "8ba0u90uwe9fuq90j3490tj0q923u12093u09gj90u130",
		},
	})

	return &Deployment{
		Id:          "FAKEID",
		Status:      DeploymentStatus_PENDING,
		Params:      mockParams,
		Connections: mockConnections,
	}, nil
}

func (s *massdriverMockServer) UploadArtifacts(context.Context, *UploadArtifactsRequest) (*Deployment, error) {
	return &Deployment{
		Id:     "FAKEID",
		Status: DeploymentStatus_COMPLETED,
	}, nil
}

// Run the implementation in a local server
func RunMockServer(port string) error {
	mdMock := NewWorkflowServer(&massdriverMockServer{})
	mux := http.NewServeMux()
	mux.Handle(mdMock.PathPrefix(), mdMock)
	return http.ListenAndServe(":"+port, mux)
}
