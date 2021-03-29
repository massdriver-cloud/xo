package massdriver

import (
	"testing"
  "bytes"
	"io/ioutil"
  http "net/http"

  mocks "xo/utils/mocks"

  proto "github.com/golang/protobuf/proto"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetDeployment(t *testing.T) {
  expectedToken := "testTokenabcd1234"
  testDeployment := Deployment{
    Id:     "1234",
    Token:  expectedToken,
  }

  respBytes, _ := proto.Marshal(&testDeployment)
  r := ioutil.NopCloser(bytes.NewReader(respBytes))
	mocks.MockDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	got, _ := GetDeployment("whatever")

	if got != expectedToken {
		t.Fatalf("expected: %v, got: %v", expectedToken, got)
	}
}
