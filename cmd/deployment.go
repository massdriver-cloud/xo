package cmd

import (
	"errors"
	"fmt"
	"os"
	"xo/src/massdriver"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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

var deploymentProvisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Manage Massdriver provision events",
	Long:  ``,
}

var deploymentDecommissionCmd = &cobra.Command{
	Use:   "decommission",
	Short: "Manage Massdriver decommission events",
	Long:  ``,
}

var deploymentProvisionStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the provision has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentProvisionStart,
	DisableFlagsInUseLine: true,
}

var deploymentProvisionCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the provision has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentProvisionComplete,
	DisableFlagsInUseLine: true,
}

var deploymentProvisionFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the provision has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentProvisionFail,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionStartCmd = &cobra.Command{
	Use:                   "start",
	Short:                 "Generate event notifying Massdriver the decommission has started",
	Long:                  descritionLong,
	RunE:                  RunDeploymentDecommissionStart,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionCompleteCmd = &cobra.Command{
	Use:                   "complete",
	Short:                 "Generate event notifying Massdriver the decommission has completed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentDecommissionComplete,
	DisableFlagsInUseLine: true,
}

var deploymentDecommissionFailCmd = &cobra.Command{
	Use:                   "fail",
	Short:                 "Generate event notifying Massdriver the decommission has failed",
	Long:                  descritionLong,
	RunE:                  RunDeploymentDecommissionFail,
	DisableFlagsInUseLine: true,
}

var deploymentArtifactsCmd = &cobra.Command{
	Use:                   "artifacts",
	Short:                 "Upload artifacts to massdriver",
	RunE:                  RunDeploymentUploadArtifacts,
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	deploymentCmd.AddCommand(deploymentProvisionCmd)
	deploymentCmd.AddCommand(deploymentDecommissionCmd)
	deploymentCmd.PersistentFlags().StringP("deployment-id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver Deployment ID")

	deploymentProvisionCmd.AddCommand(deploymentProvisionStartCmd)
	deploymentProvisionCmd.AddCommand(deploymentProvisionCompleteCmd)
	deploymentProvisionCmd.AddCommand(deploymentProvisionFailCmd)

	deploymentDecommissionCmd.AddCommand(deploymentDecommissionStartCmd)
	deploymentDecommissionCmd.AddCommand(deploymentDecommissionCompleteCmd)
	deploymentDecommissionCmd.AddCommand(deploymentDecommissionFailCmd)

	deploymentCmd.AddCommand(deploymentArtifactsCmd)
	deploymentArtifactsCmd.Flags().StringP("file", "f", "./artifacts.json", "Artifacts file")
}

func RunDeploymentProvisionStart(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while getting deployment from Massdriver")
		return err
	}
	err = mdClient.ReportProvisionStarted(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while getting deployment from Massdriver")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentProvisionComplete(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while sending provision complete event")
		return err
	}
	err = mdClient.ReportProvisionCompleted(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while sending provision complete event")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentProvisionFail(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}
	err = mdClient.ReportProvisionFailed(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentDecommissionStart(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}
	err = mdClient.ReportDecommissionStarted(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentDecommissionComplete(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}
	err = mdClient.ReportDecommissionCompleted(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentDecommissionFail(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("deployment-id")
	if id == "" {
		cmd.Help()
		fmt.Println("\nERROR: deployment-id must be set (by flags or environment variable)")
		return errors.New("deployment-id must be set (by flags or environment variable)")
	}

	log.Info().Str("deployment", id).Msg("getting deployment from Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}
	err = mdClient.ReportDecommissionFailed(id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while generating event")
		return err
	}

	log.Info().Str("deployment", id).Msg("deployment get complete")
	return nil
}

func RunDeploymentUploadArtifacts(cmd *cobra.Command, args []string) error {
	artifacts, _ := cmd.Flags().GetString("file")

	id := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")

	log.Info().Str("deployment", id).Msg("uploading artifact file to Massdriver")
	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while uploading artifact files")
		return err
	}
	err = mdClient.UploadArtifactFile(artifacts, id)
	if err != nil {
		log.Error().Err(err).Str("deployment", id).Msg("an error occurred while uploading artifact files")
		return err
	}
	log.Info().Str("deployment", id).Msg("artifact uploaded")
	return nil
}
