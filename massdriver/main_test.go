package massdriver

import (
	"bytes"
	"io/ioutil"
	http "net/http"
	"testing"

	mocks "xo/utils/mocks"

	proto "github.com/golang/protobuf/proto"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetDeployment(t *testing.T) {
	expectedToken := "nothing"
	testDeployment := Deployment{
		Id: "1234",
	}

	respBytes, _ := proto.Marshal(&testDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	mocks.MockDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	got, _ := GetDeployment("whatever", "token")

	if got != expectedToken {
		t.Fatalf("expected: %v, got: %v", expectedToken, got)
	}
}
