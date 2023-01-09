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

// cghill@DESKTOP(xo-prod)(⎈|xo-metrics:mimir):~/src/massdriver/k8s-massdriver-metrics-stack$ aws dynamodb describe-table --table-name xo-prod-secretstore-jh82
// {
//     "Table": {
//         "AttributeDefinitions": [
//             {
//                 "AttributeName": "id",
//                 "AttributeType": "S"
//             },
//             {
//                 "AttributeName": "name",
//                 "AttributeType": "S"
//             }
//         ],
//         "TableName": "xo-prod-secretstore-jh82",
//         "KeySchema": [
//             {
//                 "AttributeName": "id",
//                 "KeyType": "HASH"
//             },
//             {
//                 "AttributeName": "name",
//                 "KeyType": "RANGE"
//             }
//         ],
//         "TableStatus": "ACTIVE",
//         "CreationDateTime": "2023-01-04T16:51:20.177000-07:00",
//         "ProvisionedThroughput": {
//             "NumberOfDecreasesToday": 0,
//             "ReadCapacityUnits": 0,
//             "WriteCapacityUnits": 0
//         },
//         "TableSizeBytes": 403,
//         "ItemCount": 1,
//         "TableArn": "arn:aws:dynamodb:us-west-2:308878630280:table/xo-prod-secretstore-jh82",
//         "TableId": "a779eaf6-cd5d-4c51-bd6b-2c24aee40d3c",
//         "BillingModeSummary": {
//             "BillingMode": "PAY_PER_REQUEST",
//             "LastUpdateToPayPerRequestDateTime": "2023-01-04T16:51:20.177000-07:00"
//         }
//     }
// }
// cghill@DESKTOP(xo-prod)(⎈|xo-metrics:mimir):~/src/massdriver/k8s-massdriver-metrics-stack$ aws dynamodb scan --table-name xo-prod-secretstore-jh82
// {
//     "Items": [
//         {
//             "value": {
//                 "S": "AQICAHg4Zrh5ktyJ7hBFn4AFnT6471MGKY9x/jiHYMtfyWNL8wH2iLzIZmRcy+xFijuacirkAAAAYzBhBgkqhkiG9w0BBwagVDBSAgEAME0GCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMjF+aTzWhUtCfFQ1vAgEQgCCqiVyx8hYBzECPnbMyBqTa4SkcOlLoGHwOxtBafGe7bw=="
//             },
//             "encryption_key_id": {
//                 "S": "arn:aws:kms:us-west-2:308878630280:key/f5ed1a46-9a73-4da4-be6f-9c43bd4327e4"
//             },
//             "id": {
//                 "S": "4b4279ee-8fff-407e-9f7d-13a1eb3d75b7.5ad835fa-7d41-49d3-bde7-e3fcf105a76a"
//             },
//             "name": {
//                 "S": "SECRET_KEY_BASE"
//             }
//         }
//     ],
//     "Count": 1,
//     "ScannedCount": 1,
//     "ConsumedCapacity": null
// }
