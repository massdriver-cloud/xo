package cmd

import (
	"fmt"
	"os"
	"xo/src/deployment"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/deployments"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
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

var deploymentStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the plan has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

var deploymentCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the plan has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

var deploymentFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the plan has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentStatus,
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.AddCommand(deploymentStartCmd)
	deploymentCmd.AddCommand(deploymentCompleteCmd)
	deploymentCmd.AddCommand(deploymentFailCmd)
	deploymentCmd.PersistentFlags().StringP("id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver deployment ID")
}

func RunDeploymentStatus(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "RunDeploymentStatus")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	deploymentID, getErr := cmd.Flags().GetString("id")
	if getErr != nil {
		return telemetry.LogError(span, getErr, "unable to read id flag")
	}
	if deploymentID == "" {
		missingErr := fmt.Errorf("deployment ID not set and MASSDRIVER_DEPLOYMENT_ID environment variable is not set")
		return telemetry.LogError(span, missingErr, "an error occurred while updating deployment status")
	}

	var status deployments.Status
	switch cmd.Name() {
	case "start":
		status = deployments.StatusRunning
	case "complete":
		status = deployments.StatusCompleted
	case "fail":
		status = deployments.StatusFailed
	default:
		return fmt.Errorf("unknown deployment status: %s", cmd.Name())
	}

	log.Info().Msgf("sending deployment status: %s", status)
	mdClient, err := client.New()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	service := deployments.NewService(mdClient)
	updateErr := deployment.UpdateDeploymentStatus(ctx, service, deploymentID, status)
	if updateErr != nil {
		return telemetry.LogError(span, updateErr, "an error occurred while reporting deployment status")
	}

	log.Info().Msgf("deployment status event sent: %s", status)
	return nil
}
