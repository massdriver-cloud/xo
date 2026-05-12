# Bundle Attestations

This directory contains the attestation implementation for Massdriver bundles. Attestations use the [in-toto attestation framework](https://github.com/in-toto/attestation) to create cryptographically verifiable statements about bundles and their deployed resources.

## Overview

Attestations are stored as OCI artifacts in the same registry as the bundle, using the [OCI Referrers API](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers) to link attestations to specific bundle versions.

## Attestation Types

### Inventory Attestation

Records the infrastructure resources created by a bundle deployment.

**Predicate Type**: `https://massdriver.cloud/attestations/inventory/v1`

**Fields**:
- `deploymentId`: Unique identifier for the deployment
- `bundleName`: Name of the bundle
- `bundleVersion`: Version of the bundle deployed
- `project`: Project name
- `environment`: Environment name (e.g., production, staging)
- `account`: Cloud account information
  - `cloud`: Cloud provider (aws, azure, gcp)
  - `accountId`: Cloud account identifier
  - `region`: Cloud region
- `resources`: Array of created resources
  - `type`: Resource type (e.g., `aws:rds:db-instance`)
  - `id`: Resource identifier
  - `name`: Resource name
  - `tags`: Key-value tags
- `generatedAt`: Timestamp when inventory was collected
- `producer`: Tool that generated the inventory
  - `tool`: Tool name (e.g., terraform, pulumi)
  - `version`: Tool version

## Usage

### Creating an Inventory Attestation

```bash
# Create an inventory JSON file (see examples/inventory/example-inventory.json)
xo attest inventory \
  --name my-bundle \
  --version v1.2.3 \
  --inventory-file ./inventory.json
```

### Inventory JSON Format

See `examples/inventory/example-inventory.json` for a complete example:

```json
{
  "deploymentId": "deploy-abc123",
  "bundleName": "my-rds-bundle",
  "bundleVersion": "v1.2.3",
  "project": "my-project",
  "environment": "production",
  "account": {
    "cloud": "aws",
    "accountId": "123456789012",
    "region": "us-east-1"
  },
  "resources": [
    {
      "type": "aws:rds:db-instance",
      "id": "my-production-db",
      "name": "production-database",
      "tags": {
        "Environment": "production"
      }
    }
  ],
  "generatedAt": "2025-12-10T00:00:00Z",
  "producer": {
    "tool": "terraform",
    "version": "1.5.0"
  }
}
```

## Architecture

### Files

- `statement.go`: Core in-toto statement structure
- `inventory.go`: Inventory predicate definition and statement creation
- `oci.go`: OCI artifact operations for pushing attestations
- `*_test.go`: Unit tests

### OCI Storage

Attestations are stored as OCI artifacts with:
- **Artifact Type**: `application/vnd.massdriver.attestation.inventory.v1+json`
- **Subject**: References the bundle manifest digest
- **Annotations**: 
  - `cloud.massdriver.attestation-type`: Type of attestation
  - `cloud.massdriver.bundle-name`: Bundle name
  - `cloud.massdriver.bundle-version`: Bundle version

## Future Enhancements

1. **Signing**: Add cryptographic signatures using Sigstore/cosign
2. **Compliance Attestations**: Add security/compliance scan results
3. **Query/Discovery**: Tools to query and filter attestations
4. **Verification**: Verify attestation integrity and signatures
5. **SBOM**: Software Bill of Materials for bundle dependencies
6. **Provenance**: Build provenance using SLSA

## Resources

- [in-toto Attestation Framework](https://github.com/in-toto/attestation)
- [in-toto Statement Spec](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md)
- [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec)
- [SLSA Provenance](https://slsa.dev/provenance/)
