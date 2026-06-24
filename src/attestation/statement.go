package attestation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	v1 "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// StatementType is the in-toto v1 statement type URI.
const StatementType = v1.StatementTypeUri

// DeploymentSubject builds the single-subject slice for a deployment-tier
// statement. The in-toto v1 spec requires every subject to carry a digest; a
// deployment has no byte content, so we digest its canonical URI to bind the
// statement to that exact deployment.
func DeploymentSubject(uri string) []*v1.ResourceDescriptor {
	sum := sha256.Sum256([]byte(uri))
	return []*v1.ResourceDescriptor{{
		Uri:    uri,
		Digest: map[string]string{"sha256": hex.EncodeToString(sum[:])},
	}}
}

// StructFromJSON converts a JSON-serializable value into a protobuf Struct for
// use as an in-toto predicate body.
func StructFromJSON(v any) (*structpb.Struct, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}

// DescriptorsToValue marshals resource descriptors via protojson into a slice of
// plain maps, so they can be embedded inside a predicate Struct.
func DescriptorsToValue(ds []*v1.ResourceDescriptor) ([]any, error) {
	if len(ds) == 0 {
		return nil, nil
	}
	out := make([]any, 0, len(ds))
	for _, d := range ds {
		b, err := protojson.Marshal(d)
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}
