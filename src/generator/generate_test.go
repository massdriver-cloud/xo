package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"xo/src/generator"
)

func TestGenerate(t *testing.T) {
	//TODO: We should be mocking the filesystem here.
	//The testing/testFS package isn't quite there yet and afero although cool seems like it has implications
	//for the broader application.
	bundleData := generator.TemplateData{
		Name:        "aws-vpc",
		Description: "a vpc",
		TemplateDir: "./testdata/templates",
		BundleDir:   "./testdata/bundle",
		Provisioner: "terraform",
	}

	err := generator.Generate(bundleData)

	if err != nil {
		t.Errorf("Generating the bundle encounted an unexpected error %v", err)
	}

	bundleYamlPath := fmt.Sprintf("%s/bundle.yaml", bundleData.BundleDir)
	content, err := ioutil.ReadFile(bundleYamlPath)

	if err != nil {
		t.Errorf("Failed to create bundle.yaml")
	}

	if !strings.Contains(string(content), "title: aws-vpc") {
		t.Errorf("Data failed to render in the generated template")
	}

	readmePath := fmt.Sprintf("%s/README.md", bundleData.BundleDir)
	content, err = ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("Failed to create Readme.md")
	}

	if !strings.Contains(string(content), "a vpc") {
		t.Errorf("Data failed to render in the generated template")
	}

	t.Cleanup(func() {
		os.RemoveAll(bundleYamlPath)
		os.RemoveAll(readmePath)
	})
}
