package cmd

import (
	"io"
	"os"
	"xo/src/provisioners"
	tf "xo/src/provisioners/terraform"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var provisionerCmd = &cobra.Command{
	Use:   "provisioner",
	Short: "Manage provisioners",
	Long:  ``,
}

var provisionerAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate auth file(s) for provisioners",
	Long:  ``,
	RunE:  runProvisionerAuth,
}

var provisionerTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Commands specific to terraform provisioner",
	Long:  ``,
}

var provisionerTerraformReportCmd = &cobra.Command{
	Use:   "report-progress",
	Short: "Report provisioner progress to Massdriver",
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

func init() {
	rootCmd.AddCommand(provisionerCmd)

	provisionerCmd.AddCommand(provisionerAuthCmd)
	provisionerAuthCmd.PersistentFlags().StringP("schema", "s", "schema-connections.json", "Connections schema file")
	provisionerAuthCmd.PersistentFlags().StringP("connections", "c", "connections.tf.json", "Connections json file")
	provisionerAuthCmd.PersistentFlags().StringP("output", "o", "./auth", "Output dir path")

	provisionerCmd.AddCommand(provisionerTerraformCmd)

	provisionerTerraformCmd.AddCommand(provisionerTerraformReportCmd)
	provisionerTerraformReportCmd.Flags().StringP("file", "f", "", "File to extract ('-' for stdin)")
	provisionerTerraformReportCmd.MarkFlagRequired("file")

	provisionerTerraformCmd.AddCommand(provisionerTerraformBackendCmd)
	provisionerTerraformBackendCmd.PersistentFlags().StringP("output", "o", "./backend.tf.json", "Output file path")
	provisionerTerraformBackendCmd.AddCommand(provisionerTerraformBackendS3Cmd)
	provisionerTerraformBackendS3Cmd.Flags().StringP("bucket", "b", "", "S3 bucket (required)")
	provisionerTerraformBackendS3Cmd.Flags().StringP("key", "k", "", "Path to the state file inside the S3 Bucket (required)")
	provisionerTerraformBackendS3Cmd.Flags().StringP("region", "r", "us-west-2", "AWS Region")
	provisionerTerraformBackendS3Cmd.Flags().StringP("dynamodb-table", "d", "", "DynamoDB state lock table")
	provisionerTerraformBackendS3Cmd.Flags().StringP("shared-credentials-file", "s", "", "Shared credentials file path")
	provisionerTerraformBackendS3Cmd.Flags().StringP("profile", "p", "", "Name of AWS profile")
	provisionerTerraformBackendS3Cmd.MarkFlagRequired("bucket")
	provisionerTerraformBackendS3Cmd.MarkFlagRequired("key")
}

func runProvisionerAuth(cmd *cobra.Command, args []string) error {
	connectionsPath, _ := cmd.Flags().GetString("connections")
	schemaPath, _ := cmd.Flags().GetString("schema")
	authPath, _ := cmd.Flags().GetString("output")

	log.Info().Msg("Generating auth files...")
	if _, err := os.Stat(authPath); os.IsNotExist(err) {
		err := os.Mkdir(authPath, 0777)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while creating auth directory")
			return err
		}
	}
	err := provisioners.GenerateAuthFiles(schemaPath, connectionsPath, authPath)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while generating auth files")
		return err
	}
	log.Info().Msg("Auth files generated")
	return nil
}

func runProvisionerTerraformReport(cmd *cobra.Command, args []string) error {
	file, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	deploymentId := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentId == "" {
		log.Warn().Msg("Deployment ID is empty (pulled from MASSDRIVER_DEPLOYMENT_ID environment variable)")
	}
	deploymentToken := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	if deploymentToken == "" {
		log.Warn().Msg("Deployment token is empty (pulled from MASSDRIVER_DEPLOYMENT_TOKEN environment variable)")
	}

	var input io.Reader
	if file == "-" {
		input = os.Stdin
	} else {
		inputFile, err := os.Open(file)
		if err != nil {
			log.Error().Err(err).Msg("an error occurred while opening file")
			return err
		}
		defer inputFile.Close()
		input = inputFile
	}

	return tf.ReportProgressFromLogs(deploymentId, deploymentToken, input)
}

func runProvisionerTerraformBackendS3(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	bucket, _ := cmd.Flags().GetString("bucket")
	key, _ := cmd.Flags().GetString("key")
	region, _ := cmd.Flags().GetString("region")
	dynamoDbTable, _ := cmd.Flags().GetString("dynamodb-table")
	sharedCredentialsFile, _ := cmd.Flags().GetString("shared-credentials-file")
	profile, _ := cmd.Flags().GetString("profile")

	log.Info().
		Str("provisioner", "terraform").
		Str("output", output).
		Str("bucket", bucket).
		Str("key", key).
		Str("region", region).
		Str("dynamodb-table", dynamoDbTable).
		Str("shared-credentials-file", sharedCredentialsFile).
		Str("profile", profile).Msg("Generating state file")

	return tf.GenerateBackendS3File(output, bucket, key, region, dynamoDbTable, sharedCredentialsFile, profile)
}
