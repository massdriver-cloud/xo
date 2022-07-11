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

	"github.com/open-policy-agent/opa/rego"
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
		util.LogError(err, span ,"an error occurred while reading OPA result")
	}

	var output OPAOutput

	err = json.Unmarshal(buf, &output)
	if err != nil {
		util.LogError(err, span ,"an error occurred while parsing OPA result")
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
				util.LogError(err, span ,"an error occurred while parsing OPA result")
				errCounter++
			}

			if event != nil {
				pubErr := client.PublishEventToSNS(event)
				if pubErr != nil {
					util.LogError(err, span, "an error occurred while sending OPA result to massdriver")
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

func expressionToEventPayloadResourceProgress(expression *rego.ExpressionValue, deploymentId string) (*massdriver.EventPayloadResourceProgress, error) {
	expVal, ok := expression.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for expression value: %v expected map got %T", expression.Value, expression.Value)
	}
	// cannot type assert map[string]string so we do this in our own function iteratively.
	r, strErr := stringifiedMapValues(expVal)
	if strErr != nil {
		return nil, strErr
	}

	payload := new(massdriver.EventPayloadResourceProgress)

	payload.DeploymentId = deploymentId
	payload.ResourceId, ok = r["resource_id"]
	if !ok {
		return nil, fmt.Errorf("missing expected key resource_id got %v", expression.Value)
	}
	payload.ResourceType, ok = r["resource_type"]
	if !ok {
		return nil, fmt.Errorf("missing expected key resource_type got %v", expression.Value)
	}
	payload.ResourceName, ok = r["resource_name"]
	if !ok {
		return nil, fmt.Errorf("missing expected key resource_name got %v", expression.Value)
	}
	payload.ResourceKey, ok = r["resource_key"]
	if !ok {
		return nil, fmt.Errorf("missing expected key resource_key got %v", expression.Value)
	}

	return payload, nil
}

func convertResultExpressionToMassdriverEvent(span trace.Span, expression *rego.ExpressionValue, deploymentId string) (*massdriver.Event, error) {
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
