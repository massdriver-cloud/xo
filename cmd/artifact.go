package cmd

import (
	"fmt"
	"os"
	"xo/src/artifact"
	"xo/src/bundle"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact tools",
	Long:  ``,
}

var artifactPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes an artifact during provisioning",
	RunE:  runArtifactPublish,
}

var artifactDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an artifact during decommission",
	RunE:  runArtifactDelete,
}

func init() {
	rootCmd.AddCommand(artifactCmd)

	artifactCmd.AddCommand(artifactPublishCmd)
	artifactPublishCmd.Flags().StringP("file", "f", "", "Path to the artifact file to send (use '-' for stdin)")
	artifactPublishCmd.Flags().StringP("field", "d", "", "Artifact field in the massdriver.yaml file")
	artifactPublishCmd.Flags().StringP("name", "n", "", "Human friendly name of the artifact")
	artifactPublishCmd.Flags().StringP("massdriver-file", "m", "../massdriver.yaml", "Path to massdriver.yaml file")
	artifactPublishCmd.Flags().StringP("schema-file", "s", "../schema-artifacts.json", "Path to artifact schema file")
	artifactPublishCmd.MarkFlagRequired("file")
	artifactPublishCmd.MarkFlagRequired("field")
	artifactPublishCmd.MarkFlagRequired("name")

	artifactCmd.AddCommand(artifactDeleteCmd)
	artifactDeleteCmd.Flags().StringP("field", "d", "", "Artifact field in the massdriver.yaml file")
	artifactDeleteCmd.Flags().StringP("name", "n", "", "Human friendly name of the artifact")
	artifactDeleteCmd.Flags().StringP("massdriver-file", "m", "../massdriver.yaml", "Path to massdriver.yaml file")
	artifactDeleteCmd.MarkFlagRequired("field")
	artifactDeleteCmd.MarkFlagRequired("name")
}

func runArtifactPublish(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runArtifactPublish")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	artFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	field, err := cmd.Flags().GetString("field")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	artName, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	massYamlPath, err := cmd.Flags().GetString("massdriver-file")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	schemasPath, err := cmd.Flags().GetString("schema-file")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var artFile *os.File
	if artFilePath == "-" {
		artFile = os.Stdin
	} else {
		artFile, err = os.Open(artFilePath)
		if err != nil {
			fmt.Println(err)
		}
		defer artFile.Close()
	}

	schemasFile, err := os.Open(schemasPath)
	if err != nil {
		log.Error().Err(err).Msg("unable to open artifacts schemas file")
		return err
	}

	log.Info().Msg("Validating artifact " + field + "...")
	valid, err := artifact.Validate(field, artFile, schemasFile)
	if !valid || err != nil {
		log.Error().Err(err).Msg("artifact is invalid")
		return err
	}
	log.Info().Msg("Artifact is valid!")

	log.Info().Msg("Publishing artifact " + field + "...")
	bun, err := bundle.ParseBundle(massYamlPath)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while opening massdriver.yaml")
		return err
	}

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while initializing Massdriver client")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = artifact.Publish(ctx, mdClient, artFile, &bun, field, artName)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while publishing artifact")
		return err
	}
	log.Info().Msg("Artifact " + field + " published")

	return err
}

func runArtifactDelete(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runArtifactDelete")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	field, err := cmd.Flags().GetString("field")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while deleting artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	artName, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while deleting artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	massYamlPath, err := cmd.Flags().GetString("massdriver-file")
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while deleting artifact")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Info().Msg("Deleting artifact " + field + "...")
	bun, err := bundle.ParseBundle(massYamlPath)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while opening massdriver.yaml")
		return err
	}

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while initializing Massdriver client")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = artifact.Delete(ctx, mdClient, &bun, field, artName)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while deleting artifact")
		return err
	}
	log.Info().Msg("Artifact " + field + " deleted")

	return err
}
