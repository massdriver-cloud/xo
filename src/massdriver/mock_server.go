package massdriver

import (
	"context"
	"net/http"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type massdriverMockServer struct{}

func (s *massdriverMockServer) StartDeployment(context.Context, *mdproto.StartDeploymentRequest) (*mdproto.StartDeploymentResponse, error) {
	mockParams, _ := structpb.NewStruct(map[string]interface{}{
		"name": "value",
	})
	mockConnections, _ := structpb.NewStruct(map[string]interface{}{
		"default": map[string]interface{}{
			"aws_access_key_id":     "ACOVIBUOISKLWJEFKJL",
			"aws_secret_access_key": "8ba0u90uwe9fuq90j3490tj0q923u12093u09gj90u130",
		},
	})

	return &mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Id:     "FAKEID",
			Status: mdproto.DeploymentStatus_PENDING,
			Bundle: &mdproto.Bundle{
				Type: "test-bundle",
			},
			Params:      mockParams,
			Connections: mockConnections,
		},
	}, nil
}

func (s *massdriverMockServer) UploadArtifacts(context.Context, *mdproto.UploadArtifactsRequest) (*mdproto.UploadArtifactsResponse, error) {
	return &mdproto.UploadArtifactsResponse{}, nil
}

func (s *massdriverMockServer) UpdateResourceStatus(context.Context, *mdproto.UpdateResourceStatusRequest) (*mdproto.UpdateResourceStatusResponse, error) {
	return &mdproto.UpdateResourceStatusResponse{}, nil
}

func RunMockServer(port string) error {
	mdMock := mdproto.NewWorkflowServiceServer(&massdriverMockServer{}, twirp.WithServerPathPrefix("/rpc/twirp"))
	mux := http.NewServeMux()
	mux.Handle(mdMock.PathPrefix(), mdMock)
	return http.ListenAndServe(":"+port, mux)
}
