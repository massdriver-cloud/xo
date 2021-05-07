package generator

import (
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

	for _, file := range files {
		templatePath := data.TemplateDir + "/" + file.Name()
		filePath := data.BundleDir + "/" + file.Name()
		tmpl, err := template.ParseFiles(templatePath)

		if err != nil {
			return err
		}

		fileToWrite, err := os.Create(filePath)

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
