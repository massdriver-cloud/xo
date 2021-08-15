package massdriver

import (
	"context"
	fmt "fmt"
	http "net/http"

	"github.com/twitchtv/twirp"
)

func UpdateResource(deploymentId, token, resourceId, resourceType, resourceStatus string) error {
	status, err := convertStatus(resourceStatus)
	if err != nil {
		return err
	}

	request := UpdateResourceStatusRequest{
		DeploymentId:   deploymentId,
		ResourceId:     resourceId,
		ResourceType:   resourceType,
		ResourceStatus: status,
	}

	header := make(http.Header)
	header.Set("Authorization", "Bearer "+token)

	ctx := context.Background()
	ctx, err = twirp.WithHTTPRequestHeaders(ctx, header)
	if err != nil {
		return err
	}

	md := NewWorkflowProtobufClient(s.URL, Client, twirp.WithClientPathPrefix("/rpc/twirp"))
	_, err = md.UpdateResourceStatus(ctx, &request)
	return err
}

func convertStatus(statusString string) (ResourceStatus, error) {
	switch statusString {
	case "provisioned":
		return ResourceStatus_PROVISIONED, nil
	case "deleted":
		return ResourceStatus_DELETED, nil
	default:
		return ResourceStatus_DELETED, fmt.Errorf("unknown resource status: %v", statusString)
	}
}
