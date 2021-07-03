package cmd

import (
	"fmt"
	"io/ioutil"
	tf "xo/src/provisioners/terraform"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// converts at schema.inputs.json into terraform variables.tf file

var provisionerCmd = &cobra.Command{
	Use:   "provisioner",
	Short: "Manage provisioners",
	Long:  ``,
}

var provisionerTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Commands specific to terraform provisioner",
	Long:  ``,
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

var provisionerCompileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile JSON Schema to provisioner variable definition file.",
	Long:  ``,
	RunE:  runProvisionerCompile,
}

func init() {
	rootCmd.AddCommand(provisionerCmd)
	provisionerCmd.AddCommand(provisionerTerraformCmd)

	provisionerTerraformCmd.AddCommand(provisionerTerraformBackendCmd)
	provisionerTerraformBackendCmd.PersistentFlags().StringP("output", "o", "./backend.tf.json", "Output file path")
	provisionerTerraformBackendCmd.AddCommand(provisionerTerraformBackendS3Cmd)
	provisionerTerraformBackendS3Cmd.Flags().StringP("bucket", "b", "", "S3 bucket (required)")
	provisionerTerraformBackendS3Cmd.Flags().StringP("mrn", "m", "", "Massdriver MRN (required)")
	provisionerTerraformBackendS3Cmd.Flags().StringP("region", "r", "us-west-2", "AWS Region")
	provisionerTerraformBackendS3Cmd.Flags().StringP("dynamodb-table", "d", "xo-terraform-lock-table", "DynamoDB state lock table")
	provisionerTerraformBackendS3Cmd.Flags().StringP("shared-credentials-file", "s", "/secrets/xo.aws.creds", "Shared credentials file path")
	provisionerTerraformBackendS3Cmd.Flags().StringP("profile", "p", "xo-iac", "Name of AWS profile")
	provisionerTerraformBackendS3Cmd.MarkFlagRequired("bucket")
	provisionerTerraformBackendS3Cmd.MarkFlagRequired("mrn")

	provisionerTerraformCmd.AddCommand(provisionerCompileCmd)
	provisionerCompileCmd.Flags().StringP("schema", "s", "./schema.json", "Path to JSON Schema")
	provisionerCompileCmd.Flags().StringP("output", "o", "./variables.tf.json", "Output path. Use - for STDOUT")
}

func runProvisionerCompile(cmd *cobra.Command, args []string) error {
	var compiled string
	var err error

	provisioner := cmd.Parent().Use
	schema, _ := cmd.Flags().GetString("schema")
	outputPath, _ := cmd.Flags().GetString("output")

	log.Debug().
		Str("provisioner", provisioner).
		Str("schemaPath", schema).Msg("Compiling schema.")

	switch provisioner {
	case "terraform":
		compiled, err = tf.Compile(schema)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("unsupported argument %s the single argument 'terraform' is supported", provisioner)
		log.Error().Err(err).Msg("Compilation failed.")
		return err
	}

	return writeVariableFile(compiled, outputPath)
}

func writeVariableFile(compiled string, outPath string) error {
	switch outPath {
	case "-":
		fmt.Println(compiled)
		return nil
	default:
		data := []byte(compiled)
		err := ioutil.WriteFile(outPath, data, 0644)
		return err
	}
}

func runProvisionerTerraformBackendS3(cmd *cobra.Command, args []string) error {

	output, _ := cmd.Flags().GetString("output")
	bucket, _ := cmd.Flags().GetString("bucket")
	mrn, _ := cmd.Flags().GetString("mrn")
	region, _ := cmd.Flags().GetString("region")
	dynamoDbTable, _ := cmd.Flags().GetString("dynamodb-table")
	sharedCredentialsFile, _ := cmd.Flags().GetString("shared-credentials-file")
	profile, _ := cmd.Flags().GetString("profile")

	log.Debug().
		Str("provisioner", "terraform").
		Str("output", output).
		Str("bucket", bucket).
		Str("mrn", mrn).
		Str("region", region).
		Str("dynamodb-table", dynamoDbTable).
		Str("shared-credentials-file", sharedCredentialsFile).
		Str("profile", profile).Msg("Generating state file")

	return tf.GenerateBackendS3File(output, bucket, mrn, region, dynamoDbTable, sharedCredentialsFile, profile)
}
