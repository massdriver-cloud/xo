# Attestation Design

How Massdriver produces signed, verifiable reports about bundles and their
deployments: the trust mechanism, the three-tier model, where each attestation
is stored, and the predicate schemas. Built on the
[in-toto attestation v1](https://github.com/in-toto/attestation) types
(`github.com/in-toto/attestation/go/v1`).

## 1. Mental model

An **attestation** is a signed, machine-readable claim:

> "I, _producer_, claim _predicate_ about _subject_."

Everything we ship is the **same envelope** with a different payload — one
attestation pipeline, N predicate schemas.

| Term | Meaning | in-toto field |
|---|---|---|
| **Attestation** | The signed statement (the envelope) | `Statement` (wrapped in a DSSE envelope once signed) |
| **Subject** | *What* the claim is about, identified by uri + digest | `subject[]` (`ResourceDescriptor`) |
| **Predicate** | The content of the claim | `predicate` |
| **Predicate type** | URI naming which kind of claim this is | `predicateType` |

## 2. Three-tier model

Three entities organize everything — they determine where each attestation is
stored and how it is queried.

| Tier | Mutability | Lives in | Carries | Linking |
|---|---|---|---|---|
| **Bundle** | Immutable per version | OCI registry | SBOM, SLSA build provenance | OCI Referrers API |
| **Instance** | Mutable / long-lived | Postgres | *nothing* — the **grouping key** | parent FK |
| **Deployment** | Immutable event | Postgres | **provenance, compliance** | FK keyed by deployment ID |

In Massdriver terms:

- **Bundle** — the reusable template artifact. Same contents wherever deployed.
- **Instance** — a long-lived entity (a deployed component in a target). One
  instance has 1..N deployments over its lifetime.
- **Deployment** — a single apply against an instance. Immutable once complete.

Deployment-tier attestations are stored in Postgres and indexed by deployment
ID, with the **instance as the grouping dimension**: "how did compliance evolve
for instance X" = list X's deployments and diff their attestations. What each
attestation names as its in-toto *subject* differs by type — see §4.

## 3. OCI vs. Postgres

The attestation format and signing are storage-agnostic; OCI is one
delivery/indexing option, not a requirement.

- **Bundle-tier** (SBOM, SLSA build provenance) — the subject is the bundle (an
  OCI artifact), so these attach to the bundle in OCI via the Referrers API.
- **Deployment-tier** (provenance, compliance) — published to the Massdriver API,
  indexed by deployment ID, queryable by instance ID. OCI is not involved.

## 4. Subjects

The in-toto v1 spec requires every subject to carry a digest.

- **Provenance** is SLSA, so its subjects are the **resources the deployment
  produced** — each an in-toto `ResourceDescriptor` whose `uri` is the cloud
  resource id and whose `sha256` digest is the hash of that resource's
  deploy-time configuration. The digest binds the attestation to that exact
  config; it is *not* recomputable from the live cloud resource — an inherent
  limit of provenance over non-content-addressable infrastructure. The deployment
  itself is identified inside the predicate (`runDetails.metadata.invocationId`
  and `externalParameters`).
- **Compliance** is a Massdriver predicate; its subject is the **deployment**,
  identified by URI with a digest of that URI string:

  ```
  massdriver://<org>/<project>/<env>/<instance-id>/deployments/<deployment-id>
  ```

- **Bundle-tier** subjects use the bundle manifest digest directly.

## 5. Predicates

Two deployment-tier reports: **provenance** (how it was made + what it produced)
and **compliance** (how secure it was).

### 5.1 Provenance

Predicate type: `https://slsa.dev/provenance/v1` (genuine SLSA provenance v1).

A deployment apply *is* the build: its inputs are params/connections and the
bundle; its builder is the (versioned) provisioner; its outputs are the cloud
resources. Those resources are the statement **subjects** (§4) — SLSA records
outputs in the subject, not the predicate. Because this is real SLSA, security
teams can consume it as such, and Massdriver can target a SLSA *level* once
signing lands (§6).

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "uri": "arn:aws:rds:...:db:prod",
      "name": "production-database",
      "digest": { "sha256": "<hash of the resource's deploy-time config>" },
      "annotations": { "type": "aws:db-instance", "md:instance": "inst-7f3a", "md:project": "my-project" }
    }
  ],
  "predicateType": "https://slsa.dev/provenance/v1",
  "predicate": {
    "buildDefinition": {
      "buildType": "https://massdriver.cloud/deploy/v1",
      "externalParameters": { "instance": "inst-7f3a", "project": "my-project", "environment": "production" },
      "internalParameters": { "provisioner": "terraform" },
      "resolvedDependencies": [ { "uri": "pkg:bundle/my-rds-bundle@v1.2.3", "digest": { "sha256": "..." } } ]
    },
    "runDetails": {
      "builder": { "id": "https://massdriver.cloud/orchestrator" },
      "metadata": { "invocationId": "deploy-abc123", "finishedOn": "2026-06-22T00:00:00Z" }
    }
  }
}
```

- The subjects are the low-level cloud objects created — distinct from the
  platform's `resource` (the high-level output of an instance). They are a
  point-in-time record of what the apply produced, not a live tracker; drift and
  post-apply mutations are out of scope.
- Each subject's `annotations` carry **Massdriver-assigned** metadata (the
  normalized `type` and `md:*` attributes), never scraped cloud tags/labels —
  credentials and cloud/account context are not assumed.
- `type` is normalized from the raw provisioner type (`aws_db_instance` →
  `aws:db-instance`).
- A deploy that produces no managed resources falls back to a single subject: the
  deployment itself (URI + digest of that URI).
- Subjects are extracted per provisioner (Terraform, Helm, Bicep, or a generic
  caller-supplied list); only this extraction step is provisioner-specific — the
  predicate and envelope are identical across tools. See §7.

### 5.2 Compliance

Predicate type: `https://massdriver.cloud/attestations/compliance/v1`

