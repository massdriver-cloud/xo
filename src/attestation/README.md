# Attestations

Cryptographically verifiable [in-toto](https://github.com/in-toto/attestation)
statements about deployments. Built on the canonical
`github.com/in-toto/attestation/go/v1` types. See [`DESIGN.md`](./DESIGN.md) for
the full model (three-tier subjects, storage, signing); this README is the quick
reference.

## Model in one paragraph

An attestation is a signed claim — `subject` (what it's about) + `predicate`
(the claim). Both deployment-tier attestations are published to the Massdriver
API and indexed by deployment ID. Two types:

- **Provenance** (`https://slsa.dev/provenance/v1`) — genuine SLSA provenance.
  The deployment apply is the build; its inputs (params/connections, bundle) and
  the versioned provisioner are the predicate, and the **produced cloud resources
  are the subjects** (each digested by its deploy-time config). One signed
  statement is both "how it was made" and "what it produced."
- **Compliance** (`https://massdriver.cloud/attestations/compliance/v1`) — the
  security posture at deploy time; embeds scanner output (SARIF). Subject is the
  deployment.

## Usage

Provenance has a subcommand per provisioner; each reads that tool's state/output
and emits the same SLSA provenance:

```bash
# Terraform / OpenTofu
terraform show -json > tfshow.json
xo attest provenance terraform --id "$MASSDRIVER_DEPLOYMENT_ID" --state-file ./tfshow.json

# Helm
helm get manifest my-release > manifest.yaml
xo attest provenance helm --id "$MASSDRIVER_DEPLOYMENT_ID" --manifest-file ./manifest.yaml

# Bicep (Azure deployment stacks)
az stack group show -n my-stack -g my-rg -o json > stack.json
xo attest provenance bicep --id "$MASSDRIVER_DEPLOYMENT_ID" --stack-file ./stack.json

# Generic — custom provisioner supplies its own subjects (or none)
xo attest provenance generic --id "$MASSDRIVER_DEPLOYMENT_ID" --provisioner my-tool --subjects-file ./subjects.json

# Compliance — wrap a scanner's SARIF output
xo attest compliance --id "$MASSDRIVER_DEPLOYMENT_ID" --scanner checkov --results-file ./checkov.sarif.json
```

Common deployment-context flags (`--id/--instance/--name/--version/--project/
--environment/--organization`) are shared across all `attest` commands and
default from the orchestrator environment (`MASSDRIVER_DEPLOYMENT_ID`,
`MASSDRIVER_PACKAGE_ID`, `MASSDRIVER_BUNDLE_NAME`, `MASSDRIVER_ORGANIZATION_ID`).
See `examples/attestations/` for sample inputs.

## Layout

Shared `attestation/` package:

- `statement.go` — statement/subject helpers, predicate `structpb` building
- `context.go` — shared deployment identity envelope + `DeploymentURI`
- `publish.go` — publish attestations to the Massdriver API
- `oci.go` — push attestations to OCI (for future bundle-tier attestations)

Per-type subpackages:

- `provenance/` — SLSA provenance predicate, statement builder, shared subject helpers + an `Extractor` interface, with one subpackage per provisioner:
  - `provenance/terraform`, `provenance/helm`, `provenance/bicep`, `provenance/generic`
- `compliance/` — compliance predicate, SARIF summarization

## Not yet implemented

- **Signing** — statements are not yet wrapped in a DSSE envelope / signed
  (cosign). Until then they are verifiable records, not tamper-evident.
- **API endpoints** — `publish.go` is a stub; it serializes and logs the payload
  in place of the (not-yet-existing) API call.
- **Compliance summary** — `summary` is not yet computed from SARIF.
- **Bundle-tier attestations** — SBOM and SLSA build provenance (OCI Referrers).
