package env

import (
	"context"
	"xo/src/telemetry"

	"github.com/itchyny/gojq"
	"go.opentelemetry.io/otel"
)

func GenerateEnvs(ctx context.Context, envs map[string]string, params, connections map[string]interface{}) (map[string]string, error) {
	_, span := otel.Tracer("xo").Start(ctx, "GenerateEnvs")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	result := map[string]string{}

	combined := map[string]interface{}{
		"params":      params,
		"connections": connections,
	}

	for name, query := range envs {
		gojqQuery, err := gojq.Parse(query)
		if err != nil {
			return result, err
		}

		iter := gojqQuery.Run(combined)

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
