package massdriver

import (
	bytes "bytes"
	ioutil "io/ioutil"
	http "net/http"
	"testing"
	mocks "xo/src/utils/mocks"

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
	artifact1 := Artifact{
		Metadata: &ArtifactMetadata{
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
	artifact2 := Artifact{
		Metadata: &ArtifactMetadata{
			ProviderResourceId: "A9DA1D78-4B93-420D-9B1B-289A164A7400",
			Type:               "aws-ec2-subnet",
			Name:               "Your new Subnet",
		},
		Data:  data2,
		Specs: specs2,
	}
	wantArtifacts := []*Artifact{
		&artifact1,
		&artifact2,
	}

	mockDeployment := Deployment{}

	uar := new(UploadArtifactsRequest)
	var header *http.Header
	respBytes, _ := proto.Marshal(&mockDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	mocks.MockDoFunc = func(req *http.Request) (*http.Response, error) {
		reqBytes, _ := ioutil.ReadAll(req.Body)
		header = &req.Header
		proto.Unmarshal(reqBytes, uar)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	wantId := "fakeid"
	token := "faketoken"
	wantHeader := "Bearer " + token
	UploadArtifactFile("testdata/artifacts.json", wantId, token)

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
	if header.Get("Authorization") != wantHeader {
		t.Errorf("got %s want %s", header.Get("Authorization"), wantHeader)
	}
}
