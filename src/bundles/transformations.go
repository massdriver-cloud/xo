package bundles

import (
	"errors"
)

func applyTransformations(schema map[string]interface{}) error {

	err := addSetIdToObjectArrays(schema)
	if err != nil {
		return err
	}
	err = disableAdditionalPropertiesInObjects(schema)
	if err != nil {
		return err
	}

	for _, v := range schema {
		_, isObject := v.(map[string]interface{})
		if isObject {
			err = applyTransformations(v.(map[string]interface{}))
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func addSetIdToObjectArrays(schema map[string]interface{}) error {
	if schema["type"] == "array" {
		itemsInterface, found := schema["items"]
		if !found {
			return errors.New("found array without items")
		}
		items := itemsInterface.(map[string]interface{})
		if items["type"] == "object" {
			propertiesInterface, found := items["properties"]
			if !found {
				return errors.New("found object without properties")
			}
			properties := propertiesInterface.(map[string]interface{})
			properties["md_set_id"] = map[string]string{"type": "string"}

			requiredInterface, found := items["required"]
			if !found {
				items["required"] = []string{"md_set_id"}
			} else {
				required := requiredInterface.([]interface{})
				items["required"] = append(required, "md_set_id")
			}
		}
	}
	return nil
}

func disableAdditionalPropertiesInObjects(schema map[string]interface{}) error {
	if schema["type"] == "object" {
		_, found := schema["additionalProperties"]
		if !found {
			schema["additionalProperties"] = false
		}
	}
	return nil
}
