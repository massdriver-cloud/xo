package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"xo/src/massdriver"
	"xo/src/telemetry"
	"xo/src/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ReportResults(ctx context.Context, client *massdriver.MassdriverClient, deploymentId string, stream io.Reader) error {
	_, span := otel.Tracer("xo").Start(ctx, "provisioners.opa.ReportResults")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		util.LogError(err, span, "an error occurred while reading OPA result")
	}

	var output OPAOutput

	err = json.Unmarshal(buf, &output)
	if err != nil {
		util.LogError(err, span, "an error occurred while parsing OPA result")
		return err
	}

	// aggregate errors so returning early doesn't prevent publishing of events due to error on an earlier expression.
	errCounter := 0
	for _, result := range output.Result {
		for _, expression := range result.Expressions {
			// normally we'd pass context instead of a span (probably an antipattern) but for readability
			// its much easier to annotate the same span, so we pass that instead
			event, convErr := convertResultExpressionToMassdriverEvent(span, expression, deploymentId)
			if convErr != nil {
				util.LogError(err, span, "an error occurred while parsing OPA result")
				errCounter++
			}

			if event != nil {
				pubErr := client.PublishEvent(event)
				if pubErr != nil {
					util.LogError(err, span, "an error occurred while sending OPA event to massdriver")
					errCounter++
				}
			}

			diagnostic, convErr := expressionToEventPayloadDiagnostic(expression, deploymentId)
			if convErr != nil {
				util.LogError(err, span, "an error occurred while parsing OPA result")
				errCounter++
			}

			if diagnostic != nil {
				event := massdriver.NewEvent(massdriver.EVENT_TYPE_PROVISIONER_ERROR)
				event.Payload = diagnostic
				pubErr := client.PublishEvent(event)
				if pubErr != nil {
					util.LogError(err, span, "an error occurred while sending OPA diagnostic to massdriver")
					errCounter++
				}
			}
		}
	}
	if errCounter > 0 {
		// this serves the purpose of returning an error to the caller to inidicate the were a non-zero expressions that were not converted / published as events.
		// the individual errors were logged above during the loop.
		return fmt.Errorf("%d errors occurred while reporting OPA results", errCounter)
	}
	return nil
}

func stringifiedMapValues(m map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string)
	var isStr bool
	for k, v := range m {
		out[k], isStr = v.(string)
		if !isStr {
			if n, isInt := v.(int); isInt {
				out[k] = strconv.Itoa(n)
				continue
			}
			// for some reason ints are sometimes float64 in this map[string]interface{}
			if n, isFloat := v.(float64); isFloat {
				out[k] = strconv.Itoa(int(n))
				continue
			}

			return nil, fmt.Errorf("unexpected type for key %s: %v expected string or int got %T", k, v, v)
		}
	}
	return out, nil
}

func expressionToEventPayloadResourceProgress(expression OPAExpression, deploymentId string) (*massdriver.EventPayloadResourceProgress, error) {
	payload := new(massdriver.EventPayloadResourceProgress)

	payload.DeploymentId = deploymentId
	payload.ResourceId = expression.Value.ResourceID
	payload.ResourceType = expression.Value.ResourceType
	payload.ResourceName = expression.Value.ResourceName
	payload.ResourceKey = expression.Value.ResourceKey

	return payload, nil
}

func expressionToEventPayloadDiagnostic(expression OPAExpression, deploymentId string) (*massdriver.EventPayloadDiagnostic, error) {
	payload := new(massdriver.EventPayloadDiagnostic)

	payload.DeploymentId = deploymentId
	payload.Message = expressionTextToRule(expression.Text)
	details, err := json.Marshal(expression.Value)
	if err == nil {
		payload.Details = string(details)
	} else {
		// if for some reason we can't marshal the expression value, we'll just print the golang representation of the value as this will give more info than nothing.
		payload.Details = fmt.Sprintf("%v", expression.Value)
	}
	payload.Level = "error"

	return payload, nil
}

func convertResultExpressionToMassdriverEvent(span trace.Span, expression OPAExpression, deploymentId string) (*massdriver.Event, error) {
	event := massdriver.NewEvent(massdriver.EVENT_TYPE_RESOURCE_UPDATE_FAILED)
	opaViolationPayload, convErr := expressionToEventPayloadResourceProgress(expression, deploymentId)
	if convErr != nil {
		util.LogError(convErr, span, "an error occurred while converting OPA result to massdriver event")
		return nil, convErr
	}

	err := fmt.Errorf("%s on %#v", expressionTextToRule(expression.Text), opaViolationPayload)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	event.Payload = opaViolationPayload

	return event, nil
}

func expressionTextToRule(text string) string {
	switch text {
	case "data.terraform.deletion_violations":
		return "Deletion Violation"
	default:
		return text
	}
}
