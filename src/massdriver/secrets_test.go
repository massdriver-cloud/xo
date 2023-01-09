package massdriver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"xo/src/massdriver"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type mockDynamoDBClient struct {
	secrets map[string][]massdriver.Secret
}
type mockKMSClient struct{}

func (m *mockDynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	key := params.ExpressionAttributeValues[":id"].(*types.AttributeValueMemberS).Value
	items := make([]map[string]types.AttributeValue, len(m.secrets[key]))
	for idx := range m.secrets[key] {
		item, err := attributevalue.MarshalMap(m.secrets[key][idx])
		if err != nil {
			return nil, err
		}
		items[idx] = item
	}
	output := dynamodb.QueryOutput{
		Count: int32(len(m.secrets[key])),
		Items: items,
	}
	return &output, nil
}

func (m *mockKMSClient) Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	output := kms.DecryptOutput{
		Plaintext: []byte(fmt.Sprintf("%s+%s", *params.KeyId, strings.ToUpper(string(params.CiphertextBlob)))),
	}
	return &output, nil
}

func TestGetSecrets(t *testing.T) {
	type testData struct {
		name   string
		table  string
		client massdriver.MassdriverClient
		want   map[string]string
	}
	tests := []testData{
		{
			name:  "default and preview",
			table: "table",
			client: massdriver.MassdriverClient{
				Specification: &massdriver.Specification{
					ManifestID: "manifestId",
					PackageID:  "packageId",
				},
				DynamoDBClient: &mockDynamoDBClient{
					secrets: map[string][]massdriver.Secret{
						"manifestId.default": {{
							ID:              "simple",
							Name:            "name",
							Value:           "czNjcmV0",
							EncryptionKeyId: "key-id",
						}},
						"manifestId.packageId": {{
							ID:              "simple",
							Name:            "package",
							Value:           "YW4wdGhlcnMzY3JldA==",
							EncryptionKeyId: "key-id",
						}},
					},
				},
				KMSClient: &mockKMSClient{},
			},
			want: map[string]string{
				"name":    "key-id+S3CRET",
				"package": "key-id+AN0THERS3CRET",
			},
		},
		{
			name:  "preview envrionment (no package secrets)",
			table: "table",
			client: massdriver.MassdriverClient{
				Specification: &massdriver.Specification{
					ManifestID: "manifestId",
					PackageID:  "packageId",
					TargetMode: "preview",
				},
				DynamoDBClient: &mockDynamoDBClient{
					secrets: map[string][]massdriver.Secret{
						"manifestId.default": {{
							ID:              "simple",
							Name:            "name",
							Value:           "czNjcmV0",
							EncryptionKeyId: "key-id",
						}},
						"manifestId.packageId": {{
							ID:              "simple",
							Name:            "package",
							Value:           "YW4wdGhlcnMzY3JldA==",
							EncryptionKeyId: "key-id",
						}},
					},
				},
				KMSClient: &mockKMSClient{},
			},
			want: map[string]string{
				"name": "key-id+S3CRET",
			},
		},
		{
			name:  "collision preview wins",
			table: "table",
			client: massdriver.MassdriverClient{
				Specification: &massdriver.Specification{
					ManifestID: "manifestId",
					PackageID:  "packageId",
				},
				DynamoDBClient: &mockDynamoDBClient{
					secrets: map[string][]massdriver.Secret{
						"manifestId.default": {{
							ID:              "simple",
							Name:            "name",
							Value:           "czNjcmV0",
							EncryptionKeyId: "key-id",
						}},
						"manifestId.packageId": {{
							ID:              "simple",
							Name:            "name",
							Value:           "YW4wdGhlcnMzY3JldA==",
							EncryptionKeyId: "key-id",
						}},
					},
				},
				KMSClient: &mockKMSClient{},
			},
			want: map[string]string{
				"name": "key-id+AN0THERS3CRET",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.client.GetSecrets(context.Background())
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("want: %v, got: %v", len(tc.want), len(got))
			}

			for k, v := range tc.want {
				if got[k] != v {
					t.Fatalf("want: %v, got: %v", v, got[k])
				}
			}
		})
	}
}

func TestFetchSecretsFromDynamoDB(t *testing.T) {
	type testData struct {
		name    string
		table   string
		keyId   string
		secrets map[string][]massdriver.Secret
	}
	tests := []testData{
		{
			name:  "foo",
			table: "table",
			keyId: "key",
			secrets: map[string][]massdriver.Secret{
				"key": {{
					ID:              "simple",
					Name:            "name",
					Value:           "encrypted-value",
					EncryptionKeyId: "key-id",
				}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDynamoDBClient := mockDynamoDBClient{
				secrets: tc.secrets,
			}
			got, err := massdriver.FetchSecretsFromDynamoDB(context.Background(), &mockDynamoDBClient, tc.table, tc.keyId)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(tc.secrets) != len(got) {
				t.Fatalf("expected: %v, got: %v", len(tc.secrets), len(got))
			}

			for idx := range tc.secrets[tc.keyId] {
				if got[idx] != tc.secrets[tc.keyId][idx] {
					t.Fatalf("want: %v, got: %v", got[idx], tc.secrets[tc.keyId][idx])
				}
			}
		})
	}
}

func TestDecryptValue(t *testing.T) {
	type testData struct {
		name           string
		keyId          string
		encryptedValue string
		want           string
	}
	tests := []testData{
		{
			name:           "simple",
			keyId:          "key",
			encryptedValue: "YWJjZA==",
			want:           "key+ABCD",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := mockKMSClient{}
			got, err := massdriver.DecryptValueWithKMS(context.Background(), &mockClient, tc.encryptedValue, tc.keyId)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if got != tc.want {
				t.Fatalf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}
