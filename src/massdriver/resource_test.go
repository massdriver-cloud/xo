package massdriver_test

import (
	bytes "bytes"
	ioutil "io/ioutil"
	http "net/http"
	"testing"
	"xo/src/massdriver"

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
	type testWant struct {
		request    *massdriver.UpdateResourceStatusRequest
		authHeader string
	}
	type test struct {
		name  string
		input testInput
		want  testWant
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
			want: testWant{
				request: &massdriver.UpdateResourceStatusRequest{
					DeploymentId:   "depId",
					ResourceId:     "resId",
					ResourceType:   "resType",
					ResourceStatus: massdriver.ResourceStatus_PROVISIONED,
				},
				authHeader: "Bearer token123",
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
			want: testWant{
				request: &massdriver.UpdateResourceStatusRequest{
					DeploymentId:   "depId",
					ResourceId:     "resId",
					ResourceType:   "resType",
					ResourceStatus: massdriver.ResourceStatus_DELETED,
				},
				authHeader: "Bearer token123",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := new(massdriver.UpdateResourceStatusRequest)
			var header *http.Header
			respBytes, _ := proto.Marshal(&massdriver.UpdateResourceStatusResponse{})
			r := ioutil.NopCloser(bytes.NewReader(respBytes))
			massdriver.MockDoFunc = func(req *http.Request) (*http.Response, error) {
				reqBytes, _ := ioutil.ReadAll(req.Body)
				header = &req.Header
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

			if !proto.Equal(got, test.want.request) {
				t.Fatalf("got: %+v, want: %+v", got, &test.want.request)
			}
			if header.Get("Authorization") != test.want.authHeader {
				t.Errorf("got %s want %s", header.Get("Authorization"), test.want.authHeader)
			}
		})
	}

}
