package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets",
	Long:  ``,
}

var secretsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secrets",
	RunE:  runSecretsGet,
}

func init() {
	rootCmd.AddCommand(secretsCmd)

	secretsCmd.AddCommand(secretsGetCmd)
	secretsGetCmd.Flags().StringP("output", "o", "", "Output file path")
}

func runSecretsGet(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runSecretsGet")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	out, _ := cmd.Flags().GetString("output")

	mdClient, err := massdriver.InitializeMassdriverClient()
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while sending deployment status event: deploymentStatus")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Info().Msg("Fetching secrets.")

	secrets, err := mdClient.GetSecrets(ctx)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while fetching secrets")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	log.Info().Msg("Secrets retreived, writing to output.")

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

	jsonBytes, marshalErr := json.MarshalIndent(secrets, "", "  ")
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msg("an error occurred marshaling secrets to JSON")
		span.RecordError(marshalErr)
		span.SetStatus(codes.Error, marshalErr.Error())
		return marshalErr
	}
	_, writeErr := output.Write(jsonBytes)
	if writeErr != nil {
		log.Error().Err(writeErr).Msg("an error occurred while writing secrets to file")
		span.RecordError(writeErr)
		span.SetStatus(codes.Error, writeErr.Error())
		return writeErr
	}

	log.Info().Msg("Secrets written successfully!")

	return nil
}
