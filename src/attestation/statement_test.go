package attestation

import (
	"testing"

	v1 "github.com/in-toto/attestation/go/v1"
)

func TestStructFromJSON(t *testing.T) {
	type sample struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	s, err := StructFromJSON(sample{A: "x", B: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Fields["a"].GetStringValue() != "x" {
		t.Errorf("expected a=x, got %v", s.Fields["a"])
	}
	if s.Fields["b"].GetNumberValue() != 2 {
		t.Errorf("expected b=2, got %v", s.Fields["b"])
	}
}

func TestDescriptorsToValue(t *testing.T) {
	got, err := DescriptorsToValue([]*v1.ResourceDescriptor{
		{Uri: "pkg:bundle/foo@v1", Digest: map[string]string{"sha256": "abc"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 descriptor, got %d", len(got))
	}
	m := got[0].(map[string]any)
	if m["uri"] != "pkg:bundle/foo@v1" {
		t.Errorf("expected uri preserved, got %v", m["uri"])
	}

	if v, err := DescriptorsToValue(nil); err != nil || v != nil {
		t.Errorf("expected nil for empty input, got %v (err %v)", v, err)
	}
}

func TestDeploymentSubject(t *testing.T) {
	subs := DeploymentSubject("massdriver://deploy-1")
	if len(subs) != 1 || subs[0].Uri != "massdriver://deploy-1" {
		t.Fatalf("unexpected subject: %+v", subs)
	}
	if len(subs[0].Digest["sha256"]) != 64 {
		t.Errorf("expected a sha256 digest of the uri, got %q", subs[0].Digest["sha256"])
	}
}

func TestDeploymentURI(t *testing.T) {
	tests := []struct {
		name string
		args [5]string
		want string
	}{
		{"full context", [5]string{"org-1", "proj", "prod", "inst-1", "deploy-9"}, "massdriver://org-1/proj/prod/inst-1/deployments/deploy-9"},
		{"partial context skips empty segments", [5]string{"", "", "", "inst-1", "deploy-9"}, "massdriver://inst-1/deployments/deploy-9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeploymentURI(tt.args[0], tt.args[1], tt.args[2], tt.args[3], tt.args[4])
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
