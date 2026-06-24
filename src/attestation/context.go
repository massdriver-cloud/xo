package attestation

import (
	"net/url"
	"strings"
	"time"
)

// DeploymentContext is the shared identity envelope embedded in deployment-tier
// predicates that aren't SLSA-shaped (compliance). The deployment predicate
// carries the same identity in SLSA-native fields instead.
type DeploymentContext struct {
	DeploymentID string       `json:"deploymentId"`
	InstanceID   string       `json:"instanceId,omitempty"`
	Project      string       `json:"project,omitempty"`
	Environment  string       `json:"environment,omitempty"`
	Bundle       BundleRef    `json:"bundle"`
	GeneratedAt  time.Time    `json:"generatedAt"`
	Producer     ProducerInfo `json:"producer"`
}

// BundleRef links a deployment-tier attestation back to the bundle (OCI) tier.
type BundleRef struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Digest  string `json:"digest,omitempty"`
}

// ProducerInfo identifies the tool that generated an attestation.
type ProducerInfo struct {
	Tool    string `json:"tool"`
	Version string `json:"version,omitempty"`
}

// DeploymentURI builds the in-toto subject URI for a deployment. Empty segments
// are skipped so it degrades gracefully when context is only partially known.
func DeploymentURI(org, project, environment, instanceID, deploymentID string) string {
	segments := []string{}
	for _, s := range []string{org, project, environment, instanceID} {
		if s != "" {
			segments = append(segments, url.PathEscape(s))
		}
	}
	uri := "massdriver://" + strings.Join(segments, "/")
	if deploymentID != "" {
		uri += "/deployments/" + url.PathEscape(deploymentID)
	}
	return uri
}
