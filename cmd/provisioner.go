package cmd

import (
	"errors"
	"io"
	"os"
	"xo/src/massdriver"
	"xo/src/provisioners"
	"xo/src/provisioners/helm"
	"xo/src/provisioners/opa"
	tf "xo/src/provisioners/terraform"
	"xo/src/telemetry"
	"xo/src/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var provisionerCmd = &cobra.Command{
	Use:   "provisioner",
	Short: "Manage provisioners",
	Long:  ``,
}

var provisionerAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate secure AWS credential file for provisioning",
	Long:  `This command will generate a set of AWS credentials in ini format which can be passed to the actual provisioning step. These credentials would be narrowly scoped to just this provisioning run so the bundle can't access unauthorized data.`,
	RunE:  runProvisionerAuth,
}

var provisionerHelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Commands specific to helm provisioner",
	Long:  ``,
}

var provisionerHelmReportCmd = &cobra.Command{
	Use:   "report-progress",
	Short: "Report Helm provisioner progress to Massdriver",
	Long:  ``,
	RunE:  runProvisionerHelmReport,
}

var provisionerOPACmd = &cobra.Command{
	Use:   "opa",
	Short: "Commands specific to opa provisioner",
	Long:  ``,
}

var provisionerOPAReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report opa results to Massdriver",
	Long:  ``,
	RunE:  runProvisionerOPAReport,
}

var provisionerTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Commands specific to terraform provisioner",
	Long:  ``,
}

var provisionerTerraformReportCmd = &cobra.Command{
	Use:   "report-progress",
	Short: "Report Terraform provisioner progress to Massdriver",
	Long:  ``,
	RunE:  runProvisionerTerraformReport,
}

var provisionerTerraformBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Generate a terraform backend config",
	Long:  ``,
}

var provisionerTerraformBackendS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Generate a terraform s3 backend config",
	Long:  ``,
	RunE:  runProvisionerTerraformBackendS3,
}

var provisionerTerraformBackendHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Generate a terraform HTTP backend config",
	Long:  ``,
	RunE:  runProvisionerTerraformBackendHTTP,
}

func init() {
	rootCmd.AddCommand(provisionerCmd)

	provisionerCmd.AddCommand(provisionerAuthCmd)
	provisionerAuthCmd.PersistentFlags().StringP("role", "r", os.Getenv("MASSDRIVER_PROVISIONER_ROLE_ARN"), "AWS Role ARN to assume for provisioning (custom policy will be generated)")
	provisionerAuthCmd.PersistentFlags().StringP("external-id", "d", os.Getenv("MASSDRIVER_PROVISIONER_ROLE_EXTERNAL_ID"), "External ID to use when assuming the provisioner role")
	provisionerAuthCmd.PersistentFlags().StringP("access-key", "a", os.Getenv("MASSDRIVER_PROVISIONER_WORKFLOW_ACCESS_KEY_ID"), "AWS Access Key ID to use for assuming provisioner role")
	provisionerAuthCmd.PersistentFlags().StringP("secret-key", "s", os.Getenv("MASSDRIVER_PROVISIONER_WORKFLOW_SECRET_ACCESS_KEY"), "AWS Secret Access Key for assuming the provisioner role")
	provisionerAuthCmd.PersistentFlags().StringP("output", "o", "", "Output file path")

	provisionerCmd.AddCommand(provisionerOPACmd)
	provisionerOPACmd.AddCommand(provisionerOPAReportCmd)
	provisionerOPAReportCmd.Flags().StringP("file", "f", "", "File to extract ('-' for stdin)")
	provisionerOPAReportCmd.MarkFlagRequired("file")

	provisionerCmd.AddCommand(provisionerHelmCmd)
	provisionerHelmCmd.AddCommand(provisionerHelmReportCmd)
	provisionerHelmReportCmd.Flags().StringP("file", "f", "", "File to extract ('-' for stdin)")
	provisionerHelmReportCmd.MarkFlagRequired("file")

	provisionerCmd.AddCommand(provisionerTerraformCmd)

	provisionerTerraformCmd.AddCommand(provisionerTerraformReportCmd)
	provisionerTerraformReportCmd.Flags().StringP("file", "f", "", "File to extract ('-' for stdin)")
	provisionerTerraformReportCmd.MarkFlagRequired("file")

	provisionerTerraformCmd.AddCommand(provisionerTerraformBackendCmd)
	provisionerTerraformBackendCmd.PersistentFlags().StringP("output", "o", "./backend.tf.json", "Output file path")
	provisionerTerraformBackendCmd.AddCommand(provisionerTerraformBackendS3Cmd)
	provisionerTerraformBackendS3Cmd.Flags().StringP("step", "s", "", "Bundle Step")
	provisionerTerraformBackendS3Cmd.MarkFlagRequired("step")
	provisionerTerraformBackendCmd.AddCommand(provisionerTerraformBackendHTTPCmd)
	provisionerTerraformBackendHTTPCmd.Flags().StringP("step", "s", "", "Bundle Step")
	provisionerTerraformBackendHTTPCmd.MarkFlagRequired("step")
}

