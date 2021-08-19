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
				"path/test-template.txt": `testing the template one and two`,
			},
		},
	}
	provisioners.OutputGenerator = outputToTestBuffer
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testAuthOutput = map[string]*bytes.Buffer{}

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
