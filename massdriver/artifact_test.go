package massdriver

import (
	bytes "bytes"
	ioutil "io/ioutil"
	http "net/http"
	"testing"
	mocks "xo/utils/mocks"

	proto "google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestUploadArtifactFile(t *testing.T) {
	secrets1, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-abc",
	})
	specs1, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact1 := Artifact{
		Id:      "A6728066-189C-4405-9020-CCB168F28E7D",
		Type:    "aws-ec2-vpc",
		Name:    "Your new VPC",
		Secrets: secrets1,
		Specs:   specs1,
	}
	secrets2, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-xyz",
	})
	specs2, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact2 := Artifact{
		Id:      "A9DA1D78-4B93-420D-9B1B-289A164A7400",
		Type:    "aws-ec2-subnet",
		Name:    "Your new Subnet",
		Secrets: secrets2,
		Specs:   specs2,
	}
	wantArtifacts := []*Artifact{
		&artifact1,
		&artifact2,
	}

	mockDeployment := Deployment{}

	uar := new(UploadArtifactsRequest)
	respBytes, _ := proto.Marshal(&mockDeployment)
	r := ioutil.NopCloser(bytes.NewReader(respBytes))
	mocks.MockDoFunc = func(req *http.Request) (*http.Response, error) {
		reqBytes, _ := ioutil.ReadAll(req.Body)
		proto.Unmarshal(reqBytes, uar)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	wantId := "fakeid"
	wantToken := "faketoken"
	UploadArtifactFile("testdata/artifacts.json", wantId, wantToken)

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
	if uar.GetToken() != wantToken {
		t.Errorf("got %s want %s", uar.GetToken(), wantToken)
	}
}

func TestCreateArtifactsFromJsonBytes(t *testing.T) {
	jsonBlob := []byte(`
		[{
			"id": "A6728066-189C-4405-9020-CCB168F28E7D",
			"type": "aws-ec2-vpc",
			"name": "Your new VPC",
			"secrets": {
				"id": "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-abc"
			},
			"specs": {
				"aws": {
					"region": "us-east1"
				}
			}
		},
		{
			"id": "A9DA1D78-4B93-420D-9B1B-289A164A7400",
			"type": "aws-ec2-subnet",
			"name": "Your new Subnet",
			"secrets": {
				"id": "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-xyz"
			},
			"specs": {
				"aws": {
					"region": "us-east1"
				}
			}
		}]
    `)

	got, _ := createArtifactsFromJsonBytes(jsonBlob)

	secrets1, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-abc",
	})
	specs1, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact1 := Artifact{
		Id:      "A6728066-189C-4405-9020-CCB168F28E7D",
		Type:    "aws-ec2-vpc",
		Name:    "Your new VPC",
		Secrets: secrets1,
		Specs:   specs1,
	}
	secrets2, _ := structpb.NewStruct(map[string]interface{}{
		"id": "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-xyz",
	})
	specs2, _ := structpb.NewStruct(map[string]interface{}{
		"aws": map[string]interface{}{
			"region": "us-east1",
		},
	})
	artifact2 := Artifact{
		Id:      "A9DA1D78-4B93-420D-9B1B-289A164A7400",
		Type:    "aws-ec2-subnet",
		Name:    "Your new Subnet",
		Secrets: secrets2,
		Specs:   specs2,
	}
	want := []*Artifact{
		&artifact1,
		&artifact2,
	}

	if len(got) != len(want) {
		t.Fatalf("Array lengths different. expected: %v, got: %v", len(got), len(want))
	}
	for ii := 0; ii < len(got); ii++ {
		if !proto.Equal(got[ii], want[ii]) {
			t.Fatalf("expected: %+v, got: %+v", got[ii].String(), want[ii].String())
		}
	}
}
