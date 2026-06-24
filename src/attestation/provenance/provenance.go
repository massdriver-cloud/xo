package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"xo/src/attestation"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// PredicateType is SLSA provenance v1. The statement subjects are the
	// resources the deployment produced; the predicate describes how.
	PredicateType = "https://slsa.dev/provenance/v1"

	// BuildType identifies a Massdriver deployment (IaC apply) as the build.
	BuildType = "https://massdriver.cloud/deploy/v1"
	// BuilderID identifies the Massdriver orchestrator as the builder.
	BuilderID = "https://massdriver.cloud/orchestrator"
)

// Extractor turns a provisioner's state/output into SLSA provenance subjects
// (the produced resources). One implementation per IaC tool.
type Extractor interface {
	Subjects(input []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error)
}

// Predicate is a SLSA v1 provenance predicate: how a deployment was made. The
// resources it produced are the statement subjects, not part of the predicate.
type Predicate struct {
	BuildType            string
	ExternalParameters   map[string]any
	InternalParameters   map[string]any
	ResolvedDependencies []*v1.ResourceDescriptor
	Builder              string
	InvocationID         string
	StartedOn            string
	FinishedOn           string
}

func (p Predicate) toStruct() (*structpb.Struct, error) {
	buildDef := map[string]any{"buildType": p.BuildType}
	if len(p.ExternalParameters) > 0 {
		buildDef["externalParameters"] = p.ExternalParameters
	}
	if len(p.InternalParameters) > 0 {
		buildDef["internalParameters"] = p.InternalParameters
	}
	deps, err := attestation.DescriptorsToValue(p.ResolvedDependencies)
	if err != nil {
		return nil, err
	}
	if deps != nil {
		buildDef["resolvedDependencies"] = deps
	}

	metadata := map[string]any{"invocationId": p.InvocationID}
	if p.StartedOn != "" {
		metadata["startedOn"] = p.StartedOn
	}
	if p.FinishedOn != "" {
		metadata["finishedOn"] = p.FinishedOn
	}

	return structpb.NewStruct(map[string]any{
		"buildDefinition": buildDef,
		"runDetails": map[string]any{
			"builder":  map[string]any{"id": p.Builder},
			"metadata": metadata,
		},
	})
}

// NewStatement builds a SLSA provenance statement: the subjects are the produced
// resources, the predicate describes how they were built.
func NewStatement(subjects []*v1.ResourceDescriptor, predicate Predicate) (*v1.Statement, error) {
	if len(subjects) == 0 {
		return nil, fmt.Errorf("at least one subject (produced resource) is required")
	}
	if predicate.InvocationID == "" {
		return nil, fmt.Errorf("deployment ID is required")
	}

	predicateStruct, err := predicate.toStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to build predicate: %w", err)
	}

	return &v1.Statement{
		Type:          attestation.StatementType,
		PredicateType: PredicateType,
		Subject:       subjects,
		Predicate:     predicateStruct,
	}, nil
}

// NewSubject builds an in-toto subject for one produced resource. Extractors use
// this so subject construction stays uniform across provisioners.
func NewSubject(uri, name string, digest, annotations map[string]string) (*v1.ResourceDescriptor, error) {
	rd := &v1.ResourceDescriptor{Uri: uri, Name: name, Digest: digest}
	if len(annotations) > 0 {
		anyMap := make(map[string]any, len(annotations))
		for k, v := range annotations {
			anyMap[k] = v
		}
		ann, err := structpb.NewStruct(anyMap)
		if err != nil {
			return nil, fmt.Errorf("failed to build subject annotations: %w", err)
		}
		rd.Annotations = ann
	}
	return rd, nil
}

// ConfigDigest returns the sha256 of a value's canonical JSON encoding
// (encoding/json sorts map keys, so it is deterministic). Extractors use it to
// bind a subject to the resource's deploy-time configuration.
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
