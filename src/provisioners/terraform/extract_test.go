package terraform_test

import (
	"os"
	"testing"
	"xo/src/provisioners/terraform"
)

func TestExtract(t *testing.T) {
	type testData struct {
		name     string
		input    string
		expected string
	}
	tests := []testData{
		{
			name:     "empty schema",
			input:    "testdata/terraform-output.json",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			err = terraform.Extract(input)
			if err != nil {
				t.Fatalf("%d, unexpected error", err)
			}
			want := tc.expected
			got := ""

			if got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}
