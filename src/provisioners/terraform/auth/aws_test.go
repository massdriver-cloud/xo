package terraform

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestGenerateAwsAuth(t *testing.T) {
	type test struct {
		name string
		path string
		want string
	}

	tests := []test{
		{
			name: "empty auth",
			path: "testdata/credentials-no-auth.json",
			want: "",
		},
		{
			name: "single aws auth",
			path: "testdata/credentials-aws-auth.json",
			want: `[default]
aws_access_key_id=FAKEFAKEFAKEFAKE
aws_secret_access_key=FAKEfakeFAKEfakeFAKEfake

`,
		},
		{
			name: "multiple aws auths",
			path: "testdata/credentials-aws-multi.json",
			want: `[default]
aws_access_key_id=FAKEFAKEFAKEFAKE
aws_secret_access_key=FAKEfakeFAKEfakeFAKEfake

[another]
aws_access_key_id=ANOTHERANOTHER
aws_secret_access_key=ANOTHERfakeANOTHERfake

`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input, _ := ioutil.ReadFile(tc.path)
			output := bytes.Buffer{}

			GenerateAwsAuth(input, &output)

			got := output.String()

			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}
