package opa

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type OPAOutput struct {
	Result []OPAResult `json:"result"`
}

type OPAResult struct {
	Expressions []OPAExpression `json:"expressions"`
}

type OPAExpression struct {
	Value OPAResource `json:"value"`
	Text  string      `json:"text"`
}

type OPAResource struct {
	ResourceID   string `json:"resource_id"`
	ResourceKey  string `json:"resource_key"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
}

func (r *OPAResource) UnmarshalJSON(data []byte) error {

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "resource_id":
			r.ResourceID = v.(string)
		case "resource_key":
			r.ResourceKey = stringify(v)
		case "resource_name":
			r.ResourceName = v.(string)
		case "resource_type":
			r.ResourceType = v.(string)
		}
	}

	return nil
}

func stringify(x interface{}) string {
	switch v := x.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.Itoa(int(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}
