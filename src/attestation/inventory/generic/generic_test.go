package generic

import "testing"

func TestExtractorResources(t *testing.T) {
	input := `[
	  { "uri": "custom://thing/1", "name": "thing-1", "type": "custom:thing", "digest": { "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }, "attributes": { "tier": "db" } },
	  { "uri": "custom://thing/2", "name": "thing-2" }
	]`

	resources, err := Extractor{}.Resources([]byte(input), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	first := resources[0]
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
	if len(resources[1].Digest["sha256"]) != 64 {
		t.Errorf("expected synthesized identity digest, got %q", resources[1].Digest["sha256"])
	}
}

func TestExtractorResources_EmptyInputYieldsNoResources(t *testing.T) {
	resources, err := Extractor{}.Resources(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resources != nil {
		t.Errorf("expected nil resources for empty input, got %v", resources)
	}
}

func TestExtractorResources_MissingURI(t *testing.T) {
	if _, err := (Extractor{}).Resources([]byte(`[{"name":"no-uri"}]`), nil); err == nil {
		t.Fatal("expected error when a resource lacks a uri")
	}
}
