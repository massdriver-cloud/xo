package attestation

import (
	"encoding/json"
	"testing"
)

func TestStatement_Marshal(t *testing.T) {
	stmt := Statement{
		Type:          StatementTypeV1,
		PredicateType: "https://example.com/predicate/v1",
		Subject: []Subject{
			{
				Name: "test-artifact",
				Digest: map[string]string{
					"sha256": "abc123",
				},
			},
		},
		Predicate: json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(stmt)
	if err != nil {
		t.Fatalf("failed to marshal statement: %v", err)
	}

	var unmarshaled Statement
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal statement: %v", err)
	}

	if unmarshaled.Type != stmt.Type {
		t.Errorf("expected Type %s, got %s", stmt.Type, unmarshaled.Type)
	}
	if unmarshaled.PredicateType != stmt.PredicateType {
		t.Errorf("expected PredicateType %s, got %s", stmt.PredicateType, unmarshaled.PredicateType)
	}
	if len(unmarshaled.Subject) != len(stmt.Subject) {
		t.Errorf("expected %d subjects, got %d", len(stmt.Subject), len(unmarshaled.Subject))
	}
}

func TestSubject_RequiredFields(t *testing.T) {
	subject := Subject{
		Name: "artifact",
		Digest: map[string]string{
			"sha256": "abc123",
		},
	}

	data, err := json.Marshal(subject)
	if err != nil {
		t.Fatalf("failed to marshal subject: %v", err)
	}

	// Verify that name and digest are present (not omitted)
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, ok := result["name"]; !ok {
		t.Error("name field should be present")
	}
	if _, ok := result["digest"]; !ok {
		t.Error("digest field should be present")
	}
}
