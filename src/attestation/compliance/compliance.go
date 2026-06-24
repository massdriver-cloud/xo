package compliance

import (
	"encoding/json"
	"fmt"

	"xo/src/attestation"

	v1 "github.com/in-toto/attestation/go/v1"
)

const PredicateType = "https://massdriver.cloud/attestations/compliance/v1"

// Predicate carries security/compliance scan results for a deployment. Results
// embed the scanners' native output (SARIF where available) so the predicate
// stays compatible with many tools without bespoke schemas.
type Predicate struct {
	attestation.DeploymentContext
	Scanners []Scanner       `json:"scanners,omitempty"`
	Summary  Summary         `json:"summary"`
	Results  json.RawMessage `json:"results,omitempty"`
}

type Scanner struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Ruleset string `json:"ruleset,omitempty"`
}

type Summary struct {
	Passed      int            `json:"passed"`
	Failed      int            `json:"failed"`
	Skipped     int            `json:"skipped"`
	LevelCounts map[string]int `json:"levelCounts,omitempty"` // by SARIF level: error/warning/note
}

// NewStatement wraps a compliance predicate in an in-toto statement whose subject
// is the deployment (identified by URI).
func NewStatement(subjectURI string, predicate Predicate) (*v1.Statement, error) {
	if subjectURI == "" {
		return nil, fmt.Errorf("subject URI is required")
	}
	if predicate.DeploymentID == "" {
		return nil, fmt.Errorf("deployment ID is required")
	}

	predicateStruct, err := attestation.StructFromJSON(predicate)
	if err != nil {
		return nil, fmt.Errorf("failed to build predicate: %w", err)
	}

	return &v1.Statement{
		Type:          attestation.StatementType,
		PredicateType: PredicateType,
		Subject:       attestation.DeploymentSubject(subjectURI),
		Predicate:     predicateStruct,
	}, nil
}
