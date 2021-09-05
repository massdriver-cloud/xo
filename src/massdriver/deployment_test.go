package massdriver

import (
	bytes "bytes"
	json "encoding/json"
	"io/ioutil"
	http "net/http"
	"testing"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	proto "google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestStartDeployment(t *testing.T) {
	testParams, _ := structpb.NewStruct(map[string]interface{}{
		"aws_region": "us-east-1",
		"some_key":   true,
		"other_key":  27,
		"nested_key": map[string]interface{}{
			"key_a": "value_a",
			"key_b": 123.456,
		},
	})
	testConnections, _ := structpb.NewStruct(map[string]interface{}{
		"default": map[string]interface{}{
			"aws_access_key_id":     "ACOVIBUOISKLWJEFKJL",
			"aws_secret_access_key": "8ba0u90uwe9fuq90j3490tj0q923u12093u09gj90u130",
		},
	})
	testResponse := mdproto.StartDeploymentResponse{
		Deployment: &mdproto.Deployment{
			Id:          "1234",
			Status:      mdproto.DeploymentStatus_PENDING,
			Params:      testParams,
			Connections: testConnections,
		},
	}

	respBytes, _ := proto.Marshal(&testResponse)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	MockDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	got, err := StartDeployment("id", "token")
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	if !proto.Equal(got, testResponse.Deployment) {
		t.Fatalf("expected: %+v, got: %+v", got.String(), testResponse.Deployment.String())
	}
}

func TestWriteSchema(t *testing.T) {
	testParams, _ := structpb.NewStruct(map[string]interface{}{
		"aws_region": "us-east-1",
		"some_key":   true,
		"other_key":  27,
		"nested_key": map[string]interface{}{
			"key_a": "value_a",
			"key_b": 123.456,
		},
	})

	buf := bytes.Buffer{}
	writeSchema(testParams, &buf)

	wantString := `{"aws_region":"us-east-1","nested_key":{"key_a":"value_a","key_b":123.456},"other_key":27,"some_key":true}`

	gotBytes := new(bytes.Buffer)
	wantBytes := new(bytes.Buffer)
	json.Compact(gotBytes, buf.Bytes())
	json.Compact(wantBytes, []byte(wantString))

	if gotBytes.String() != wantBytes.String() {
		t.Fatalf("want: %v, got: %v", wantBytes, gotBytes)
	}
}