Compatible with multiple scanners (Checkov, Wiz, Snyk, tfsec, …) by embedding
[SARIF](https://sarifweb.azurewebsites.net/), the format scanners already emit.
The raw SARIF is embedded verbatim; `scanners` and `summary` are derived from it
(via `github.com/owenrumney/go-sarif`) for fast querying. `summary.levelCounts`
buckets findings by the standard SARIF level (error/warning/note) — tool-specific
severity scores are not interpreted. Suppressed results count as skipped. It
embeds the shared deployment context.

```json
{
  "deploymentId": "deploy-abc123",
  "instanceId": "inst-7f3a",
  "project": "my-project",
  "environment": "production",
  "bundle": { "name": "my-rds-bundle", "version": "v1.2.3", "digest": "sha256:..." },
  "generatedAt": "2026-06-22T00:00:00Z",
  "producer": { "tool": "xo" },
  "scanners": [ { "name": "checkov", "version": "3.2.0" } ],
  "summary": { "passed": 142, "failed": 3, "skipped": 5, "levelCounts": { "error": 1, "warning": 2 } },
  "results": "<embedded SARIF, or a digest+URI pointer to it>"
}
```

## 6. Signing (not yet built)

An unsigned attestation is a verifiable record, not tamper-evident. The trust
layer wraps each `Statement` in a **DSSE** envelope signed with
**Sigstore/cosign** (optionally logged to Rekor). Signing is storage-agnostic, so
deployment-tier attestations in Postgres are signed exactly like bundle-tier
ones in OCI.

## 7. Implementation

Shared package `attestation/` (used by both attestation types):

- `statement.go` — in-toto statement helpers: subject construction, predicate
  `structpb` building (statements serialize via `protojson`).
- `context.go` — shared deployment identity envelope and `DeploymentURI`.
- `publish.go` — publishes attestations to the Massdriver API. **Stub:** the API
  endpoints do not exist yet, so it serializes and logs the payload.
- `oci.go` — pushes attestations to OCI (for bundle-tier attestations).

Per-type subpackages:

- `attestation/provenance/` — SLSA provenance predicate, statement builder, and
  shared subject helpers (`NewSubject`, `ConfigDigest`, `IdentityDigest`) plus an
  `Extractor` interface (`bytes → subjects`). Subject extraction is
  provisioner-pluggable; one subpackage per IaC tool:
  - `provenance/terraform` — `terraform show -json` (via `hashicorp/terraform-json`)
  - `provenance/helm` — `helm get manifest` (rendered Kubernetes objects)
  - `provenance/bicep` — `az stack ... show -o json` (Azure deployment stacks)
  - `provenance/generic` — caller-supplied subjects (custom provisioners), or none
- `attestation/compliance/` — compliance predicate and SARIF summarization.

Commands: `xo attest provenance <terraform|helm|bicep|generic>` (each reads that
provisioner's state/output and emits SLSA provenance) and `xo attest compliance`
(scanner results → compliance attestation). See `examples/attestations/`.

## 8. References

- [in-toto Attestation Framework](https://github.com/in-toto/attestation) ·
  [Statement v1](https://github.com/in-toto/attestation/blob/main/spec/v1/statement.md) ·
  [ResourceDescriptor](https://github.com/in-toto/attestation/blob/main/spec/v1/resource_descriptor.md)
- [SLSA Provenance](https://slsa.dev/provenance/) (deployment + bundle tiers)
- [SARIF](https://sarifweb.azurewebsites.net/)
- [terraform-json](https://pkg.go.dev/github.com/hashicorp/terraform-json) — types for `terraform show -json`
- [Sigstore / cosign](https://docs.sigstore.dev/) · [DSSE](https://github.com/secure-systems-lab/dsse)
- [OCI Distribution Spec — Referrers API](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers)
