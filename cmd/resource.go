package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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

	resourceFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read file flag")
	}
	field, err := cmd.Flags().GetString("field")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read field flag")
	}
	resourceName, err := cmd.Flags().GetString("name")
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

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	var resourceFile *os.File
	if resourceFilePath == "-" {
		resourceFile = os.Stdin
	} else {
		resourceFile, err = os.Open(resourceFilePath)
		if err != nil {
			return telemetry.LogError(span, err, "unable to open resource file")
		}
		defer resourceFile.Close()
	}
	resourceBytes, err := io.ReadAll(resourceFile)
	if err != nil {
		return telemetry.LogError(span, err, "unable to read resource file")
	}

	schemasFile, err := os.Open(schemasPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open schemas file")
	}
	defer schemasFile.Close()

	log.Info().Msg("Validating resource...")
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

	log.Info().Msg("Publishing resource...")
	bun, err := bundle.ParseBundle(massYamlPath)
	if err != nil {
		return telemetry.LogError(span, err, "unable to open massdriver.yaml")
	}

	err = resource.Publish(ctx, provClient.Resources, resourceMap, &bun, field, resourceName)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while publishing resource")
	}
	log.Info().Msgf("Resource \"%s-%s\" published", provClient.Config().InstanceID, field)

	return nil
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

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}

	if id == "" {
		instanceId := provClient.Config().InstanceID
		// TODO: This can be removed once self-hosted users are updated
		if instanceId == "" {
			packageName := os.Getenv("MASSDRIVER_PACKAGE_NAME")
			if packageName == "" {
				missingErr := fmt.Errorf("id field not set and both MASSDRIVER_INSTANCE_ID and MASSDRIVER_PACKAGE_NAME environment variables are not set")
				return telemetry.LogError(span, missingErr, "an error occurred while deleting resource")
			}
			instanceId = packageName[:strings.LastIndex(packageName, "-")]
		}
		id = instanceId + "-" + field
	}

	log.Info().Msg("Deleting resource...")
	err = resource.Delete(ctx, provClient.Resources, id)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while deleting resource")
	}
	log.Info().Msgf("Resource \"%s\" deleted", id)

	return nil
}
