package terraform

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type terraformLog struct {
	Level     string        `json:"@level"`
	Message   string        `json:"@message"`
	Module    string        `json:"@module"`
	Timestamp string        `json:"@timestamp"`
	Hook      terraformHook `json:"hook"`
	Type      string        `json:"type"`
}

type terraformHook struct {
	Resource terraformResourceAddr `json:"resource"`
	Action   string                `json:"action"`
	IDKey    string                `json:"id_key,omitempty"`
	IDValue  string                `json:"id_value,omitempty"`
	Elapsed  float64               `json:"elapsed_seconds"`
}

type terraformResourceAddr struct {
	Addr            string                  `json:"addr"`
	Module          string                  `json:"module"`
	Resource        string                  `json:"resource"`
	ImpliedProvider string                  `json:"implied_provider"`
	ResourceType    string                  `json:"resource_type"`
	ResourceName    string                  `json:"resource_name"`
	ResourceKey     ctyjson.SimpleJSONValue `json:"resource_key"`
}

func Extract(stream io.Reader) error {
	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		var record terraformLog

		err := json.Unmarshal([]byte(scanner.Text()), &record)
		if err != nil {
			return err
		}

		fmt.Println(record.Hook.Action)
	}

	return nil
}
