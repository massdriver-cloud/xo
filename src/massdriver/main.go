package massdriver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"xo/src/api"

	"github.com/Khan/genqlient/graphql"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

var MassdriverURL = "https://api.massdriver.cloud/"

// EventPublisher will know how to publish an event to a specific target (sns, logs etc.)
type EventPublisher interface {
	Publish(ctx context.Context, event *Event) error
}

// SnsInterface allows for mocking the sns client in the tests without needing aws
type SnsInterface interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type MassdriverClient struct {
	GQLCLient     graphql.Client
	Specification *Specification
	Publisher     EventPublisher
}

type Specification struct {
	Action         string `envconfig:"ACTION"`
	BundleID       string `envconfig:"BUNDLE_ID" required:"true"`
	BundleName     string `envconfig:"BUNDLE_NAME"`
	BundleType     string `envconfig:"BUNDLE_TYPE"`
	DeploymentID   string `envconfig:"DEPLOYMENT_ID" required:"true"`
	EventTopicARN  string `envconfig:"EVENT_TOPIC_ARN" required:"true"`
	ManifestID     string `envconfig:"MANIFEST_ID"`
	OrganizationID string `envconfig:"ORGANIZATION_ID" required:"true"`
	PackageID      string `envconfig:"PACKAGE_ID" required:"true"`
	PackageName    string `envconfig:"PACKAGE_NAME" required:"true"`
	TargetMode     string `envconfig:"TARGET_MODE"`
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

	graphqlEndpoint, gqlErr := url.JoinPath(client.Specification.URL, "api")
	if gqlErr != nil {
		return nil, gqlErr
	}
	client.GQLCLient = api.NewClient(graphqlEndpoint, client.Specification.DeploymentID, client.Specification.Token)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// If the ARN doesn't exist, assume we are running locally
	if os.Getenv("MASSDRIVER_EVENT_TOPIC_ARN") == "" {
		client.Publisher = &localPublisher{}
	} else {
		client.Publisher = &SNSPublisher{
			SNSClient:     sns.NewFromConfig(cfg),
			Specification: client.Specification,
		}
	}

	return client, nil
}

func GetSpecification() (*Specification, error) {
	// If the ARN doesn't exist, assume we are running locally
	if os.Getenv("MASSDRIVER_EVENT_TOPIC_ARN") == "" {
		return &Specification{}, nil
	}
	spec := Specification{}
	err := envconfig.Process("massdriver", &spec)
	return &spec, err
}

func (c MassdriverClient) PublishEvent(event *Event) error {
	return c.Publisher.Publish(context.Background(), event)
}

type SNSPublisher struct {
	Specification *Specification
	SNSClient     SnsInterface
}

func (l *SNSPublisher) Publish(ctx context.Context, event *Event) error {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	jsonString := string(jsonBytes)

	deduplicationId := uuid.New().String()

	input := sns.PublishInput{
		Message:                &jsonString,
		MessageDeduplicationId: &deduplicationId,
		MessageGroupId:         &l.Specification.DeploymentID,
		TopicArn:               &l.Specification.EventTopicARN,
	}

	_, err = l.SNSClient.Publish(context.Background(), &input)
	return err
}

type localPublisher struct{}

func (l *localPublisher) Publish(ctx context.Context, event *Event) error {
	out, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
