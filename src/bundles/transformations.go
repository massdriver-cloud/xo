package bundles

import (
	"reflect"
	"xo/src/jsonschema"
)

func applyTransformations(oj *jsonschema.OrderedJSON) error {

	err := addSetIdToObjectArrays(oj)
	if err != nil {
		return err
	}

	for _, elem := range *oj {
		if reflect.TypeOf(elem.Value) == reflect.TypeOf(jsonschema.OrderedJSON{}) {
			local := elem.Value.(jsonschema.OrderedJSON)
			err = applyTransformations(&local)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func addSetIdToObjectArrays(oj *jsonschema.OrderedJSON) error {
	if oj.Type() == "array" {
		items := oj.GetItems()
		if items.Type() == "object" {
			properties := items.GetProperties()
			required := items.GetRequired()
			mdSetId := jsonschema.OrderedJSON{
				jsonschema.OrderedJSONElement{Key: "type", Value: "string"},
			}
			properties = append(properties, jsonschema.OrderedJSONElement{Key: "md_set_id", Value: mdSetId})
			required = append(required, "md_set_id")
			items.SetProperties(properties)
			items.SetRequired(required)
		}
		oj.SetItems(items)
	}
	return nil
}
