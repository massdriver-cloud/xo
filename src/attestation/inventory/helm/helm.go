// Package helm extracts inventory resources from `helm get manifest`
// output (the rendered, multi-document Kubernetes manifest).
package helm

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"xo/src/attestation/inventory"

	v1 "github.com/in-toto/attestation/go/v1"
	yaml "gopkg.in/yaml.v3"
)

// Extractor reads the multi-document YAML emitted by `helm get manifest`.
type Extractor struct{}

func (Extractor) Resources(manifest []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(manifest))

	var resources []*v1.ResourceDescriptor
	for {
		var object map[string]any
		err := decoder.Decode(&object)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse helm manifest: %w", err)
		}
		if len(object) == 0 {
			continue // empty document between `---` separators
		}

		kind, _ := object["kind"].(string)
		apiVersion, _ := object["apiVersion"].(string)
		var name, namespace string
		if meta, ok := object["metadata"].(map[string]any); ok {
			name, _ = meta["name"].(string)
			namespace, _ = meta["namespace"].(string)
		}
		if kind == "" || name == "" {
			continue
		}

		digest, err := inventory.ConfigDigest(object)
		if err != nil {
			return nil, fmt.Errorf("failed to digest manifest object: %w", err)
		}

		annotations := map[string]string{"type": k8sType(apiVersion, kind)}
		if namespace != "" {
			annotations["k8s:namespace"] = namespace
		}
		for k, v := range attributes {
			annotations[k] = v
		}

		resource, err := inventory.NewResource(k8sURI(apiVersion, kind, namespace, name), name, digest, annotations)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// k8sType derives a normalized type from a Kubernetes object's group + kind,
// e.g. apiVersion "apps/v1" kind "Deployment" -> "k8s:apps:deployment";
// apiVersion "v1" kind "Service" -> "k8s:core:service".
func k8sType(apiVersion, kind string) string {
	group := "core"
	if i := strings.Index(apiVersion, "/"); i >= 0 {
		group = apiVersion[:i]
	}
	return "k8s:" + strings.ToLower(group) + ":" + strings.ToLower(kind)
}

// k8sURI builds a stable identifier for a Kubernetes object.
func k8sURI(apiVersion, kind, namespace, name string) string {
	parts := []string{apiVersion, kind}
	if namespace != "" {
		parts = append(parts, namespace)
	}
	parts = append(parts, name)
	return "k8s://" + strings.Join(parts, "/")
}
