package massdriver

import (
	"context"
	fmt "fmt"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	"github.com/twitchtv/twirp"
)

func UpdateResource(deploymentId, token, resourceId, resourceType, resourceStatus string) error {
	status, err := convertStatus(resourceStatus)
	if err != nil {
		return err
	}

	request := mdproto.UpdateResourceStatusRequest{
		DeploymentId:    deploymentId,
		DeploymentToken: token,
		ResourceId:      resourceId,
		ResourceType:    resourceType,
		ResourceStatus:  status,
	}

	md := mdproto.NewWorkflowServiceProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	_, err = md.UpdateResourceStatus(context.Background(), &request)
	return err
}

func convertStatus(statusString string) (mdproto.ResourceStatus, error) {
	switch statusString {
	case "provisioned":
		return mdproto.ResourceStatus_PROVISIONED, nil
	case "deleted":
		return mdproto.ResourceStatus_DELETED, nil
	default:
		return mdproto.ResourceStatus_DELETED, fmt.Errorf("unknown resource status: %v", statusString)
	}
}
