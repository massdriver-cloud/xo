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

func TestUpdateResource(t *testing.T) {
	type testInput struct {
		deploymentId   string
		token          string
		resourceId     string
		resourceType   string
		resourceStatus string
	}

	type test struct {
		name  string
		input testInput
		want  mdproto.UpdateResourceStatusRequest
	}
	tests := []test{
		{
			name: "provisioned",
			input: testInput{
				deploymentId:   "depId",
				token:          "token123",
				resourceId:     "resId",
				resourceType:   "resType",
				resourceStatus: "provisioned",
			},
			want: mdproto.UpdateResourceStatusRequest{
				DeploymentId:    "depId",
				DeploymentToken: "token123",
				ResourceId:      "resId",
				ResourceType:    "resType",
				ResourceStatus:  mdproto.ResourceStatus_PROVISIONED,
			},
		},
		{
			name: "deleted",
			input: testInput{
				deploymentId:   "depId",
				token:          "token123",
				resourceId:     "resId",
				resourceType:   "resType",
				resourceStatus: "deleted",
			},
			want: mdproto.UpdateResourceStatusRequest{
				DeploymentId:    "depId",
				DeploymentToken: "token123",
				ResourceId:      "resId",
				ResourceType:    "resType",
				ResourceStatus:  mdproto.ResourceStatus_DELETED,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := new(mdproto.UpdateResourceStatusRequest)
			respBytes, _ := proto.Marshal(&mdproto.UpdateResourceStatusResponse{})
			r := ioutil.NopCloser(bytes.NewReader(respBytes))
			massdriver.MockDoFunc = func(req *http.Request) (*http.Response, error) {
				reqBytes, _ := ioutil.ReadAll(req.Body)
				proto.Unmarshal(reqBytes, got)
				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			}

			massdriver.UpdateResource(
				test.input.deploymentId,
				test.input.token,
				test.input.resourceId,
				test.input.resourceType,
				test.input.resourceStatus,
			)

			if !proto.Equal(got, &test.want) {
				t.Fatalf("got: %+v, want: %+v", got, test.want)
			}
		})
	}

}
