package massdriver

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

var MassdriverURL = "https://api.massdriver.cloud/"

type MassdriverClient struct {
	Specification *Specification
}

type Specification struct {
	Action       string `envconfig:"ACTION"`
	BundleName   string `envconfig:"BUNDLE_NAME"`
	DeploymentID string `envconfig:"DEPLOYMENT_ID" required:"true"`
	// TODO: make this required once PACKAGE_NAME is fully deprecated
	InstanceID     string `envconfig:"INSTANCE_ID"`
	OrganizationID string `envconfig:"ORGANIZATION_ID" required:"true"`
	PackageName    string `envconfig:"PACKAGE_NAME"`
	Token          string `envconfig:"TOKEN" required:"true"`
	URL            string `envconfig:"URL"`
}

func InitializeMassdriverClient() (*MassdriverClient, error) {
	client := new(MassdriverClient)

	var specErr error
	client.Specification, specErr = GetSpecification()
	if specErr != nil {
		return nil, specErr
	}

	if client.Specification.URL == "" {
		client.Specification.URL = MassdriverURL
	}

	// TODO need to rework auth, for now just assume deployment id and token are present
	deployment_id := os.Getenv("MASSDRIVER_DEPLOYMENT_ID")
	deployment_token := os.Getenv("MASSDRIVER_TOKEN")
	if deployment_id == "" || deployment_token == "" {
		return nil, fmt.Errorf("MASSDRIVER_DEPLOYMENT_ID and MASSDRIVER_TOKEN must be set")
	}

	return client, nil
}

func GetSpecification() (*Specification, error) {
	spec := Specification{}
	err := envconfig.Process("massdriver", &spec)
	return &spec, err
}
