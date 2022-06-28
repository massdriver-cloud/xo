package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"xo/src/massdriver"
	"xo/src/telemetry"

	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OPAOutput struct {
	Result rego.ResultSet `json:"result"`
}

func ReportResults(ctx context.Context, client *massdriver.MassdriverClient, deploymentId string, stream io.Reader) error {
	_, span := otel.Tracer("xo").Start(ctx, "provisioners.opa.ReportResults")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while reading OPA result")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	var output OPAOutput

	err = json.Unmarshal(buf, &output)
	if err != nil {
		log.Error().Err(err).Msg("an error occurred while parsing OPA result")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	for _, result := range output.Result {
		for _, expression := range result.Expressions {
			// normally we'd pass context instead of a span (probably an antipattern) but for readability
			// its much easier to annotate the same span, so we pass that instead
			event := convertResultExpressionToMassdriverEvent(span, expression, deploymentId)

			if event != nil {
				err = client.PublishEventToSNS(event)
				if err != nil {
					log.Error().Err(err).Msg("an error occurred while sending OPA result to massdriver")
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
			}
		}
	}

	return nil
}

func convertResultExpressionToMassdriverEvent(span trace.Span, expression *rego.ExpressionValue, deploymentId string) *massdriver.Event {
	event := massdriver.NewEvent(massdriver.EVENT_TYPE_PROVISIONER_VIOLATION)
	violation := new(massdriver.EventPayloadOPAViolation)
	violation.DeploymentId = deploymentId
	violation.Rule = expressionTextToRule(expression.Text)
	violation.Value = expression.Value

	err := fmt.Errorf("%s on %v", violation.Rule, violation.Value)

	event.Payload = violation

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return event
}

func expressionTextToRule(text string) string {
	switch text {
	case "data.terraform.deletion_violations":
		return "Deletion Violation"
	default:
		return text
	}
}
