package provenance

import (
	"encoding/json"
	"testing"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func sampleSubjects() []*v1.ResourceDescriptor {
	ann, _ := structpb.NewStruct(map[string]any{"type": "aws:db-instance", "md:instance": "inst-1"})
	return []*v1.ResourceDescriptor{
		{Uri: "my-db", Name: "production-db", Digest: map[string]string{"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}, Annotations: ann},
	}
}

func samplePredicate() Predicate {
	return Predicate{
		BuildType:            BuildType,
		ExternalParameters:   map[string]any{"instance": "inst-1"},
		InternalParameters:   map[string]any{"provisioner": "terraform"},
		ResolvedDependencies: []*v1.ResourceDescriptor{{Uri: "pkg:bundle/my-rds-bundle@v1.2.3"}},
		Builder:              BuilderID,
		InvocationID:         "deploy-123",
	}
}

func TestNewStatement_Success(t *testing.T) {
	stmt, err := NewStatement(sampleSubjects(), samplePredicate())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stmt.PredicateType != PredicateType {
		t.Errorf("expected PredicateType %s, got %s", PredicateType, stmt.PredicateType)
	}
	if stmt.PredicateType != "https://slsa.dev/provenance/v1" {
		t.Errorf("expected SLSA provenance predicate type, got %s", stmt.PredicateType)
	}
	if len(stmt.Subject) != 1 || stmt.Subject[0].Uri != "my-db" {
		t.Fatalf("expected the produced resource as subject, got %+v", stmt.Subject)
	}
	if err := stmt.Validate(); err != nil {
		t.Fatalf("statement failed in-toto validation: %v", err)
	}

	b, err := protojson.Marshal(stmt)
	if err != nil {
		t.Fatalf("failed to marshal statement: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	pred := m["predicate"].(map[string]any)
	if _, ok := pred["assets"]; ok {
		t.Error("predicate should not contain an assets field; produced resources are the subjects")
	}
	if _, ok := pred["buildDefinition"]; !ok {
		t.Error("predicate should contain buildDefinition")
	}
}

func TestNewStatement_ValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		subjects []*v1.ResourceDescriptor
		pred     Predicate
		wantErr  string
	}{
		{"no subjects", nil, samplePredicate(), "at least one subject (produced resource) is required"},
		{"missing deployment id", sampleSubjects(), Predicate{}, "deployment ID is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStatement(tt.subjects, tt.pred)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
