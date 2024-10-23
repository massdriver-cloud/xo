package cmd

import (
	"errors"
	"os"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var descritionLong = `
	Publishes an event to AWS SNS, which distributes the event to SQS subscribers.

	This command is designed to be executed in automation, and therefore takes inputs
	from environment variables. Specifically, the following environment variables
	are read and used to populate event data:

	MASSDRIVER_DEPLOYMENT_ID
	MASSDRIVER_EVENT_TOPIC_ARN
	MASSDRIVER_PROVISIONER

	Be sure these environment variables are set, and you have access to the SNS topic.
	`

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Massdriver deployment events",
	Long:  ``,
}

var deploymentPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Manage Massdriver decommission events",
	Long:  ``,
}

var deploymentProvisionCmd = &cobra.Command{
	Use:     "provision",
	Aliases: []string{"apply"},
	Short:   "Manage Massdriver provision events",
	Long:    ``,
}

var deploymentDecommissionCmd = &cobra.Command{
	Use:     "decommission",
	Aliases: []string{"destroy"},
	Short:   "Manage Massdriver decommission events",
	Long:    ``,
}

var deploymentPlanStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the plan has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentPlanCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the plan has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentPlanFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the plan has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentProvisionStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the provision has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentProvisionCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the provision has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentProvisionFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the provision has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the decommission has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the decommission has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the decommission has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.AddCommand(deploymentProvisionCmd)
	deploymentCmd.AddCommand(deploymentDecommissionCmd)

	deploymentPlanCmd.AddCommand(deploymentPlanStartCmd)
	deploymentPlanCmd.AddCommand(deploymentPlanCompleteCmd)
	deploymentPlanCmd.AddCommand(deploymentPlanFailCmd)

	deploymentProvisionCmd.AddCommand(deploymentProvisionStartCmd)
	deploymentProvisionCmd.AddCommand(deploymentProvisionCompleteCmd)
	deploymentProvisionCmd.AddCommand(deploymentProvisionFailCmd)

	deploymentDecommissionCmd.AddCommand(deploymentDecommissionStartCmd)
	deploymentDecommissionCmd.AddCommand(deploymentDecommissionCompleteCmd)
	deploymentDecommissionCmd.AddCommand(deploymentDecommissionFailCmd)
}

func RunDeploymentStatus(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "RunDeploymentStatus")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	deploymentId := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentId == "" {
		err := errors.New("MASSDRIVER_DEPLOYMENT_ID environment variable must be set")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	deploymentStatus := cmd.Parent().Use + "_" + cmd.Use

	log.Info().Msgf("sending deployment status event: %s", deploymentStatus)
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msgf("an error occurred while sending deployment status event: %s", deploymentStatus)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = mdClient.ReportDeploymentStatus(ctx, deploymentId, deploymentStatus)
	if err != nil {
		log.Error().Err(err).Msgf("an error occurred while sending deployment status event: %s", deploymentStatus)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Info().Msgf("deployment status event sent: %s", deploymentStatus)
	return nil
}
