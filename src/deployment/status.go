package deployment

import (
	"context"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/deployments"
	"github.com/rs/zerolog/log"
)

type DeploymentService interface {
	UpdateDeploymentStatus(ctx context.Context, id string, status deployments.Status) (*deployments.Deployment, error)
}

func UpdateDeploymentStatus(ctx context.Context, svc DeploymentService, id string, status deployments.Status) error {
	_, err := svc.UpdateDeploymentStatus(ctx, id, status)
	if err != nil {
		if deployments.IsInvalidTransitionError(err) {
			// we don't want to fail the whole process if we get an invalid transition error
			log.Warn().Err(err).Msg("Invalid transition error")
			return nil
		}
		return err
	}
	return nil
}
