package transform

import "testing"

func TestTransform(t *testing.T) {
	type TestCase struct {
		Name        string
		Input       string
		Transformer string
		Want        string
	}
	cases := []TestCase{
		{
			Name:        "Transforms json",
			Input:       `{"greeting":"hello", "target":"jimmy", "punctuation":"!"}`,
			Transformer: "./testdata/hello.js",
			Want:        `{"statement":"hello jimmy!"}`,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			want := test.Want
			got, _ := Transform(test.Input, test.Transformer)

			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
	}
}
