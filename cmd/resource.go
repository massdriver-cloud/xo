package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"xo/src/bundle"
	"xo/src/resource"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

const bundleDir = "/massdriver/bundle"

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

var resourceProvisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Publishes the current step's resources and removes stale ones",
	Long: `Renders and publishes a resource for every jq template in the current step, then
deletes any resource on the instance whose field is no longer declared in massdriver.yaml
(e.g. after a bundle update removed or renamed it).

Templates are evaluated against {params, connections, envs, secrets, outputs} and are
looked up, in priority order, from:
  1. steps[].resources.<field>.template in massdriver.yaml
  2. resource_<field>.jq in the step directory
  3. artifact_<field>.jq in the step directory (legacy)`,
	RunE: runResourceProvision,
}

var resourceDecommissionCmd = &cobra.Command{
	Use:   "decommission",
	Short: "Deletes the current step's resources",
	Long: `Deletes the resources this step manages (those it has an inline or file template
for) plus any resource whose field is no longer declared in massdriver.yaml. Resources
managed by other steps are left for those steps. Run it only after the step's
decommission has succeeded.`,
	RunE: runResourceDecommission,
}

func init() {
	rootCmd.AddCommand(resourceCmd)

	resourceCmd.AddCommand(resourcePublishCmd)
	resourcePublishCmd.Flags().StringP("file", "f", "", "Path to the json formatted resource file to send (use '-' for stdin)")
	resourcePublishCmd.Flags().StringP("field", "d", "", "Resource field in the massdriver.yaml file")
	resourcePublishCmd.Flags().StringP("name", "n", "", "Human friendly name of the resource")
	resourcePublishCmd.MarkFlagRequired("file")
	resourcePublishCmd.MarkFlagRequired("field")
	resourcePublishCmd.MarkFlagRequired("name")

	resourceCmd.AddCommand(resourceDeleteCmd)
	resourceDeleteCmd.Flags().StringP("field", "d", "", "Resource field in the massdriver.yaml file")
	resourceDeleteCmd.Flags().StringP("id", "i", "", "Resource identifier")
	resourceDeleteCmd.MarkFlagRequired("field")

	resourceCmd.AddCommand(resourceProvisionCmd)
	resourceProvisionCmd.Flags().StringP("outputs", "o", "", "Path to the json formatted provisioner outputs (use '-' for stdin)")
	resourceProvisionCmd.Flags().StringP("name", "n", resource.DefaultNameFormat, "Human friendly resource name ({field} and {instance} are substituted)")
	resourceProvisionCmd.Flags().String("params", "/massdriver/params.json", "Path to params json")
	resourceProvisionCmd.Flags().String("connections", "/massdriver/connections.json", "Path to connections json")
	resourceProvisionCmd.Flags().String("envs", "/massdriver/envs.json", "Path to envs json")
	resourceProvisionCmd.Flags().String("secrets", "/massdriver/secrets.json", "Path to secrets json")
	resourceProvisionCmd.MarkFlagRequired("outputs")

	resourceCmd.AddCommand(resourceDecommissionCmd)
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
	cmd.SilenceUsage = true

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

	resourceMap := make(map[string]any)
	unmarshalErr := json.Unmarshal(resourceBytes, &resourceMap)
	if unmarshalErr != nil {
		return telemetry.LogError(span, unmarshalErr, "unable to unmarshal resource bytes")
	}

	log.Info().Msg("Publishing resource...")
	err = resource.Publish(ctx, provClient.Resources, resourceMap, field, resourceName)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while publishing resource")
	}
	log.Info().Msgf("Resource \"%s.%s\" published", provClient.Config().InstanceID, field)

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
	cmd.SilenceUsage = true

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
		id = instanceId + "." + field
	}

	log.Info().Msg("Deleting resource...")
	err = resource.Delete(ctx, provClient.Resources, id)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while deleting resource")
	}
	log.Info().Msgf("Resource \"%s\" deleted", id)

	return nil
}

// resolveStep parses massdriver.yaml and returns the bundle, the current step
// path and the step's directory, where its jq template files live.
func resolveStep(provConfig provisioning.Config) (bundle.Bundle, string, string, error) {
	stepPath := provConfig.StepPath
	if stepPath == "" {
		return bundle.Bundle{}, "", "", fmt.Errorf("MASSDRIVER_STEP_PATH is not set")
	}

	bun, err := bundle.ParseBundle(filepath.Join(bundleDir, "massdriver.yaml"))
	if err != nil {
		return bundle.Bundle{}, "", "", fmt.Errorf("unable to parse massdriver.yaml: %w", err)
	}
	return bun, stepPath, filepath.Join(bundleDir, stepPath), nil
}

func runResourceProvision(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runResourceProvision")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	outputsPath, err := cmd.Flags().GetString("outputs")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read outputs flag")
	}
	nameFormat, err := cmd.Flags().GetString("name")
	if err != nil {
		return telemetry.LogError(span, err, "unable to read name flag")
	}
	cmd.SilenceUsage = true

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}
	provConfig := provClient.Config()

	bun, stepPath, templateDir, err := resolveStep(provConfig)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while provisioning resources")
	}
	log.Info().Msgf(`Provisioning resources for instance "%s"...`, provConfig.InstanceID)

	inputs := map[string]any{}
	for _, key := range []string{"params", "connections", "envs", "secrets"} {
		path, _ := cmd.Flags().GetString(key)
		if inputs[key], err = readJSONObject(path); err != nil {
			return telemetry.LogError(span, err, fmt.Sprintf("unable to read %s", key))
		}
	}
	if inputs["outputs"], err = readJSONObject(outputsPath); err != nil {
		return telemetry.LogError(span, err, "unable to read outputs")
	}

	result, err := resource.Provision(ctx, provClient.Resources, resource.ProvisionInput{
		InstanceID:  provConfig.InstanceID,
		Bundle:      bun,
		StepPath:    stepPath,
		TemplateDir: templateDir,
		Inputs:      inputs,
		NameFormat:  nameFormat,
	})
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while provisioning resources")
	}
	log.Info().Msgf("Resource provision complete: %d published, %d deleted", len(result.Published), len(result.Deleted))

	return nil
}

func runResourceDecommission(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runResourceDecommission")
	telemetry.SetSpanAttributes(span)
	defer span.End()
	cmd.SilenceUsage = true

	provClient, err := provisioning.NewClient()
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while initializing Massdriver client")
	}
	provConfig := provClient.Config()

	bun, stepPath, templateDir, err := resolveStep(provConfig)
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while decommissioning resources")
	}
	log.Info().Msgf(`Decommissioning resources for instance "%s"...`, provConfig.InstanceID)

	result, err := resource.Decommission(ctx, provClient.Resources, resource.DecommissionInput{
		InstanceID:  provConfig.InstanceID,
		Bundle:      bun,
		StepPath:    stepPath,
		TemplateDir: templateDir,
	})
	if err != nil {
		return telemetry.LogError(span, err, "an error occurred while decommissioning resources")
	}
	log.Info().Msgf("Resource decommission complete: %d deleted", len(result.Deleted))

	return nil
}

func readJSONObject(path string) (map[string]any, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	obj := map[string]any{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("%s is not a valid JSON object: %w", path, err)
	}
	return obj, nil
}
