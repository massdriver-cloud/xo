package massdriver

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type SnsInterface interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type MassdriverClient struct {
	Specification *Specification
	SNSClient     SnsInterface
}

type Specification struct {
	Action                    string `envconfig:"ACTION"`
	BundleBucket              string `envconfig:"BUNDLE_BUCKET"`
	BundleID                  string `envconfig:"BUNDLE_ID"`
	BundleName                string `envconfig:"BUNDLE_NAME"`
	BundleOwnerOrganizationID string `envconfig:"BUNDLE_OWNER_ORGANIZATION_ID"`
	BundleType                string `envconfig:"BUNDLE_TYPE"`
	DeploymentID              string `envconfig:"DEPLOYMENT_ID"`
	DynamoDBStateLockTable    string `envconfig:"DYNAMODB_STATE_LOCK_TABLE"`
	EventTopicARN             string `envconfig:"EVENT_TOPIC_ARN"`
	OrganizationID            string `envconfig:"ORGANIZATION_ID"`
	PackageID                 string `envconfig:"PACKAGE_ID"`
	PackageName               string `envconfig:"PACKAGE_NAME"`
	Provisioner               string `envconfig:"PROVISIONER"`
	S3StateBucket             string `envconfig:"S3_STATE_BUCKET"`
	S3StateRegion             string `envconfig:"S3_STATE_REGION"`
	Token                     string `envconfig:"TOKEN"`
	URL                       string `envconfig:"URL"`
}

func InitializeMassdriverClient() (*MassdriverClient, error) {
	client := new(MassdriverClient)

	var specErr error
	client.Specification, specErr = GetSpecification()
	if specErr != nil {
		return nil, specErr
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client.SNSClient = sns.NewFromConfig(cfg)

	return client, nil
}

func GetSpecification() (*Specification, error) {
	spec := Specification{}
	err := envconfig.Process("massdriver", &spec)
	return &spec, err
}

func (c MassdriverClient) PublishEventToSNS(event *Event) error {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	jsonString := string(jsonBytes)

	deduplicationId := uuid.New().String()

	input := sns.PublishInput{
		Message:                &jsonString,
		MessageDeduplicationId: &deduplicationId,
		MessageGroupId:         &c.Specification.DeploymentID,
		TopicArn:               &c.Specification.EventTopicARN,
	}

	_, err = c.SNSClient.Publish(context.Background(), &input)
	return err
}
