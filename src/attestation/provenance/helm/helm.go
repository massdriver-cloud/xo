// Package helm extracts SLSA provenance subjects from `helm get manifest`
// output (the rendered, multi-document Kubernetes manifest).
package helm

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"xo/src/attestation/provenance"

	v1 "github.com/in-toto/attestation/go/v1"
	yaml "gopkg.in/yaml.v3"
)

// Extractor reads the multi-document YAML emitted by `helm get manifest`.
type Extractor struct{}

func (Extractor) Subjects(manifest []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(manifest))

	var subjects []*v1.ResourceDescriptor
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

		digest, err := provenance.ConfigDigest(object)
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

		subject, err := provenance.NewSubject(k8sURI(apiVersion, kind, namespace, name), name, digest, annotations)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
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
