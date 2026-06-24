package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"xo/src/attestation"
	"xo/src/attestation/compliance"
	"xo/src/attestation/provenance"
	"xo/src/attestation/provenance/bicep"
	"xo/src/attestation/provenance/generic"
	"xo/src/attestation/provenance/helm"
	"xo/src/attestation/provenance/terraform"
	"xo/src/telemetry"

	v1 "github.com/in-toto/attestation/go/v1"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "Attest to the state of a deployment",
	Long:  `Create and publish attestations for a deployment, such as provenance and compliance reports`,
}

var attestProvenanceCmd = &cobra.Command{
	Use:   "provenance",
	Short: "Attest to the provenance of a deployment",
	Long:  `Create a SLSA provenance attestation. Use the per-provisioner subcommand matching the IaC tool that performed the apply.`,
}

var attestProvenanceTerraformCmd = &cobra.Command{
	Use:     "terraform",
	Short:   "Provenance from `terraform show -json` output",
	Aliases: []string{"opentofu", "tofu"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProvenance(cmd, "terraform", terraform.Extractor{}, "state-file", true)
	},
}

var attestProvenanceHelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Provenance from `helm get manifest` output",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProvenance(cmd, "helm", helm.Extractor{}, "manifest-file", true)
	},
}

var attestProvenanceBicepCmd = &cobra.Command{
	Use:   "bicep",
	Short: "Provenance from `az stack ... show -o json` output",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProvenance(cmd, "bicep", bicep.Extractor{}, "stack-file", true)
	},
}

var attestProvenanceGenericCmd = &cobra.Command{
	Use:   "generic",
	Short: "Provenance for a custom provisioner (subjects supplied directly, or none)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("provisioner")
		if name == "" {
			name = "generic"
		}
		return runProvenance(cmd, name, generic.Extractor{}, "subjects-file", false)
	},
}

var attestComplianceCmd = &cobra.Command{
	Use:   "compliance",
	Short: "Attest to the compliance posture of a deployment",
	Long:  `Create a compliance attestation from scanner results and publish it to Massdriver`,
	RunE:  runComplianceAttest,
}

func init() {
	rootCmd.AddCommand(attestCmd)

	// Shared deployment context. Defaults come from the orchestrator's environment.
	attestCmd.PersistentFlags().StringP("id", "i", os.Getenv("MASSDRIVER_DEPLOYMENT_ID"), "Massdriver deployment ID")
	attestCmd.PersistentFlags().String("instance", os.Getenv("MASSDRIVER_PACKAGE_ID"), "Massdriver instance (package) ID")
	attestCmd.PersistentFlags().String("organization", os.Getenv("MASSDRIVER_ORGANIZATION_ID"), "Massdriver organization ID")
	attestCmd.PersistentFlags().String("project", "", "Project name")
	attestCmd.PersistentFlags().String("environment", "", "Environment name")
	attestCmd.PersistentFlags().StringP("name", "n", os.Getenv("MASSDRIVER_BUNDLE_NAME"), "Bundle name")
	attestCmd.PersistentFlags().StringP("version", "v", "", "Bundle version")

	attestCmd.AddCommand(attestProvenanceCmd)
	attestProvenanceCmd.AddCommand(attestProvenanceTerraformCmd)
	attestProvenanceCmd.AddCommand(attestProvenanceHelmCmd)
	attestProvenanceCmd.AddCommand(attestProvenanceBicepCmd)
	attestProvenanceCmd.AddCommand(attestProvenanceGenericCmd)

	attestProvenanceTerraformCmd.Flags().StringP("state-file", "s", "", "Path to `terraform show -json` output")
	attestProvenanceHelmCmd.Flags().StringP("manifest-file", "m", "", "Path to `helm get manifest` output")
	attestProvenanceBicepCmd.Flags().String("stack-file", "", "Path to `az stack ... show -o json` output")
	attestProvenanceGenericCmd.Flags().StringP("subjects-file", "f", "", "Path to a JSON array of produced resources (optional)")
	attestProvenanceGenericCmd.Flags().String("provisioner", "", "Name of the custom provisioner (recorded in the provenance)")

	attestCmd.AddCommand(attestComplianceCmd)
	attestComplianceCmd.Flags().StringP("results-file", "f", "", "Path to scanner results (SARIF/JSON)")
	attestComplianceCmd.Flags().String("scanner", "", "Scanner name (e.g. checkov)")
}

// deploymentContextFromFlags assembles the shared identity envelope and the
// deployment subject URI from the persistent attest flags / environment.
func deploymentContextFromFlags(cmd *cobra.Command) (attestation.DeploymentContext, string, error) {
	deploymentID, _ := cmd.Flags().GetString("id")
	if deploymentID == "" {
		return attestation.DeploymentContext{}, "", fmt.Errorf("deployment ID is required (--id or MASSDRIVER_DEPLOYMENT_ID)")
	}

	instanceID, _ := cmd.Flags().GetString("instance")
	org, _ := cmd.Flags().GetString("organization")
	project, _ := cmd.Flags().GetString("project")
	environment, _ := cmd.Flags().GetString("environment")
	bundleName, _ := cmd.Flags().GetString("name")
	bundleVersion, _ := cmd.Flags().GetString("version")

	dctx := attestation.DeploymentContext{
		DeploymentID: deploymentID,
		InstanceID:   instanceID,
		Project:      project,
		Environment:  environment,
		Bundle:       attestation.BundleRef{Name: bundleName, Version: bundleVersion},
		GeneratedAt:  time.Now(),
		Producer:     attestation.ProducerInfo{Tool: "xo"},
	}

	subjectURI := attestation.DeploymentURI(org, project, environment, instanceID, deploymentID)
	return dctx, subjectURI, nil
}

