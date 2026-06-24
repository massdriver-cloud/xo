package attestation

import (
	"context"
	"fmt"

	v1 "github.com/in-toto/attestation/go/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

// Publisher publishes deployment-tier attestations to the Massdriver API,
// indexed by deployment ID.
//
// NOTE: the API does not expose attestation endpoints yet. This is a stub that
// serializes and logs the payload in place of the real POST. When the endpoints
// land, wire the Massdriver client in here and replace the log call.
type Publisher struct {
	// client *client.Client // TODO(api): real Massdriver client once endpoints exist
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

// Publish sends a (eventually signed) attestation statement for a deployment.
func (p *Publisher) Publish(ctx context.Context, deploymentID string, stmt *v1.Statement) error {
	if deploymentID == "" {
		return fmt.Errorf("deploymentID is required")
	}
	if stmt == nil {
		return fmt.Errorf("statement is required")
	}

	// TODO(signing): wrap the statement in a DSSE envelope and sign (cosign)
	// before publishing — see DESIGN.md.

	payload, err := protojson.Marshal(stmt)
	if err != nil {
		return fmt.Errorf("failed to marshal attestation: %w", err)
	}

	// TODO(api): POST the signed attestation to Massdriver, keyed by deploymentID.
	log.Info().
		Str("deploymentId", deploymentID).
		Str("predicateType", stmt.PredicateType).
		Int("bytes", len(payload)).
		Msg("publishing attestation to Massdriver API (stub — endpoint not yet implemented)")
	log.Debug().RawJSON("attestation", payload).Msg("attestation payload")

	return nil
}
