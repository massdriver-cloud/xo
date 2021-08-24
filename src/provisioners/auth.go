package provisioners

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"xo/src/jsonschema"

	"github.com/itchyny/gojq"
	"gopkg.in/ini.v1"
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

			outputData := data[name]
			if prop.GenerateAuthFile.Template != nil {
				outputData, err = renderTemplate(data[name], *prop.GenerateAuthFile.Template)
				if err != nil {
					return err
				}
			}

			switch prop.GenerateAuthFile.Format {
			case "ini":
				err = renderINI(out, outputData)
				if err != nil {
					return err
				}
			case "json":
				err = renderJSON(out, outputData)
				if err != nil {
					return err
				}
			case "yaml":
				err = renderYAML(out, outputData)
				if err != nil {
					return err
				}
			default:
				return errors.New("unrecognized file format " + prop.GenerateAuthFile.Format)
			}
		}
	}
	return nil
}

func renderTemplate(data interface{}, template string) (interface{}, error) {
	query, err := gojq.Parse(template)
	if err != nil {
		return nil, err
	}
	iter := query.Run(data)
	// this iterator object is strange. I think its for iterating through the results
	// of complex manipulations (like turning a map into an array, and iterating through
	// the array elements). Since we're just doing simple transformations for now, I'm
	// making the assumption that everything we need is in the first iterator element.
	// This logic **PROBABLY** will fall apart if we start doing complex jq queries,
	// and we'll have to revisit it then.
	v, _ := iter.Next()
	if err, ok := v.(error); ok {
		return nil, err
	}
	return v, nil
}

func renderINI(out io.Writer, data interface{}) error {
	cfg := ini.Empty()

	err := generateIni(cfg, data, ini.DefaultSection)
	if err != nil {
		return err
	}

	ini.PrettyFormat = false
	_, err = cfg.WriteTo(out)
	return err
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

// The ini package we are using doesn't currently have the ability to Marshal/Unmarshal
// ini files to/from map[string]interface{}. There is a github issue:
// https://github.com/go-ini/ini/issues/275 but it isn't getting much traction. For now,
// our use-case is simple enough we can write our own "marshaler".
func generateIni(cfg *ini.File, data interface{}, sectionName string) error {
	d := data.(map[string]interface{})
	for k, v := range d {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			var newSectionName string
			if sectionName == ini.DefaultSection {
				newSectionName = k
			} else {
				newSectionName = sectionName + "." + k
			}
			cfg.NewSection(newSectionName)
			return generateIni(cfg, v, newSectionName)
		default:
			section, err := cfg.GetSection(sectionName)
			if err != nil {
				return err
			}
			section.Key(k).SetValue(fmt.Sprintf("%v", v))
		}
	}
	return nil
}
