package massdriver_test

import (
	bytes "bytes"
	ioutil "io/ioutil"
	http "net/http"
	"testing"
	"xo/src/massdriver"

	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
	proto "google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestUploadArtifactFile(t *testing.T) {
	data1, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-abc",
	})
	specs1, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact1 := mdproto.Artifact{
		Metadata: &mdproto.ArtifactMetadata{
			ProviderResourceId: "A6728066-189C-4405-9020-CCB168F28E7D",
			Type:               "aws-ec2-vpc",
			Name:               "Your new VPC",
		},
		Data:  data1,
		Specs: specs1,
	}
	data2, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-xyz",
	})
	specs2, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact2 := mdproto.Artifact{
		Metadata: &mdproto.ArtifactMetadata{
			ProviderResourceId: "A9DA1D78-4B93-420D-9B1B-289A164A7400",
			Type:               "aws-ec2-subnet",
			Name:               "Your new Subnet",
		},
		Data:  data2,
		Specs: specs2,
	}
	wantArtifacts := []*mdproto.Artifact{
		&artifact1,
		&artifact2,
	}

	mockDeployment := mdproto.Deployment{}

	uar := new(mdproto.UploadArtifactsRequest)
	respBytes, _ := proto.Marshal(&mockDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	massdriver.MockDoFunc = func(req *http.Request) (*http.Response, error) {
		reqBytes, _ := ioutil.ReadAll(req.Body)
		proto.Unmarshal(reqBytes, uar)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	wantId := "fakeid"
	token := "faketoken"
	err := massdriver.UploadArtifactFile("testdata/artifacts.json", wantId, token)
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}

	gotArtifacts := uar.GetArtifacts()
	if len(gotArtifacts) != len(wantArtifacts) {
		t.Fatalf("Array lengths different. expected: %v, got: %v", len(gotArtifacts), len(wantArtifacts))
	}
	for ii := 0; ii < len(gotArtifacts); ii++ {
		if !proto.Equal(gotArtifacts[ii], wantArtifacts[ii]) {
			t.Fatalf("expected: %+v, got: %+v", gotArtifacts[ii].String(), wantArtifacts[ii].String())
		}
	}
	if uar.GetDeploymentId() != wantId {
		t.Errorf("got %s want %s", uar.GetDeploymentId(), wantId)
	}
}
