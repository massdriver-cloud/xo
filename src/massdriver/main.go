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
	Specification Specification
	SNSClient     SnsInterface
}

type Specification struct {
	DeploymentID  string `envconfig:"DEPLOYMENT_ID"`
	EventTopicARN string `envconfig:"EVENT_TOPIC_ARN"`
	Provisioner   string `envconfig:"PROVISIONER"`
}

func InitializeMassdriverClient() (*MassdriverClient, error) {
	client := new(MassdriverClient)
	err := envconfig.Process("massdriver", &client.Specification)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client.SNSClient = sns.NewFromConfig(cfg)

	return client, nil
}

func (c MassdriverClient) PublishEventToSNS(event *Event) error {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	jsonString := string(jsonBytes)

	deduplicationId := uuid.New().String()
	groupId := event.Payload.GetDeploymentId()

	input := sns.PublishInput{
		Message:                &jsonString,
		MessageDeduplicationId: &deduplicationId,
		MessageGroupId:         &groupId,
		TopicArn:               &c.Specification.EventTopicARN,
	}

	_, err = c.SNSClient.Publish(context.Background(), &input)
	return err
}
