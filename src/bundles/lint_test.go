package bundles_test

import (
	"testing"
	"xo/src/bundles"
)

func TestTransform(t *testing.T) {
	type TestCase struct {
		Name string
		Path string
		Want bool
	}
	cases := []TestCase{
		{
			Name: "Lints a good bundle",
			Path: "./testdata/linting/good-bundle",
			Want: true,
		},
		{
			Name: "Non-existant bundle",
			Path: "./testdata/linting/lol-no-bundle",
			Want: false,
		},
		{
			Name: "Missing bundle.yaml",
			Path: "./testdata/linting/missing-bundle",
			Want: false,
		},
		{
			Name: "Invalid bundle.yaml",
			Path: "./testdata/linting/invalid-bundle",
			Want: false,
		},
		{
			Name: "Missing schema",
			Path: "./testdata/linting/missing-schema",
			Want: false,
		},
		{
			Name: "Invalid schema",
			Path: "./testdata/linting/invalid-schema",
			Want: false,
		},
		{
			Name: "Missing stories",
			Path: "./testdata/linting/missing-stories",
			Want: false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			want := test.Want
			got, _ := bundles.LintBundle(test.Path)

			if got != want {
				t.Errorf("got %t want %t", got, want)
			}
		})
	}
}
