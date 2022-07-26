package jsonschema

import "encoding/json"

func GetJSONSchema(path string) (Schema, error) {
	schema := Schema{}
	sl := Load(path)

	schemaSrc, err := sl.LoadJSON()
	if err != nil {
		return schema, err
	}

	byteData, errMarshal := json.Marshal(schemaSrc)
	if errMarshal != nil {
		return schema, errMarshal
	}

	if errUnmarsh := json.Unmarshal(byteData, &schema); errUnmarsh != nil {
		return schema, errUnmarsh
	}
	return schema, nil
}
