package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func Generate(data TemplateData) {
	bundleDir := fmt.Sprintf("%s/%s", data.BundleDir, data.Name)
	currentDirectory := ""

	filepath.WalkDir(data.TemplateDir, func(path string, info fs.DirEntry, err error) error {

		if info.IsDir() {
			if isRootPath(path, data.TemplateDir) {
				os.MkdirAll(bundleDir, 0777)
				return nil
			}

			subDirectory := fmt.Sprintf("%s/%s", bundleDir, info.Name())
			os.Mkdir(subDirectory, 0777)
			currentDirectory = fmt.Sprintf("%s/", info.Name())
			return nil
		}

		renderPath := fmt.Sprintf("%s/%s%s", bundleDir, currentDirectory, info.Name())
		renderTemplate(path, renderPath, data)

		return nil
	})
}

func renderTemplate(path, renderPath string, data TemplateData) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	fileToWrite, err := os.Create(renderPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	tmpl.Execute(fileToWrite, data)

	fileToWrite.Close()
}

func isRootPath(rootPath, currentPath string) bool {
	return rootPath == currentPath
}