func runProvisionerAuth(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerAuth")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	out, _ := cmd.Flags().GetString("output")
	roleArn, _ := cmd.Flags().GetString("role")
	externalId, _ := cmd.Flags().GetString("external-id")
	accessKey, _ := cmd.Flags().GetString("access-key")
	secretKey, _ := cmd.Flags().GetString("secret-key")

	if roleArn == "" {
		err := errors.New("role ARN is empty (nothing in MASSDRIVER_PROVISIONER_ROLE_ARN environment variable)")
		util.LogError(err, span, "error while generating provisioner auth")
		return err
	}

	var output io.Writer
	if out == "" {
		output = os.Stdout
	} else {
		if _, err := os.Stat(out); errors.Is(err, os.ErrNotExist) {
			outputFile, fileErr := os.Create(out)
			if fileErr != nil {
				log.Error().Err(fileErr).Msg("an error occurred while creating file")
				span.RecordError(fileErr)
				span.SetStatus(codes.Error, fileErr.Error())
				return fileErr
			}
			defer outputFile.Close()
			output = outputFile
		}
	}

	spec, specErr := massdriver.GetSpecification()
	if specErr != nil {
		log.Error().Err(specErr).Msg("an error occurred while extracting Massdriver specification")
		span.RecordError(specErr)
		span.SetStatus(codes.Error, specErr.Error())
		return specErr
	}

	log.Info().Msg("Generating secure AWS credentials for provisioning...")

	var cfg aws.Config
	var cfgErr error

	if len(accessKey) > 0 && len(secretKey) > 0 {
		log.Info().Msg("Using credentials that were passed to perform AssumeRole")
		cfg, cfgErr = config.LoadDefaultConfig(ctx,
			config.WithRegion("us-west-2"),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					accessKey, secretKey, "",
				),
			),
		)
	} else {
		log.Info().Msg("No credentials passed, falling back to default credentials for AssumeRole")
		cfg, cfgErr = config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	}
	if cfgErr != nil {
		log.Error().Err(cfgErr).Msg("an error occured while configuring AWS client")
		span.RecordError(cfgErr)
		span.SetStatus(codes.Error, cfgErr.Error())
		return cfgErr
	}

	stsClient := sts.NewFromConfig(cfg)

	genErr := provisioners.GenerateProvisionerAWSCredentials(ctx, output, stsClient, spec, roleArn, externalId)
	if genErr != nil {
		log.Error().Err(genErr).Msg("an error occured while generating secure credentials for provisioning")
		span.RecordError(genErr)
		span.SetStatus(codes.Error, genErr.Error())
		return genErr
	}

	log.Info().Msg("Credentials generated.")

	return nil
}

func runProvisionerHelmReport(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerHelmReport")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	deploymentId := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentId == "" {
		log.Warn().Msg("Deployment ID is empty (nothing in MASSDRIVER_DEPLOYMENT_ID environment variable)")
	}

	var input io.Reader
	if file == "-" {
		input = os.Stdin
	} else {
		inputFile, err := os.Open(file)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while opening file")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		defer inputFile.Close()
		input = inputFile
	}

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while initializing Massdriver client")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = helm.ReportProgressFromLogs(ctx, mdClient, deploymentId, input)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while reporting progress")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func runProvisionerOPAReport(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerOPAReport")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	deploymentId := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentId == "" {
		log.Warn().Msg("Deployment ID is empty (nothing in MASSDRIVER_DEPLOYMENT_ID environment variable)")
	}

	var input io.Reader
	if file == "-" {
		input = os.Stdin
	} else {
		inputFile, err := os.Open(file)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while opening file")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		defer inputFile.Close()
		input = inputFile
	}

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while initializing Massdriver client")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = opa.ReportResults(ctx, mdClient, deploymentId, input)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while reporting progress")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func runProvisionerTerraformReport(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerTerraformReport")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	deploymentId := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentId == "" {
		log.Warn().Msg("Deployment ID is empty (nothing in MASSDRIVER_DEPLOYMENT_ID environment variable)")
	}

	var input io.Reader
	if file == "-" {
		input = os.Stdin
	} else {
		inputFile, err := os.Open(file)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while opening file")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		defer inputFile.Close()
		input = inputFile
	}

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while initializing Massdriver client")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = tf.ReportProgressFromLogs(ctx, mdClient, deploymentId, input)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while reporting progress")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func runProvisionerTerraformBackendS3(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerTerraformBackendS3")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	output, _ := cmd.Flags().GetString("output")
	step, _ := cmd.Flags().GetString("step")

	spec, specErr := massdriver.GetSpecification()
	if specErr != nil {
		log.Error().Err(specErr).Msg("an error occurred while extracting Massdriver specification")
		span.RecordError(specErr)
		span.SetStatus(codes.Error, specErr.Error())
		return specErr
	}

	log.Info().Msg("Generating state file")

	return tf.GenerateBackendS3File(ctx, output, spec, step)
}

func runProvisionerTerraformBackendHTTP(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runProvisionerTerraformBackendHTTP")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	output, _ := cmd.Flags().GetString("output")
	step, _ := cmd.Flags().GetString("step")

	spec, specErr := massdriver.GetSpecification()
	if specErr != nil {
		log.Error().Err(specErr).Msg("an error occurred while extracting Massdriver specification")
		span.RecordError(specErr)
		span.SetStatus(codes.Error, specErr.Error())
		return specErr
	}

	log.Info().Msg("Generating state file")

	return tf.GenerateBackendHTTPFile(ctx, output, spec, step)
}
