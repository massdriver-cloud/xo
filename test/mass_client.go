package testmass

import (
	"context"
	"xo/src/massdriver"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type testClient struct {
	MassClient massdriver.MassdriverClient
}

func NewMassdriverTestClient(deploymentId string) *testClient {
	return &testClient{MassClient: massdriver.MassdriverClient{
		Specification: &massdriver.Specification{
			DeploymentID: deploymentId,
		},
		Publisher: &massdriver.SNSPublisher{
			Specification: &massdriver.Specification{
				DeploymentID: deploymentId,
			},
			SNSClient: &SNSTestClient{},
		},
	},
	}
}

func (tc *testClient) GetSNS() *SNSTestClient {
	// Not a fan of this but it gets the job done
	client := tc.MassClient.Publisher.(*massdriver.SNSPublisher)
	sns := client.SNSClient.(*SNSTestClient)
	return sns
}

func (tc *testClient) GetSNSMessages() []string {
	return tc.GetSNS().messages
}

type SNSTestClient struct {
	Input    *sns.PublishInput
	messages []string
}

func (c *SNSTestClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	c.Input = params
	c.messages = append(c.messages, *params.Message)
	return &sns.PublishOutput{}, nil
}
