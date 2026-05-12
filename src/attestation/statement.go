package attestation

import "encoding/json"

// Subject identifies the artifact(s) being attested to
type Subject struct {
	Name   string            `json:"name"`
	Digest map[string]string `json:"digest"` // e.g. {"sha256": "abcd..."}
}

// Statement represents an in-toto attestation statement v1
// Spec: https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md
type Statement struct {
	Type          string          `json:"_type"`
	Subject       []Subject       `json:"subject"`
	PredicateType string          `json:"predicateType"`
	Predicate     json.RawMessage `json:"predicate"`
}

const (
	// StatementTypeV1 is the in-toto statement type identifier
	StatementTypeV1 = "https://in-toto.io/Statement/v1"
)
