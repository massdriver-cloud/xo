package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"xo/src/bundle"
	"xo/src/resource"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var resourceCmd = &cobra.Command{
	Use:     "resource",
	Short:   "Resource tools",
	Long:    ``,
	Aliases: []string{"artifact"},
}

var resourcePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publishes a resource during provisioning",
	RunE:  runResourcePublish,
}

var resourceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a resource during decommission",
	RunE:  runResourceDelete,
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	resourceCmd.AddCommand(resourcePublishCmd)
	resourcePublishCmd.Flags().StringP("file", "f", "", "Path to the json formatted resource file to send (use '-' for stdin)")
	resourcePublishCmd.Flags().StringP("field", "d", "", "Resource field in the massdriver.yaml file")
	resourcePublishCmd.Flags().StringP("name", "n", "", "Human friendly name of the resource")
	resourcePublishCmd.Flags().StringP("massdriver-file", "m", "../massdriver.yaml", "Path to massdriver.yaml file")
	resourcePublishCmd.Flags().StringP("schema-file", "s", "../schema-artifacts.json", "Path to artifact schema file")
	resourcePublishCmd.MarkFlagRequired("file")
	resourcePublishCmd.MarkFlagRequired("field")
	resourcePublishCmd.MarkFlagRequired("name")

	resourceCmd.AddCommand(resourceDeleteCmd)
	resourceDeleteCmd.Flags().StringP("field", "d", "", "Resource field in the massdriver.yaml file")
	resourceDeleteCmd.Flags().StringP("id", "i", "", "Resource identifier")
	resourceDeleteCmd.Flags().StringP("massdriver-file", "m", "../massdriver.yaml", "Path to massdriver.yaml file")
	resourceDeleteCmd.MarkFlagRequired("field")
}

func runResourcePublish(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runResourcePublish")
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
			return telemetry.LogError(span, err, "unable to open resource file")
		}
		defer artFile.Close()
	}
	resourceBytes, err := io.ReadAll(artFile)
	if err != nil {
		return telemetry.LogError(span, err, "unable to read resource file")
	}

	schemasFile, err := os.Open(schemasPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open schemas file")
	}

	log.Info().Msg("Validating resource " + field + "...")
	valid, err := resource.Validate(field, resourceBytes, schemasFile)
	if !valid || err != nil {
		return telemetry.LogError(span, err, "resource is invalid")
	}
	resourceMap := make(map[string]any)
	unmarshalErr := json.Unmarshal(resourceBytes, &resourceMap)
	if unmarshalErr != nil {
		return telemetry.LogError(span, unmarshalErr, "unable to unmarshal resource bytes")
	}
	log.Info().Msg("Resource is valid!")

	log.Info().Msg("Publishing resource " + field + "...")
	bun, err := bundle.ParseBundle(massYamlPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open massdriver.yaml")
	}

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	err = resource.Publish(ctx, provClient.Resources, resourceMap, &bun, field, artName)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while publishing resource")
	}
	log.Info().Msg("Resource " + field + " published")

	return err
}

func runResourceDelete(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runResourceDelete")
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
			return telemetry.LogError(span, missingErr, "an error occurred while deleting resource")
		}
		id = packageName + "-" + field
	}

	log.Info().Msg("Deleting resource " + id + "...")

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	err = resource.Delete(ctx, provClient.Resources, id, field)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while deleting resource")
	}
	log.Info().Msg("Resource " + id + " deleted")

	return err
}
