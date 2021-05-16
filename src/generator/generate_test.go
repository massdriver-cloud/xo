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

	assertFileCreatedAndContainsText := func(t testing.TB, filename, expectedContent string) {
		t.Helper()
		content, err := ioutil.ReadFile(filename)

		if err != nil {
			t.Errorf("Failed to create %s", filename)
		}

		if !strings.Contains(string(content), expectedContent) {
			t.Errorf("Data failed to render in template %s", filename)
		}
	}

	os.Mkdir(bundleData.BundleDir, 0777)

	generator.Generate(bundleData)

	bundleYamlPath := fmt.Sprintf("%s/%s/bundle.yaml", bundleData.BundleDir, bundleData.Name)
	expectedContent := "title: aws-vpc"

	assertFileCreatedAndContainsText(t, bundleYamlPath, expectedContent)

	readmePath := fmt.Sprintf("%s/%s/README.md", bundleData.BundleDir, bundleData.Name)
	expectedContent = "a vpc"

	assertFileCreatedAndContainsText(t, readmePath, expectedContent)

	terraformPath := fmt.Sprintf("%s/%s/terraform", bundleData.BundleDir, bundleData.Name)
	mainTFPath := fmt.Sprintf("%s/main.tf", terraformPath)
	expectedContent = "random_pet"

	assertFileCreatedAndContainsText(t, mainTFPath, expectedContent)

	t.Cleanup(func() {
		os.RemoveAll(bundleData.BundleDir)
	})
}
