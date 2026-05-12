package attestation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
)

// PushStatement pushes an attestation statement as an OCI artifact referencing a subject
// repoRef: OCI repository reference (e.g., "ghcr.io/org/bundle")
// subjectDigest: digest of the bundle being attested (e.g., "sha256:abc123...")
// artifactType: OCI artifact type for the attestation
// stmt: the attestation statement to push
// annotations: optional manifest annotations
func PushStatement(ctx context.Context, repoRef string, subjectDigest string, artifactType string, stmt *Statement, annotations map[string]string) (ocispec.Descriptor, error) {
	if repoRef == "" {
		return ocispec.Descriptor{}, fmt.Errorf("repoRef is required")
	}
	if subjectDigest == "" {
		return ocispec.Descriptor{}, fmt.Errorf("subjectDigest is required")
	}
	if artifactType == "" {
		return ocispec.Descriptor{}, fmt.Errorf("artifactType is required")
	}
	if stmt == nil {
		return ocispec.Descriptor{}, fmt.Errorf("statement is required")
	}

	store := memory.New()

	stmtBytes, err := json.Marshal(stmt)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to marshal statement: %w", err)
	}

	// Create descriptor and push attestation blob
	blobDesc := content.NewDescriptorFromBytes("application/json", stmtBytes)
	blobDesc.Annotations = map[string]string{
		ocispec.AnnotationTitle: "attestation.json",
	}
	if err := store.Push(ctx, blobDesc, bytes.NewReader(stmtBytes)); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push attestation to store: %w", err)
	}

	repo, err := remote.NewRepository(repoRef)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to create repository: %w", err)
	}

	// Pack into an OCI Artifact manifest with subject
	packOpts := oras.PackManifestOptions{
		Layers:              []ocispec.Descriptor{blobDesc},
		ManifestAnnotations: annotations,
		Subject: &ocispec.Descriptor{
			MediaType: ocispec.MediaTypeImageManifest,
			Digest:    digest.Digest(subjectDigest),
		},
	}

	manifest, err := oras.PackManifest(ctx, store, oras.PackManifestVersion1_1, artifactType, packOpts)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to pack attestation manifest: %w", err)
	}

	// Copy from memory store to remote repository
	_, err = oras.Copy(ctx, store, manifest.Digest.String(), repo, manifest.Digest.String(), oras.DefaultCopyOptions)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push attestation to registry: %w", err)
	}

	return manifest, nil
}
