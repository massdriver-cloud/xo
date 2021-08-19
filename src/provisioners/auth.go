package provisioners

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"xo/src/jsonschema"

	"gopkg.in/yaml.v2"
)

var OutputGenerator func(string, string, string) (io.Writer, error)

func init() {
	OutputGenerator = outputToFile
}

func outputToFile(dir, name, ext string) (io.Writer, error) {
	return os.OpenFile(path.Join(dir, name+"."+ext), os.O_WRONLY, 0644)
}

func GenerateAuthFiles(schemaPath string, dataPath string, outputPath string) error {
	schema, err := jsonschema.GetJSONSchema(schemaPath)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	bytes, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	for name, prop := range schema.Properties {
		if prop.GenerateAuthFile != nil {
			if data[name] == nil {
				return errors.New("schema property doesn't exist in data file: " + name)
			}

			out, err := OutputGenerator(outputPath, name, prop.GenerateAuthFile.Format)
			if err != nil {
				return err
			}

			if prop.GenerateAuthFile.Template != nil {
				err = renderTemplate(out, data[name], *prop.GenerateAuthFile.Template)
				if err != nil {
					return err
				}
			} else {
				switch prop.GenerateAuthFile.Format {
				case "ini":
					err = renderINI(out, data[name])
					if err != nil {
						return err
					}
				case "json":
					err = renderJSON(out, data[name])
					if err != nil {
						return err
					}
				case "yaml":
					err = renderYAML(out, data[name])
					if err != nil {
						return err
					}
				default:
					return errors.New("unrecognized file format " + prop.GenerateAuthFile.Format)
				}
			}
		}
	}
	return nil
}

func renderTemplate(out io.Writer, data interface{}, templateFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(out, path.Base(templateFile), data)
}

func renderINI(out io.Writer, data interface{}) error {
	return errors.New("not implemented yet")
}

func renderJSON(out io.Writer, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = out.Write(bytes)
	return err
}

func renderYAML(out io.Writer, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	_, err = out.Write(bytes)
	return err
}