// runProvenance extracts subjects with the given provisioner extractor and
// publishes a SLSA provenance attestation. inputFlag names the file flag holding
// the provisioner's state/output; requireInput errors if it is unset.
func runProvenance(cmd *cobra.Command, provisioner string, extractor provenance.Extractor, inputFlag string, requireInput bool) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvenance")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	dctx, subjectURI, ctxErr := deploymentContextFromFlags(cmd)
	if ctxErr != nil {
		return telemetry.LogError(span, ctxErr, "invalid deployment context")
	}

	var input []byte
	if path, _ := cmd.Flags().GetString(inputFlag); path != "" {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return telemetry.LogError(span, readErr, "failed to read input file")
		}
		input = data
	} else if requireInput {
		return telemetry.LogError(span, fmt.Errorf("required flag --%s must be set", inputFlag), "input file is required")
	}

	subjects, extractErr := extractor.Subjects(input, massdriverAttributes(dctx))
	if extractErr != nil {
		return telemetry.LogError(span, extractErr, "failed to extract provenance subjects")
	}
	// A deploy that produced no enumerable resources still gets a record; the
	// deployment itself stands in as the subject.
	if len(subjects) == 0 {
		log.Warn().Msg("no resources extracted; recording provenance with the deployment as subject")
		subjects = attestation.DeploymentSubject(subjectURI)
	}

	pred := provenance.Predicate{
		BuildType:            provenance.BuildType,
		ExternalParameters:   externalParameters(dctx),
		InternalParameters:   map[string]any{"provisioner": provisioner},
		ResolvedDependencies: bundleDependency(dctx.Bundle),
		Builder:              provenance.BuilderID,
		InvocationID:         dctx.DeploymentID,
		FinishedOn:           dctx.GeneratedAt.UTC().Format(time.RFC3339),
	}

	stmt, stmtErr := provenance.NewStatement(subjects, pred)
	if stmtErr != nil {
		return telemetry.LogError(span, stmtErr, "failed to create provenance statement")
	}

	log.Info().Msgf("publishing provenance attestation for deployment %s via %s (%d subjects)", dctx.DeploymentID, provisioner, len(subjects))
	if pubErr := attestation.NewPublisher().Publish(ctx, dctx.DeploymentID, stmt); pubErr != nil {
		return telemetry.LogError(span, pubErr, "failed to publish attestation")
	}

	return nil
}

func runComplianceAttest(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runComplianceAttest")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	dctx, subjectURI, ctxErr := deploymentContextFromFlags(cmd)
	if ctxErr != nil {
		return telemetry.LogError(span, ctxErr, "invalid deployment context")
	}

	resultsFile, _ := cmd.Flags().GetString("results-file")
	if resultsFile == "" {
		return telemetry.LogError(span, fmt.Errorf("required flag --results-file must be set"), "results file is required")
	}
	scannerName, _ := cmd.Flags().GetString("scanner")

	resultsData, readErr := os.ReadFile(resultsFile)
	if readErr != nil {
		return telemetry.LogError(span, readErr, "failed to read results file")
	}

	// Derive scanners and a summary from the SARIF; embed the original bytes.
	scanners, summary, sumErr := compliance.SummarizeSARIF(resultsData)
	if sumErr != nil {
		return telemetry.LogError(span, sumErr, "failed to summarize scanner results")
	}
	if len(scanners) == 0 && scannerName != "" {
		scanners = []compliance.Scanner{{Name: scannerName}}
	}

	pred := compliance.Predicate{
		DeploymentContext: dctx,
		Scanners:          scanners,
		Summary:           summary,
		Results:           json.RawMessage(resultsData),
	}

	stmt, stmtErr := compliance.NewStatement(subjectURI, pred)
	if stmtErr != nil {
		return telemetry.LogError(span, stmtErr, "failed to create compliance statement")
	}

	log.Info().Msgf("publishing compliance attestation for deployment %s", dctx.DeploymentID)
	if pubErr := attestation.NewPublisher().Publish(ctx, dctx.DeploymentID, stmt); pubErr != nil {
		return telemetry.LogError(span, pubErr, "failed to publish attestation")
	}

	return nil
}

// massdriverAttributes are the platform-assigned attributes stamped onto every
// subject, in place of scraping cloud tags/labels.
func massdriverAttributes(dctx attestation.DeploymentContext) map[string]string {
	attrs := map[string]string{}
	if dctx.InstanceID != "" {
		attrs["md:instance"] = dctx.InstanceID
	}
	if dctx.Project != "" {
		attrs["md:project"] = dctx.Project
	}
	if dctx.Environment != "" {
		attrs["md:environment"] = dctx.Environment
	}
	return attrs
}

func externalParameters(dctx attestation.DeploymentContext) map[string]any {
	params := map[string]any{}
	if dctx.InstanceID != "" {
		params["instance"] = dctx.InstanceID
	}
	if dctx.Project != "" {
		params["project"] = dctx.Project
	}
	if dctx.Environment != "" {
		params["environment"] = dctx.Environment
	}
	return params
}

func bundleDependency(b attestation.BundleRef) []*v1.ResourceDescriptor {
	if b.Name == "" {
		return nil
	}
	uri := "pkg:bundle/" + b.Name
	if b.Version != "" {
		uri += "@" + b.Version
	}
	rd := &v1.ResourceDescriptor{Uri: uri}
	if b.Digest != "" {
		rd.Digest = map[string]string{"sha256": strings.TrimPrefix(b.Digest, "sha256:")}
	}
	return []*v1.ResourceDescriptor{rd}
}
