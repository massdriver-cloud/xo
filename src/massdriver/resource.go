package massdriver

import (
	"context"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"

	"github.com/twitchtv/twirp"
)

func SendProvisionerProgressUpdate(message *mdproto.ProvisionerProgressUpdateRequest) error {
	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	_, err := md.ProvisionerProgressUpdate(context.Background(), message)
	return err
}
