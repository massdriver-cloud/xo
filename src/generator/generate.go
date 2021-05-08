package generator

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/google/uuid"
)

type TemplateData struct {
	Name        string
	Description string
	Provisioner string
	TemplateDir string
	BundleDir   string
}

func (g TemplateData) Uuid() string {
	uuid, err := uuid.NewUUID()

	if err != nil {
		panic(nil)
	}

	return uuid.String()
}

func Generate(data TemplateData) error {
	files, err := getTemplateFiles(data.TemplateDir)

	if err != nil {
		return err
	}

	bundleDir := fmt.Sprintf("%s/%s", data.BundleDir, data.Name)
	terraformDir := fmt.Sprintf("%s/%s", bundleDir, "terraform")
	os.Mkdir(bundleDir, 0777)
	os.Mkdir(terraformDir, 0777)

	for _, file := range files {
		templatePath := data.TemplateDir + "/" + file.Name()
		renderPath := fmt.Sprintf("%s/%s", bundleDir, file.Name())
		tmpl, err := template.ParseFiles(templatePath)

		if err != nil {
			return err
		}

		fileToWrite, err := os.Create(renderPath)

		if err != nil {
			return err
		}

		tmpl.Execute(fileToWrite, data)

		fileToWrite.Close()
	}

	return nil
}

func getTemplateFiles(templateDir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(templateDir)

	return files, err
}
