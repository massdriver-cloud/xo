package provisioners_test

import (
	"bytes"
	"io"
	"testing"
	"xo/src/provisioners"
)

var testAuthOutput map[string]*bytes.Buffer

func outputToTestBuffer(dir, name, ext string) (io.Writer, error) {
	key := dir + "/" + name + "." + ext
	testAuthOutput[key] = new(bytes.Buffer)
	return testAuthOutput[key], nil
}

func TestGenerateAuthFiles(t *testing.T) {
	type testData struct {
		name       string
		schemaPath string
		dataPath   string
		expected   map[string]string
	}
	tests := []testData{
		{
			name:       "Test all renders",
			schemaPath: "testdata/schema-all.json",
			dataPath:   "testdata/data-all.json",
			expected: map[string]string{
				"path/test-json.json": `{"foo":"bar","hello":"world"}`,
				"path/test-yaml.yaml": `test: yaml
`,
				"path/test-ini.ini": `lol=rofl

`,
				"path/test-template.json": `{"new1":"one","new2":"two"}`,
			},
		},
		{
			name:       "Test real auth files",
			schemaPath: "testdata/schema-real.json",
			dataPath:   "testdata/data-real.json",
			expected: map[string]string{
				"path/aws-creds.ini": `[default]
aws_secret_access_key=lolroflnopasswordherefbi
aws_access_key_id=FAKEFAKEFAKEFAKE

`,
				"path/aws-role.ini": `[default]
source_profile=eks
role_arn=arn:aws:iam::123456789012:role/testrole

[eks]

`,
				"path/k8s-authentication.yaml": `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: notarealcertificateherejimmy
    server: https://my.dumb.k8s
  name: default
contexts:
- context:
    cluster: default
    user: default
  name: default
current-context: default
kind: Config
users:
- name: default
  user:
    token: goaheadandtrythistokenandseeifitworks
`,
			},
		},
	}
	provisioners.OutputGenerator = outputToTestBuffer
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testAuthOutput = make(map[string]*bytes.Buffer)

			err := provisioners.GenerateAuthFiles(tc.schemaPath, tc.dataPath, "path")
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}

			if len(tc.expected) != len(testAuthOutput) {
				t.Fatalf("expected: %v, got: %v", len(tc.expected), len(testAuthOutput))
			}

			for key, want := range tc.expected {
				got, exists := testAuthOutput[key]
				if !exists {
					t.Fatalf("expected key %v to exist", key)
				}
				if want != got.String() {
					t.Fatalf("expected: %v, got: %v", want, got.String())
				}
			}
		})
	}
}
