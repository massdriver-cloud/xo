package massdriver_test

import (
	bytes "bytes"
	ioutil "io/ioutil"
	http "net/http"
	"testing"
	"xo/src/massdriver"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"

	proto "google.golang.org/protobuf/proto"
)

func TestSendResourceProgress(t *testing.T) {
	depId := "fakeid"
	depToken := "faketoken"
	want := mdproto.ProvisionerProgressUpdateRequest{
		DeploymentId:    depId,
		DeploymentToken: depToken,
	}

	mockResponse := mdproto.ProvisionerProgressUpdateResponse{}
	got := new(mdproto.ProvisionerProgressUpdateRequest)
	respBytes, _ := proto.Marshal(&mockResponse)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	massdriver.MockDoFunc = func(req *http.Request) (*http.Response, error) {
		reqBytes, _ := ioutil.ReadAll(req.Body)
		proto.Unmarshal(reqBytes, got)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	err := massdriver.SendProvisionerProgressUpdate(&want)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	if !proto.Equal(got, &want) {
		t.Fatalf("expected: %v, got: %v", *got, want)
	}
}
