// Package bicep extracts SLSA provenance subjects from `az stack <scope> show -o
// json` output (Azure deployment stacks).
package bicep

import (
	"encoding/json"
	"fmt"
	"strings"

	"xo/src/attestation/provenance"

	v1 "github.com/in-toto/attestation/go/v1"
)

// Extractor reads the JSON emitted by `az stack ... show -o json`.
type Extractor struct{}

func (Extractor) Subjects(stackJSON []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	// The managed-resources list appears at the top level on newer CLIs and under
	// `properties` on older ones; accept either.
	var stack struct {
		Resources  []json.RawMessage `json:"resources"`
		Properties struct {
			Resources []json.RawMessage `json:"resources"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(stackJSON, &stack); err != nil {
		return nil, fmt.Errorf("failed to parse az stack output: %w", err)
	}

	resources := stack.Resources
	if len(resources) == 0 {
		resources = stack.Properties.Resources
	}

	var subjects []*v1.ResourceDescriptor
	for _, raw := range resources {
		var object map[string]any
		if err := json.Unmarshal(raw, &object); err != nil {
			return nil, fmt.Errorf("failed to parse stack resource: %w", err)
		}
		id, _ := object["id"].(string)
		if id == "" {
			continue
		}

		// az stack lists resource references; the entry is the most config we
		// have, so the digest binds to whatever fields it includes.
		digest, err := provenance.ConfigDigest(object)
		if err != nil {
			return nil, fmt.Errorf("failed to digest stack resource: %w", err)
		}

		annotations := map[string]string{"type": azureType(id)}
		for k, v := range attributes {
			annotations[k] = v
		}

		subject, err := provenance.NewSubject(id, azureName(id), digest, annotations)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// azureType derives a normalized type from an Azure resource id's provider path,
// e.g. ".../providers/Microsoft.Sql/servers/foo/databases/bar" ->
// "azure:microsoft.sql:servers:databases".
func azureType(id string) string {
	marker := "/providers/"
	idx := strings.Index(id, marker)
	if idx < 0 {
		return "azure"
	}
	segments := strings.Split(strings.Trim(id[idx+len(marker):], "/"), "/")
	if len(segments) == 0 || segments[0] == "" {
		return "azure"
	}
	parts := []string{strings.ToLower(segments[0])} // provider namespace
	for i := 1; i < len(segments); i += 2 {         // type segments alternate with names
		parts = append(parts, strings.ToLower(segments[i]))
	}
	return "azure:" + strings.Join(parts, ":")
}

// azureName returns the resource name (last id segment).
func azureName(id string) string {
	parts := strings.Split(strings.TrimRight(id, "/"), "/")
	return parts[len(parts)-1]
}
