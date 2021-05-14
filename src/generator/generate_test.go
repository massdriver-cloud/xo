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

	os.Mkdir(bundleData.BundleDir, 0777)

	generator.Generate(bundleData)

	bundleYamlPath := fmt.Sprintf("%s/%s/bundle.yaml", bundleData.BundleDir, bundleData.Name)
	content, err := ioutil.ReadFile(bundleYamlPath)

	if err != nil {
		t.Errorf("Failed to create bundle.yaml")
	}

	if !strings.Contains(string(content), "title: aws-vpc") {
		t.Errorf("Data failed to render in the generated template")
	}

	readmePath := fmt.Sprintf("%s/%s/README.md", bundleData.BundleDir, bundleData.Name)
	content, err = ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("Failed to create Readme.md")
	}

	if !strings.Contains(string(content), "a vpc") {
		t.Errorf("Data failed to render in the generated template")
	}

	terraformPath := fmt.Sprintf("%s/%s/terraform", bundleData.BundleDir, bundleData.Name)

	mainTFPath := fmt.Sprintf("%s/main.tf", terraformPath)
	content, err = ioutil.ReadFile(mainTFPath)

	if err != nil {
		t.Errorf("Failed to create main.tf")
	}

	if !strings.Contains(string(content), "random_pet") {
		t.Errorf("Data failed to render in the generated template")
	}

	t.Cleanup(func() {
		os.RemoveAll(bundleData.BundleDir)
	})
}
