package provenance

import (
	"fmt"

	"xo/src/attestation"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// PredicateType is SLSA provenance v1. The subject is the deployment itself;
	// the predicate describes how it was produced. The resources the deployment
	// produced are recorded separately as inventory, not here.
	PredicateType = "https://slsa.dev/provenance/v1"

	// BuildType identifies a Massdriver deployment (IaC apply) as the build.
	BuildType = "https://massdriver.cloud/deploy/v1"
	// BuilderID identifies the Massdriver orchestrator as the builder.
	BuilderID = "https://massdriver.cloud/orchestrator"
)

// Predicate is a SLSA v1 provenance predicate: how a deployment was made — its
// inputs, the bundle it ran, and the builder. It describes production, not the
// resources produced.
type Predicate struct {
	BuildType            string
	ExternalParameters   map[string]any
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

// NewStatement builds a SLSA provenance statement whose subject is the deployment
// (identified by URI) — the same subject compliance and inventory use. The
// predicate describes how the deployment was produced.
func NewStatement(subjectURI string, predicate Predicate) (*v1.Statement, error) {
	if subjectURI == "" {
		return nil, fmt.Errorf("subject URI is required")
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
		Subject:       attestation.DeploymentSubject(subjectURI),
		Predicate:     predicateStruct,
	}, nil
}
