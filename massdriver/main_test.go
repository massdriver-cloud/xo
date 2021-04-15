package massdriver

import (
	"bytes"
	"io/ioutil"
	http "net/http"
	"testing"

	mocks "xo/utils/mocks"

	proto "github.com/golang/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetDeployment(t *testing.T) {
	testInputs, _ := structpb.NewStruct(map[string]interface{}{
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
	testDeployment := Deployment{
		Id:          "1234",
		Status:      DeploymentStatus_PENDING,
		Inputs:      testInputs,
		Connections: testConnections,
	}

	respBytes, _ := proto.Marshal(&testDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	mocks.MockDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	got, _ := GetDeployment("id", "token")

	gotString := proto.MarshalTextString(got)
	wantString := proto.MarshalTextString(&testDeployment)

	if gotString != wantString {
		t.Fatalf("expected: %+v, got: %+v", gotString, wantString)
	}
}

func TestWriteSchema(t *testing.T) {
	testInputs, _ := structpb.NewStruct(map[string]interface{}{
		"aws_region": "us-east-1",
		"some_key":   true,
		"other_key":  27,
		"nested_key": map[string]interface{}{
			"key_a": "value_a",
			"key_b": 123.456,
		},
	})

	buf := bytes.Buffer{}
	writeSchema(testInputs, &buf)

	gotString := buf.String()
	wantString := "{\"aws_region\":\"us-east-1\",\"nested_key\":{\"key_a\":\"value_a\",\"key_b\":123.456},\"other_key\":27,\"some_key\":true}"

	if gotString != wantString {
		t.Fatalf("expected: %v, got: %v", gotString, wantString)
	}
}
