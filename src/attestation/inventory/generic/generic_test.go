package generic

import "testing"

func TestExtractorAssets(t *testing.T) {
	input := `[
	  { "uri": "custom://thing/1", "name": "thing-1", "type": "custom:thing", "digest": { "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }, "attributes": { "tier": "db" } },
	  { "uri": "custom://thing/2", "name": "thing-2" }
	]`

	assets, err := Extractor{}.Assets([]byte(input), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}

	first := assets[0]
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
	if len(assets[1].Digest["sha256"]) != 64 {
		t.Errorf("expected synthesized identity digest, got %q", assets[1].Digest["sha256"])
	}
}

func TestExtractorAssets_EmptyInputYieldsNoAssets(t *testing.T) {
	assets, err := Extractor{}.Assets(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assets != nil {
		t.Errorf("expected nil assets for empty input, got %v", assets)
	}
}

func TestExtractorAssets_MissingURI(t *testing.T) {
	if _, err := (Extractor{}).Assets([]byte(`[{"name":"no-uri"}]`), nil); err == nil {
		t.Fatal("expected error when an asset lacks a uri")
	}
}
