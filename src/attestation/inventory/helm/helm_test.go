package helm

import "testing"

const sampleManifest = `---
# Source: chart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: production
spec:
  replicas: 3
---
apiVersion: v1
kind: Service
metadata:
  name: my-app
spec:
  type: ClusterIP
---
`

func TestExtractorResources(t *testing.T) {
	resources, err := Extractor{}.Resources([]byte(sampleManifest), map[string]string{"md:instance": "inst-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Fatalf("expected 2 resources (empty trailing doc skipped), got %d", len(resources))
	}

	dep := resources[0]
	if dep.Uri != "k8s://apps/v1/Deployment/production/my-app" {
		t.Errorf("unexpected deployment uri: %s", dep.Uri)
	}
	if dep.Name != "my-app" {
		t.Errorf("expected name 'my-app', got %s", dep.Name)
	}
	if len(dep.Digest["sha256"]) != 64 {
		t.Errorf("expected sha256 digest of the object, got %q", dep.Digest["sha256"])
	}
	if got := dep.Annotations.Fields["type"].GetStringValue(); got != "k8s:apps:deployment" {
		t.Errorf("expected type 'k8s:apps:deployment', got %s", got)
	}
	if got := dep.Annotations.Fields["k8s:namespace"].GetStringValue(); got != "production" {
		t.Errorf("expected namespace annotation 'production', got %s", got)
	}
	if got := dep.Annotations.Fields["md:instance"].GetStringValue(); got != "inst-1" {
		t.Errorf("expected md:instance 'inst-1', got %s", got)
	}

	svc := resources[1]
	if got := svc.Annotations.Fields["type"].GetStringValue(); got != "k8s:core:service" {
		t.Errorf("expected core-group type 'k8s:core:service', got %s", got)
	}
}

func TestExtractorResources_InvalidYAML(t *testing.T) {
	if _, err := (Extractor{}).Resources([]byte("\tnot: [valid"), nil); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
