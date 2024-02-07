package env

import (
	"github.com/itchyny/gojq"
)

func GenerateEnvs(envs map[string]string, params, connections map[string]interface{}) (map[string]string, error) {
	result := map[string]string{}

	combined := map[string]interface{}{
		"params":      params,
		"connections": connections,
	}

	for name, query := range envs {
		foo, err := gojq.Parse(query)
		if err != nil {
			return result, err
		}

		iter := foo.Run(combined)

		value, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := value.(error); ok {
			return result, err
		}
		// Maybe check if there are multiple results and error?
		result[name] = value.(string)
	}

	return result, nil
}
