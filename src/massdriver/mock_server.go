package massdriver

import (
	"context"
	"net/http"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
)

type massdriverMockServer struct {
	params      string
	connections string
}

func (s *massdriverMockServer) StartDeployment(context.Context, *mdproto.StartDeploymentRequest) (*mdproto.StartDeploymentResponse, error) {
	return &mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Id:     "FAKEID",
			Status: mdproto.DeploymentStatus_PENDING,
			Organization: &mdproto.Organization{
				Id: "organization",
			},
			Params:           s.params,
			ConnectionParams: s.connections,
		},
	}, nil
}

func (s *massdriverMockServer) CompleteDeployment(context.Context, *mdproto.CompleteDeploymentRequest) (*mdproto.CompleteDeploymentResponse, error) {
	return &mdproto.CompleteDeploymentResponse{}, nil
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
