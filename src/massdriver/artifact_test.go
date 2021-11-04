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

func TestUploadArtifactFile(t *testing.T) {
	wantId := "fakeid"
	token := "faketoken"
	wantUAR := mdproto.CompleteDeploymentRequest{
		DeploymentId:    wantId,
		DeploymentToken: token,
		Artifacts:       `[{"foo":{"bar":"baz"}},{"hello":"world"}]`,
	}

	mockDeployment := mdproto.Deployment{}

	gotUAR := new(mdproto.CompleteDeploymentRequest)
	respBytes, _ := proto.Marshal(&mockDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	massdriver.MockDoFunc = func(req *http.Request) (*http.Response, error) {
		reqBytes, _ := ioutil.ReadAll(req.Body)
		proto.Unmarshal(reqBytes, gotUAR)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	err := massdriver.UploadArtifactFile("testdata/artifacts.json", wantId, token)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	if !proto.Equal(gotUAR, &wantUAR) {
		t.Fatalf("expected: %v, got: %v", *gotUAR, wantUAR)
	}
}
