package massdriver

import (
	"fmt"
	"net/url"
	"os"
	"xo/src/api"

	"github.com/Khan/genqlient/graphql"
	"github.com/kelseyhightower/envconfig"
)

var MassdriverURL = "https://api.massdriver.cloud/"

type MassdriverClient struct {
	GQLCLient     graphql.Client
	Specification *Specification
}

type Specification struct {
	Action           string `envconfig:"ACTION"`
	BundleID         string `envconfig:"BUNDLE_ID" required:"true"`
	BundleName       string `envconfig:"BUNDLE_NAME"`
	BundleType       string `envconfig:"BUNDLE_TYPE"`
	DeploymentID     string `envconfig:"DEPLOYMENT_ID" required:"true"`
	ManifestID       string `envconfig:"MANIFEST_ID"`
	OrganizationUUID string `envconfig:"ORGANIZATION_UUID" required:"true"`
	OrganizationID   string `envconfig:"ORGANIZATION_ID" required:"true"`
	PackageID        string `envconfig:"PACKAGE_ID" required:"true"`
	PackageName      string `envconfig:"PACKAGE_NAME" required:"true"`
	TargetMode       string `envconfig:"TARGET_MODE"`
	Token            string `envconfig:"TOKEN" required:"true"`
	URL              string `envconfig:"URL"`
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

	graphqlEndpoint, gqlErr := url.JoinPath(client.Specification.URL, "api")
	if gqlErr != nil {
		return nil, gqlErr
	}
	client.GQLCLient = api.NewClient(graphqlEndpoint, client.Specification.DeploymentID, client.Specification.Token)

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
