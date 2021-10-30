package massdriver

import (
	"context"
	"net/http"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type massdriverMockServer struct {
	params      *map[string]interface{}
	connections *map[string]interface{}
}

func (s *massdriverMockServer) StartDeployment(context.Context, *mdproto.StartDeploymentRequest) (*mdproto.StartDeploymentResponse, error) {
	mockParams, _ := structpb.NewStruct(*s.params)
	mockConnectionParams, _ := structpb.NewStruct(*s.connections)

	return &mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Id:     "FAKEID",
			Status: mdproto.DeploymentStatus_PENDING,
			Organization: &mdproto.Organization{
				Id: "organization",
			},
			Params:           mockParams,
			ConnectionParams: mockConnectionParams,
		},
	}, nil
}

func (s *massdriverMockServer) CompleteDeployment(context.Context, *mdproto.CompleteDeploymentRequest) (*mdproto.CompleteDeploymentResponse, error) {
	return &mdproto.CompleteDeploymentResponse{}, nil
}

func RunMockServer(port string, params *map[string]interface{}, connections *map[string]interface{}) error {
	mockServer := massdriverMockServer{}
	mockServer.params = params
	mockServer.connections = connections

	mdMock := mdproto.NewWorkflowServiceServer(&mockServer, twirp.WithServerPathPrefix("/rpc/twirp"))
	mux := http.NewServeMux()
	mux.Handle(mdMock.PathPrefix(), mdMock)
	return http.ListenAndServe(":"+port, mux)
}
