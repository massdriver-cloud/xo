package bundle_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"xo/src/bundle"
	"xo/src/massdriver"

	"github.com/massdriver-cloud/mass/pkg/gqlmock"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/memory"
)

func TestPullV0(t *testing.T) {
	type testData struct {
		name string
		data []byte
	}
	tests := []testData{
		{
			name: "basic",
			data: []byte(`data`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gqlClient := gqlmock.NewClientWithSingleJSONResponse(map[string]interface{}{
				"data": map[string]interface{}{
					"bundleSourceCode": map[string]interface{}{
						"source": base64.StdEncoding.EncodeToString(tc.data),
					},
				},
			})

			client := massdriver.MassdriverClient{
				GQLCLient: gqlClient,
				Specification: &massdriver.Specification{
					BundleID:         "bundleuuid1",
					OrganizationUUID: "orguuid1",
				},
			}

			buf := new(bytes.Buffer)

			pullErr := bundle.PullV0(context.Background(), &client, buf)
			if pullErr != nil {
				t.Errorf("pull failed: %v", pullErr)
			}

			got := buf.String()
			want := string(tc.data)
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}

}

func TestPullV1(t *testing.T) {
	ctx := context.Background()

	// Pre-populate source store with fake files
	source := memory.New()
	target := memory.New()
	tag := "latest"

	// Simulate pushing files
	files := map[string]string{
		"massdriver.yaml":       "kind: Bundle\nname: test",
		"schema-artifacts.json": `{"type": "object"}`,
	}

	var layers []ocispec.Descriptor
	for path, data := range files {
		desc := content.NewDescriptorFromBytes("application/octet-stream", []byte(data))
		desc.Annotations = map[string]string{
			ocispec.AnnotationTitle: path,
		}
		if err := source.Push(ctx, desc, bytes.NewReader([]byte(data))); err != nil {
			t.Fatalf("failed to push %s: %v", path, err)
		}
		layers = append(layers, desc)
	}

	// Create and tag manifest
	manifest, err := oras.PackManifest(ctx, source, oras.PackManifestVersion1_1,
		"application/vnd.massdriver.bundle.v1+json", oras.PackManifestOptions{Layers: layers})
	if err != nil {
		t.Fatalf("failed to pack manifest: %v", err)
	}
	if err := source.Tag(ctx, manifest, tag); err != nil {
		t.Fatalf("failed to tag manifest: %v", err)
	}

	tests := []struct {
		name       string
		sourceRepo oras.Target
		target     oras.Target
		tag        string
		wantFiles  []string
		wantErr    bool
	}{
		{
			name:       "successful pull",
			sourceRepo: source,
			target:     target,
			tag:        tag,
			wantFiles: []string{
				"massdriver.yaml",
				"schema-artifacts.json",
			},
			wantErr: false,
		},
		{
			name:       "missing tag",
			sourceRepo: source,
			target:     memory.New(),
			tag:        "does-not-exist",
			wantFiles:  nil,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			desc, pullErr := bundle.PullV1(ctx, tc.sourceRepo, tc.target, tc.tag)
			if (pullErr != nil) != tc.wantErr {
				t.Fatalf("PullV1WithRepo() error = %v, wantErr %v", pullErr, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			// Fetch manifest and verify titles
			rc, err := tc.target.Fetch(ctx, desc)
			if err != nil {
				t.Fatalf("Fetch error: %v", err)
			}
			var manifest ocispec.Manifest
			if err := json.NewDecoder(rc).Decode(&manifest); err != nil {
				t.Fatalf("Manifest decode error: %v", err)
			}

			gotTitles := make(map[string]bool)
			for _, l := range manifest.Layers {
				gotTitles[l.Annotations[ocispec.AnnotationTitle]] = true
			}

			for _, f := range tc.wantFiles {
				if !gotTitles[f] {
					t.Errorf("expected file %q not found in pulled layers", f)
				}
			}
		})
	}
}
