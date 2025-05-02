package cmd

import (
	"fmt"
	"io"
	"os"
	"xo/src/artifact"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/services/artifacts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact tools",
	Long:  ``,
}

var artifactPublishCmd = &cobra.Command{
	Use:           "publish",
	Short:         "Publishes an artifact during provisioning",
	RunE:          runArtifactPublish,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var artifactDeleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "Deletes an artifact during decommission",
	RunE:          runArtifactDelete,
	SilenceUsage:  true,
	SilenceErrors: true,
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
	artifactDeleteCmd.Flags().StringP("id", "i", "", "Artifact identifier")
	artifactDeleteCmd.Flags().StringP("massdriver-file", "m", "../massdriver.yaml", "Path to massdriver.yaml file")
	artifactDeleteCmd.MarkFlagRequired("field")
}

func runArtifactPublish(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runArtifactPublish")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	artFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read file flag")
	}
	field, err := cmd.Flags().GetString("field")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read field flag")
	}
	artName, err := cmd.Flags().GetString("name")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read name flag")
	}
	massYamlPath, err := cmd.Flags().GetString("massdriver-file")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read massdriver.yaml file flag")
	}
	schemasPath, err := cmd.Flags().GetString("schema-file")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read schema file flag")
	}

	var artFile *os.File
	if artFilePath == "-" {
		artFile = os.Stdin
	} else {
		artFile, err = os.Open(artFilePath)
		if err != nil {
			return telemetry.LogError(span, err, "unable to open artifact file")
		}
		defer artFile.Close()
	}
	artifactBytes, err := io.ReadAll(artFile)
	if err != nil {
		return telemetry.LogError(span, err, "unable to read artifact file")
	}

	schemasFile, err := os.Open(schemasPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open artifacts schemas file")
	}

	log.Info().Msg("Validating artifact " + field + "...")
	valid, err := artifact.Validate(field, artifactBytes, schemasFile)
	if !valid || err != nil {
		return telemetry.LogError(span, err, "artifact is invalid")
	}
	log.Info().Msg("Artifact is valid!")

	log.Info().Msg("Publishing artifact " + field + "...")
	bun, err := bundle.ParseBundle(massYamlPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open massdriver.yaml")
	}

	mdClient, err := client.New()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	artifactService := artifacts.NewService(mdClient)

	err = artifact.Publish(ctx, artifactService, artifactBytes, &bun, field, artName)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while publishing artifact")
	}
	log.Info().Msg("Artifact " + field + " published")

	return err
}

func runArtifactDelete(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runArtifactDelete")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read id flag")
	}
	field, err := cmd.Flags().GetString("field")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read field flag")
	}

	if id == "" {
		packageName := os.Getenv("MASSDRIVER_PACKAGE_NAME")
		if packageName == "" {
			missingErr := fmt.Errorf("id field not set and MASSDRIVER_PACKAGE_NAME environment variable is not set")
			return telemetry.LogError(span, missingErr, "an error occurred while deleting artifact")
		}
		id = packageName + "-" + field
	}

	log.Info().Msg("Deleting artifact " + id + "...")

	mdClient, err := client.New()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	artifactService := artifacts.NewService(mdClient)

	err = artifact.Delete(ctx, artifactService, id, field)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while deleting artifact")
	}
	log.Info().Msg("Artifact " + id + " deleted")

	return err
}
