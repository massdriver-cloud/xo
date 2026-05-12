package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"xo/src/attestation"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/client"
	sdkbundle "github.com/massdriver-cloud/massdriver-sdk-go/massdriver/platform/bundle"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"oras.land/oras-go/v2/content/file"
)

var attestCmd = &cobra.Command{
	Use:   "attest",
	Short: "Attest to the state of a bundle",
	Long:  `Create and publish attestations for bundles, such as inventory and compliance reports`,
}

var attestInventoryCmd = &cobra.Command{
	Use:     "inventory",
	Short:   "Attest to the inventory of a bundle",
	Long:    `Create an in-toto attestation statement for bundle inventory and push it to the OCI registry`,
	Aliases: []string{"inv"},
	RunE:    runInventoryAttest,
}

func init() {
	rootCmd.AddCommand(attestCmd)

	attestCmd.AddCommand(attestInventoryCmd)
	attestInventoryCmd.Flags().StringP("version", "v", "", "Bundle version")
	attestInventoryCmd.Flags().StringP("name", "n", "", "Bundle name")
	attestInventoryCmd.Flags().StringP("inventory-file", "f", "", "Path to inventory JSON file")
	viper.BindPFlag("bundle.version", attestInventoryCmd.Flags().Lookup("version"))
	viper.BindPFlag("bundle.name", attestInventoryCmd.Flags().Lookup("name"))
	viper.BindPFlag("inventory.file", attestInventoryCmd.Flags().Lookup("inventory-file"))
}

func runInventoryAttest(cmd *cobra.Command, args []string) error {
	ctx, span := otel.Tracer("xo").Start(telemetry.GetContextWithTraceParentFromEnv(), "runInventoryAttest")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	// Get bundle name and version
	bundleName := viper.GetString("bundle.name")
	if bundleName == "" {
		return telemetry.LogError(span, fmt.Errorf("required flag bundleName must be set via flag or environment variable"), "bundle name is required")
	}
	bundleVersion := viper.GetString("bundle.version")
	if bundleVersion == "" {
		return telemetry.LogError(span, fmt.Errorf("required flag bundleVersion must be set via flag or environment variable"), "bundle version is required")
	}

	// Get inventory file path
	inventoryFile := viper.GetString("inventory.file")
	if inventoryFile == "" {
		return telemetry.LogError(span, fmt.Errorf("required flag inventory-file must be set"), "inventory file is required")
	}

	log.Info().Msgf("creating inventory attestation for bundle %s:%s", bundleName, bundleVersion)

	// Read inventory data from file
	inventoryData, err := os.ReadFile(inventoryFile)
	if err != nil {
		return telemetry.LogError(span, err, "failed to read inventory file")
	}

	var inventoryPred attestation.InventoryPredicate
	if err := json.Unmarshal(inventoryData, &inventoryPred); err != nil {
		return telemetry.LogError(span, err, "failed to parse inventory JSON")
	}

	// Set metadata if not already present
	if inventoryPred.GeneratedAt.IsZero() {
		inventoryPred.GeneratedAt = time.Now()
	}
	if inventoryPred.Producer.Tool == "" {
		inventoryPred.Producer.Tool = "xo"
		inventoryPred.Producer.Version = "1.0.0" // TODO: get from build info
	}

	// Initialize Massdriver client
	mdClient, clientErr := client.New()
	if clientErr != nil {
		return telemetry.LogError(span, clientErr, "failed to create massdriver client")
	}

	// Get bundle repository
	repo, repoErr := sdkbundle.GetBundleRepository(mdClient, bundleName)
	if repoErr != nil {
		return telemetry.LogError(span, repoErr, "failed to get bundle repository")
	}

	// Pull bundle to get its digest
	fileStore, fileErr := file.New("bundle")
	if fileErr != nil {
		return telemetry.LogError(span, fileErr, "failed to create file store")
	}
	defer fileStore.Close()

	log.Info().Msgf("pulling bundle to get digest...")
	desc, pullErr := bundle.Pull(ctx, repo, fileStore, bundleVersion)
	if pullErr != nil {
		return telemetry.LogError(span, pullErr, "failed to pull bundle")
	}

	bundleDigest := desc.Digest.String()
	log.Debug().Msgf("bundle digest: %s", bundleDigest)

	// Create attestation statement
	stmt, stmtErr := attestation.NewStatementFromInventory(bundleDigest, inventoryPred)
	if stmtErr != nil {
		return telemetry.LogError(span, stmtErr, "failed to create attestation statement")
	}

	// Use bundle name as repository reference
	// The SDK should provide the full registry path via GetBundleRepository
	// For now, we'll need to extract it or pass it explicitly
	// TODO: Update when SDK provides a way to get the full repository reference
	repoRef := bundleName

	// Push attestation to registry
	log.Info().Msg("pushing inventory attestation...")
	attestDesc, pushErr := attestation.PushStatement(
		ctx,
		repoRef,
		bundleDigest,
		attestation.ArtifactTypeInventory,
		stmt,
		map[string]string{
			"cloud.massdriver.attestation-type": "inventory",
			"cloud.massdriver.bundle-name":      bundleName,
			"cloud.massdriver.bundle-version":   bundleVersion,
		},
	)
	if pushErr != nil {
		return telemetry.LogError(span, pushErr, "failed to push attestation")
	}

	log.Info().Msgf("inventory attestation published: %s", attestDesc.Digest.String())

	return nil
}
