package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type WeakSchema map[string]interface{}

type Bundle struct {
	Schema      string     `json:"schema"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Artifacts   WeakSchema `json:"artifacts"`
	Inputs      WeakSchema `json:"inputs"`
	Connections WeakSchema `json:"connections"`
}

func ParseBundle(path string) Bundle {
	bundle := Bundle{}

	data, err := ioutil.ReadFile(path)
	checkErr(err)

	err = yaml.Unmarshal([]byte(data), &bundle)
	checkErr(err)

	hydratedArtifacts := Hydrate(bundle.Artifacts)
	bundle.Artifacts = hydratedArtifacts.(map[string]interface{})

	hydratedInputs := Hydrate(bundle.Inputs)
	bundle.Inputs = hydratedInputs.(map[string]interface{})

	hydratedConnections := Hydrate(bundle.Connections)
	bundle.Connections = hydratedConnections.(map[string]interface{})

	hydratedBundle, err := json.Marshal(bundle)
	checkErr(err)
	fmt.Printf(string(hydratedBundle))
	return bundle
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
