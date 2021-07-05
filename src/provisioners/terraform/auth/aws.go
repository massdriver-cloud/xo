package terraform

import (
	"encoding/json"
	"io"

	"gopkg.in/ini.v1"
)

type AwsCredentials struct {
	AwsAccessKeyId     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
}

func GenerateAwsAuth(connections []byte, output io.Writer) error {
	var objMap map[string]json.RawMessage
	err := json.Unmarshal(connections, &objMap)
	if err != nil {
		return err
	}

	cfg := ini.Empty()
	ini.PrettyFormat = false

	for jsonKey, jsonValueBytes := range objMap {
		var awsCredentials AwsCredentials
		err = json.Unmarshal(jsonValueBytes, &awsCredentials)
		if err == nil {
			cfg.NewSection(jsonKey)
			cfg.Section(jsonKey).Key("aws_access_key_id").SetValue(awsCredentials.AwsAccessKeyId)
			cfg.Section(jsonKey).Key("aws_secret_access_key").SetValue(awsCredentials.AwsSecretAccessKey)
		}
	}

	_, err = cfg.WriteTo(output)

	return err
}
