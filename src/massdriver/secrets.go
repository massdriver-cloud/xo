package massdriver

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type DynamoDBInterface interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type KMSInterface interface {
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

type Secret struct {
	ID              string `dynamodbav:"id"`
	Name            string `dynamodbav:"name"`
	Value           string `dynamodbav:"value"`
	EncryptionKeyId string `dynamodbav:"encryption_key_id"`
}

func (c *MassdriverClient) GetSecrets(ctx context.Context) (map[string]string, error) {
	// Get default secrets
	defaultSecretsKey := fmt.Sprintf("%s.default", c.Specification.ManifestID)
	defaultSecrets, err := FetchSecretsFromDynamoDB(ctx, c.DynamoDBClient, c.Specification.SecretsTableName, defaultSecretsKey)
	if err != nil {
		return nil, err
	}

	// Get package secrets (if not a preview env)
	var packageSecrets []Secret
	if c.Specification.TargetMode != "preview" {
		packageSecretsKey := fmt.Sprintf("%s.%s", c.Specification.ManifestID, c.Specification.PackageID)
		packageSecrets, err = FetchSecretsFromDynamoDB(ctx, c.DynamoDBClient, c.Specification.SecretsTableName, packageSecretsKey)
		if err != nil {
			return nil, err
		}
	}

	combinedSecrets := append(defaultSecrets, packageSecrets...)
	result := make(map[string]string, len(combinedSecrets))

	for _, secret := range combinedSecrets {
		secretValue, err := DecryptValueWithKMS(ctx, c.KMSClient, secret.Value, secret.EncryptionKeyId)
		if err != nil {
			return result, err
		}
		result[secret.Name] = secretValue
	}

	return result, nil
}

func FetchSecretsFromDynamoDB(ctx context.Context, client DynamoDBInterface, table string, key string) ([]Secret, error) {
	input := dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("id=:id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: key},
		},
	}

	output, err := client.Query(ctx, &input)
	if err != nil {
		return nil, err
	}
	secrets := make([]Secret, output.Count)
	for idx, item := range output.Items {
		var secret Secret
		err = attributevalue.UnmarshalMap(item, &secret)
		if err != nil {
			return nil, err
		}
		secrets[idx] = secret
	}

	return secrets, nil
}

func DecryptValueWithKMS(ctx context.Context, client KMSInterface, value string, keyId string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	input := &kms.DecryptInput{
		CiphertextBlob: blob,
		KeyId:          &keyId,
	}

	result, err := client.Decrypt(ctx, input)
	if err != nil {
		return "", err
	}

	return string(result.Plaintext), nil
}
