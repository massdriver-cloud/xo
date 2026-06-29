package inventory

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"xo/src/attestation"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// PredicateType is the Massdriver-owned inventory predicate: the cloud resources
// a deployment produced. No ratified standard covers this — CycloneDX and SPDX
// are software BOMs, not cloud-resource inventories — so we own it explicitly.
const PredicateType = "https://massdriver.cloud/attestations/inventory/v1"

// Extractor turns a provisioner's state/output into the resources a deployment
// produced. One implementation per IaC tool. The result is self-reported: it is
// "the deployment reported these resources," not independently verified.
type Extractor interface {
	Resources(input []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error)
}

// Predicate is the Massdriver inventory: the resources a deployment produced.
// It embeds the shared deployment context and carries the resource list (the
// resources are described as in-toto resource descriptors — the same shape that,
// before the provenance/inventory split, served as the provenance subjects).
type Predicate struct {
	attestation.DeploymentContext
	Provisioner string `json:"provisioner,omitempty"`
	// Resources is encoded via protojson (see NewStatement), not encoding/json,
	// so it is excluded from the struct's JSON marshaling and merged in by hand.
	Resources []*v1.ResourceDescriptor `json:"-"`
}

// NewStatement wraps an inventory predicate in an in-toto statement whose subject
// is the deployment (identified by URI) — the same subject provenance and
// compliance use, so all three anchor to the same deployment event.
func NewStatement(subjectURI string, predicate Predicate) (*v1.Statement, error) {
	if subjectURI == "" {
		return nil, fmt.Errorf("subject URI is required")
	}
	if predicate.DeploymentID == "" {
		return nil, fmt.Errorf("deployment ID is required")
	}

	// Marshal the context fields with encoding/json, then graft the resources in
	// via protojson so their digests and annotations serialize correctly.
	raw, err := json.Marshal(predicate)
	if err != nil {
		return nil, fmt.Errorf("failed to build predicate: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("failed to build predicate: %w", err)
	}

	resources, err := attestation.DescriptorsToValue(predicate.Resources)
	if err != nil {
		return nil, fmt.Errorf("failed to encode resources: %w", err)
	}
	if resources == nil {
		resources = []any{}
	}
	body["resources"] = resources

	predicateStruct, err := structpb.NewStruct(body)
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

// NewResource builds an in-toto resource descriptor for one produced resource.
// Extractors use this so resource construction stays uniform across provisioners.
func NewResource(uri, name string, digest, annotations map[string]string) (*v1.ResourceDescriptor, error) {
	rd := &v1.ResourceDescriptor{Uri: uri, Name: name, Digest: digest}
	if len(annotations) > 0 {
		anyMap := make(map[string]any, len(annotations))
		for k, v := range annotations {
			anyMap[k] = v
		}
		ann, err := structpb.NewStruct(anyMap)
		if err != nil {
			return nil, fmt.Errorf("failed to build resource annotations: %w", err)
		}
		rd.Annotations = ann
	}
	return rd, nil
}

// ConfigDigest returns the sha256 of a value's canonical JSON encoding
// (encoding/json sorts map keys, so it is deterministic). Extractors use it to
// bind a resource to its deploy-time configuration.
func ConfigDigest(v any) (map[string]string, error) {
	canonical, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(canonical)
	return map[string]string{"sha256": hex.EncodeToString(sum[:])}, nil
}

// IdentityDigest returns the sha256 of an identifier string, for resources whose
// provisioner output exposes only an id (no recoverable config).
func IdentityDigest(s string) map[string]string {
	sum := sha256.Sum256([]byte(s))
	return map[string]string{"sha256": hex.EncodeToString(sum[:])}
}
