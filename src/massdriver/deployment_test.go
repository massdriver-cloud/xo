package massdriver_test

import (
	bytes "bytes"
	"io"
	"io/ioutil"
	http "net/http"
	"testing"
	"xo/src/massdriver"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	proto "google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

var testDeploymentOutput map[string]*bytes.Buffer

func outputToTestBuffer(path string) (io.Writer, error) {
	testDeploymentOutput[path] = new(bytes.Buffer)
	return testDeploymentOutput[path], nil
}

func TestStartDeployment(t *testing.T) {
	testParams, _ := structpb.NewStruct(map[string]interface{}{
		"hello": "world",
	})
	testConnections, _ := structpb.NewStruct(map[string]interface{}{
		"default": map[string]interface{}{
			"foo": "bar",
		},
	})
	testResponse := mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Bundle: &mdproto.Bundle{
				Type: "test-bundle",
			},
			Params:      testParams,
			Connections: testConnections,
		},
	}

	testDeploymentOutput = make(map[string]*bytes.Buffer)
	massdriver.OutputGenerator = outputToTestBuffer
	expectedOutput := map[string]string{
		"out/params.auto.tfvars.json": `{
  "hello": "world"
}`,
		"out/connections.auto.tfvars.json": `{
  "default": {
    "foo": "bar"
  }
}`,
		"out/bundle.txt": "test-bundle",
	}

	respBytes, _ := proto.Marshal(&testResponse)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	massdriver.MockDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	_, err := massdriver.StartDeployment("id", "token", "out")
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	for key, want := range expectedOutput {
		got, exists := testDeploymentOutput[key]
		if !exists {
			t.Fatalf("expected key %v to exist", key)
		}
		if want != got.String() {
			t.Fatalf("expected: %v, got: %v", want, got.String())
		}
	}
}
