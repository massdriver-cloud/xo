package bundle_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"xo/src/bundle"
)

func TestParseBundle(t *testing.T) {
	var got, _ = bundle.ParseBundle("./testdata/bundle.yaml")
	var want = bundle.Bundle{
		Schema:      "draft-07",
		Name:        "aws-vpc",
		Description: "Something",
		Type:        "bundle",
		Artifacts:   map[string]interface{}{},
		Params: map[string]interface{}{
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":  "string",
					"title": "Name",
				},
				"age": map[string]interface{}{
					"type":  "integer",
					"title": "Age",
				},
			},
			"required": []interface{}{
				"name",
			},
		},
		Connections: map[string]interface{}{
			"required": []interface{}{
				"default",
			},
			"properties": map[string]interface{}{
				"default": map[string]interface{}{
					"type":  "string",
					"title": "Default credential",
				},
			},
		},
		Ui: map[string]interface{}{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseBundleStepResources(t *testing.T) {
	path := filepath.Join(t.TempDir(), "massdriver.yaml")
	content := `
name: pg
steps:
  - path: src
    provisioner: opentofu
    resources:
      billing:
        template: |
          {
            url: .outputs.function_url,
            region: .params.region
          }
resources:
  billing:
    resource_type: billing
artifacts:
  properties:
    legacy:
      $ref: foo
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := bundle.ParseBundle(path)
	if err != nil {
		t.Fatal(err)
	}

	step, ok := got.Step("src")
	if !ok {
		t.Fatal("expected step src")
	}
	wantTemplate := "{\n  url: .outputs.function_url,\n  region: .params.region\n}\n"
	if step.Resources["billing"].Template != wantTemplate {
		t.Errorf("got template %q, want %q", step.Resources["billing"].Template, wantTemplate)
	}
	if _, ok := got.Step("missing"); ok {
		t.Error("expected missing step to not be found")
	}

	wantFields := map[string]bool{"billing": true, "legacy": true}
	if !reflect.DeepEqual(got.ResourceFields(), wantFields) {
		t.Errorf("got fields %v, want %v", got.ResourceFields(), wantFields)
	}
}
