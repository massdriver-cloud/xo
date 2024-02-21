package cmd

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"xo/src/bundle"
	"xo/src/env"
	"xo/src/telemetry"
	"xo/src/util"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var envsCmd = &cobra.Command{
	Use:   "envs",
	Short: "Manage environment variables",
	Long:  ``,
}

var envsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get environment variables",
	RunE:  runEnvsGet,
}

func init() {
	rootCmd.AddCommand(envsCmd)

	envsCmd.AddCommand(envsGetCmd)
	envsGetCmd.Flags().StringP("output", "o", "", "Output file path")
	envsGetCmd.Flags().StringP("massdriver-file", "m", "./bundle/massdriver.yaml", "Path to massdriver.yaml file")
	envsGetCmd.Flags().StringP("params-file", "p", "./params.json", "Path to params.json file")
	envsGetCmd.Flags().StringP("connections-file", "c", "./connections.json", "Path to connections.json file")
}

func runEnvsGet(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runEnvsGet")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	out, _ := cmd.Flags().GetString("output")
	massdriverPath, _ := cmd.Flags().GetString("massdriver-file")
	paramsPath, _ := cmd.Flags().GetString("params-file")
	connectionsPath, _ := cmd.Flags().GetString("connections-file")

	bun, err := bundle.ParseBundle(massdriverPath)
	if err != nil {
		util.LogError(err, span, "an error occurred while parsing massdriver.yaml")
		return err
	}

	params := map[string]interface{}{}
	paramsFile, err := os.ReadFile(paramsPath)
	if err != nil {
		util.LogError(err, span, "an error occurred while attempting to read params.json")
		return err
	}
	err = json.Unmarshal(paramsFile, &params)
	if err != nil {
		util.LogError(err, span, "an error occurred while parsing params.json")
		return err
	}

	connections := map[string]interface{}{}
	connectionsFile, err := os.ReadFile(connectionsPath)
	if err != nil {
		util.LogError(err, span, "an error occurred while attempting to read params.json")
		return err
	}
	err = json.Unmarshal(connectionsFile, &connections)
	if err != nil {
		util.LogError(err, span, "an error occurred while parsing params.json")
		return err
	}

	log.Info().Msg("Extracting environment variables...")

	envs, err := env.GenerateEnvs(ctx, bun.App.Envs, params, connections)
	if err != nil {
		util.LogError(err, span, "an error occurred while extracting environment variables")
		return err
	}

	log.Info().Msg("Environment variables extracted, writing to output.")

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

	jsonBytes, marshalErr := json.MarshalIndent(envs, "", "  ")
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msg("an error occurred marshaling environment variables to JSON")
		span.RecordError(marshalErr)
		span.SetStatus(codes.Error, marshalErr.Error())
		return marshalErr
	}
	_, writeErr := output.Write(jsonBytes)
	if writeErr != nil {
		log.Error().Err(writeErr).Msg("an error occurred while writing environment variables to file")
		span.RecordError(writeErr)
		span.SetStatus(codes.Error, writeErr.Error())
		return writeErr
	}

	log.Info().Msg("Environment variables written successfully!")

	return nil
}
