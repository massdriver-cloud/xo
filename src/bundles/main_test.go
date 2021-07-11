package bundles_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"xo/src/bundles"
)

func TestBuild(t *testing.T) {
	var bundle, _ = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	_ = bundle.Build("./tmp/")

	gotDir, _ := os.ReadDir("./tmp")
	got := []string{}

	for _, dirEntry := range gotDir {
		got = append(got, dirEntry.Name())
	}

	want := []string{"schema-artifacts.json", "schema-connections.json", "schema-params.json"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	defer os.RemoveAll("./tmp")
}

func TestBuildSchema(t *testing.T) {
	var bundle, _ = bundles.ParseBundle("./testdata/bundle.Build/bundle.yaml")
	var inputIo bytes.Buffer

	bundles.BuildSchema(bundle.Params, bundle.Metadata("params"), &inputIo)
	var gotJson = &map[string]interface{}{}
	_ = json.Unmarshal(inputIo.Bytes(), gotJson)

	wantBytes, _ := ioutil.ReadFile("./testdata/bundle.Build/schema-params.json")
	var wantJson = &map[string]interface{}{}
	_ = json.Unmarshal(wantBytes, wantJson)

	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("got %v, want %v", gotJson, wantJson)
	}
}

// func TestParseBundle(t *testing.T) {
// 	ppn := yaml.MapSlice{}
// 	ppn = append(ppn, yaml.MapItem{Key: "type", Value: "string"})
// 	ppn = append(ppn, yaml.MapItem{Key: "title", Value: "name"})
// 	ppa := yaml.MapSlice{}
// 	ppa = append(ppa, yaml.MapItem{Key: "type", Value: "integer"})
// 	ppa = append(ppa, yaml.MapItem{Key: "title", Value: "age"})
// 	pps := yaml.MapSlice{}
// 	pps = append(pps, yaml.MapItem{Key: "name", Value: ppn})
// 	pps = append(pps, yaml.MapItem{Key: "age", Value: ppa})
// 	pp := yaml.MapItem{Key: "properties", Value: pps}
// 	ps := yaml.MapSlice{}
// 	ps = append(ps, yaml.MapItem{Key: "required", Value: [1]string{"name"}})
// 	ps = append(ps, pp)
// 	//params := yaml.MapSlice{}
// 	//params = append(params, yaml.MapItem{Key: "params", Value: ps})

// 	cpd := yaml.MapSlice{}
// 	cpd = append(cpd, yaml.MapItem{Key: "type", Value: "string"})
// 	cpd = append(cpd, yaml.MapItem{Key: "title", Value: "Default credential"})
// 	cps := yaml.MapSlice{}
// 	cps = append(cps, yaml.MapItem{Key: "default", Value: cpd})
// 	cp := yaml.MapItem{Key: "properties", Value: cps}
// 	cs := yaml.MapSlice{}
// 	cs = append(cs, yaml.MapItem{Key: "required", Value: [1]string{"default"}})
// 	cs = append(cs, cp)
// 	//conns := yaml.MapSlice{}
// 	//conns = append(conns, yaml.MapItem{Key: "connections", Value: cs})

// 	var got, _ = bundles.ParseBundle("./testdata/bundle.yaml")
// 	var want = bundles.Bundle{
// 		Uuid:        "FC2C7101-86A6-437B-B8C2-A2391FE8C847",
// 		Schema:      "draft-07",
// 		Type:        "aws-vpc",
// 		Title:       "AWS VPC",
// 		Description: "Something",
// 		Artifacts:   yaml.MapSlice{},
// 		Params:      ps,
// 		Connections: cs,
// 		// Params: map[string]interface{}{
// 		// 	"properties": map[string]interface{}{
// 		// 		"name": map[string]interface{}{
// 		// 			"type":  "string",
// 		// 			"title": "Name",
// 		// 		},
// 		// 		"age": map[string]interface{}{
// 		// 			"type":  "integer",
// 		// 			"title": "Age",
// 		// 		},
// 		// 	},
// 		// 	"required": []interface{}{
// 		// 		"name",
// 		// 	},
// 		// },
// 		// Connections: map[string]interface{}{
// 		// 	"required": []interface{}{
// 		// 		"default",
// 		// 	},
// 		// 	"properties": map[string]interface{}{
// 		// 		"default": map[string]interface{}{
// 		// 			"type":  "string",
// 		// 			"title": "Default credential",
// 		// 		},
// 		// 	},
// 		// },
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("got %v, want %v", got, want)
// 	}
// }
