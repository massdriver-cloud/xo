// Package generic builds SLSA provenance subjects from a caller-supplied list,
// for provisioners that don't have a built-in parser. With no input, it yields
// no subjects and the caller falls back to the deployment as the subject.
package generic

import (
	"encoding/json"
	"fmt"

	"xo/src/attestation/provenance"

	v1 "github.com/in-toto/attestation/go/v1"
)

// subject is the Massdriver subject schema a custom provisioner emits.
type subject struct {
	URI        string            `json:"uri"`
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Digest     map[string]string `json:"digest,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Extractor reads a JSON array of subjects (the generic schema).
type Extractor struct{}

func (Extractor) Subjects(input []byte, attributes map[string]string) ([]*v1.ResourceDescriptor, error) {
	if len(input) == 0 {
		return nil, nil // no subjects supplied; caller falls back to the deployment
	}

	var items []subject
	if err := json.Unmarshal(input, &items); err != nil {
		return nil, fmt.Errorf("failed to parse subjects file: %w", err)
	}

	var subjects []*v1.ResourceDescriptor
	for _, item := range items {
		if item.URI == "" {
			return nil, fmt.Errorf("each subject requires a uri")
		}
		digest := item.Digest
		if len(digest) == 0 {
			digest = provenance.IdentityDigest(item.URI)
		}

		annotations := map[string]string{}
		if item.Type != "" {
			annotations["type"] = item.Type
		}
		for k, v := range item.Attributes {
			annotations[k] = v
		}
		for k, v := range attributes {
			annotations[k] = v
		}

		s, err := provenance.NewSubject(item.URI, item.Name, digest, annotations)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	return subjects, nil
}
