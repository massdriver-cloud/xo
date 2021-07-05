package terraform

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
)

func GenerateAuthFiles(connectionsPath string, outputDir string) error {
	connBytes, err := ioutil.ReadFile(connectionsPath)
	if err != nil {
		return err
	}

	var awsOutput bytes.Buffer
	err = GenerateAwsAuth(connBytes, &awsOutput)
	if err != nil {
		return err
	}
	if awsOutput.Len() > 0 {
		awsCredentialsHandle, err := os.OpenFile(path.Join(outputDir, "aws.ini"), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = awsCredentialsHandle.Write(awsOutput.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}
