package generic

import "testing"

func TestExtractorSubjects(t *testing.T) {
	input := `[
	  { "uri": "custom://thing/1", "name": "thing-1", "type": "custom:thing", "digest": { "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }, "attributes": { "tier": "db" } },
	  { "uri": "custom://thing/2", "name": "thing-2" }
	]`

	subjects, err := Extractor{}.Subjects([]byte(input), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %d", len(subjects))
	}

	first := subjects[0]
	if first.Digest["sha256"] != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Errorf("expected provided digest to be preserved, got %q", first.Digest["sha256"])
	}
	if got := first.Annotations.Fields["tier"].GetStringValue(); got != "db" {
		t.Errorf("expected custom attribute 'tier'='db', got %s", got)
	}
	if got := first.Annotations.Fields["md:instance"].GetStringValue(); got != "inst-1" {
		t.Errorf("expected md:instance 'inst-1', got %s", got)
	}

	// Second has no digest supplied -> synthesized identity digest.
	if len(subjects[1].Digest["sha256"]) != 64 {
		t.Errorf("expected synthesized identity digest, got %q", subjects[1].Digest["sha256"])
	}
}

func TestExtractorSubjects_EmptyInputYieldsNoSubjects(t *testing.T) {
	subjects, err := Extractor{}.Subjects(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subjects != nil {
		t.Errorf("expected nil subjects for empty input (caller falls back), got %v", subjects)
	}
}

func TestExtractorSubjects_MissingURI(t *testing.T) {
	if _, err := (Extractor{}).Subjects([]byte(`[{"name":"no-uri"}]`), nil); err == nil {
		t.Fatal("expected error when a subject lacks a uri")
	}
}
