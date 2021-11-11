package massdriver

import (
	"context"
	"fmt"
	"net/http"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
)

type massdriverMockServer struct {
	params      string
	connections string
}

func (s *massdriverMockServer) StartDeployment(ctx context.Context, req *mdproto.StartDeploymentRequest) (*mdproto.StartDeploymentResponse, error) {
	fmt.Printf("Received StartDeploymentRequest: %v\n", req)
	return &mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Id:     "FAKEID",
			Status: mdproto.DeploymentStatus_DEPLOYMENT_STATUS_PENDING,
			Organization: &mdproto.Organization{
				Id: "organization",
			},
			Params:           s.params,
			ConnectionParams: s.connections,
		},
	}, nil
}

func (s *massdriverMockServer) CompleteDeployment(ctx context.Context, req *mdproto.CompleteDeploymentRequest) (*mdproto.CompleteDeploymentResponse, error) {
	fmt.Printf("Received CompleteDeploymentRequest: %v\n", req)
	return &mdproto.CompleteDeploymentResponse{}, nil
}

func (s *massdriverMockServer) FailDeployment(ctx context.Context, req *mdproto.FailDeploymentRequest) (*mdproto.FailDeploymentResponse, error) {
	fmt.Printf("Received CompleteDeploymentRequest: %v\n", req)
	return &mdproto.FailDeploymentResponse{}, nil
}

func (s *massdriverMockServer) ProvisionerProgressUpdate(ctx context.Context, req *mdproto.ProvisionerProgressUpdateRequest) (*mdproto.ProvisionerProgressUpdateResponse, error) {
	fmt.Printf("Received ProvisionerProgressUpdateRequest: %v\n", req)
	return &mdproto.ProvisionerProgressUpdateResponse{}, nil
}

func RunMockServer(port string, params string, connections string) error {
	mockServer := massdriverMockServer{}
	mockServer.params = params
	mockServer.connections = connections

	mdMock := mdproto.NewWorkflowServiceServer(&mockServer, twirp.WithServerPathPrefix("/rpc/twirp"))
	mux := http.NewServeMux()
	mux.Handle(mdMock.PathPrefix(), mdMock)
	return http.ListenAndServe(":"+port, mux)
}
