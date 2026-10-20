package resource

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const DefaultNameFormat = "Resource {field} for {instance}"

type ProvisionInput struct {
	InstanceID string
	Bundle     bundle.Bundle
	StepPath   string
	// TemplateDir is searched for resource_<field>.jq and artifact_<field>.jq files.
	TemplateDir string
	// Inputs is the object templates are evaluated against:
	// {params, connections, envs, secrets, outputs}.
	Inputs map[string]any
	// NameFormat is the human friendly resource name. "{field}" and
	// "{instance}" are replaced with the resource field and instance ID.
	NameFormat string
}

// Provision publishes every resource the current step has a template for, then
// deletes resources on the instance whose field is no longer declared in the
// bundle (i.e. the resource was removed or renamed in a bundle update).
func Provision(ctx context.Context, svc ResourceService, input ProvisionInput) (*Result, error) {
	ctx, span := otel.Tracer("xo").Start(ctx, "ResourceProvision")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	result, err := provision(ctx, svc, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

func provision(ctx context.Context, svc ResourceService, input ProvisionInput) (*Result, error) {
	step, err := findStep(input.Bundle, input.StepPath)
	if err != nil {
		return nil, err
	}

	templates, warnings, err := LoadTemplates(step.Resources, input.TemplateDir)
	for _, w := range warnings {
		log.Warn().Msg(w)
	}
	if err != nil {
		return nil, err
	}

	declared := input.Bundle.ResourceFields()
	var undeclared []string
	for _, tmpl := range templates {
		if !declared[tmpl.Field] {
			undeclared = append(undeclared, fmt.Sprintf("%q (from %s)", tmpl.Field, tmpl.Origin))
		}
	}
	if len(undeclared) > 0 {
		return nil, fmt.Errorf("templates exist for resources not declared in massdriver.yaml: %s (declared: %s)",
			strings.Join(undeclared, ", "), strings.Join(sortedKeys(declared), ", "))
	}

	payloads := make(map[string]map[string]any, len(templates))
	var renderErrs []error
	for _, tmpl := range templates {
		payload, renderErr := tmpl.Render(input.Inputs)
		if renderErr != nil {
			renderErrs = append(renderErrs, renderErr)
			continue
		}
		for _, path := range NullPaths(payload) {
			log.Warn().Msgf("resource %q: template from %s produced null at %s (check for a typo or missing input)", tmpl.Field, tmpl.Origin, path)
		}
		payloads[tmpl.Field] = payload
	}
	if len(renderErrs) > 0 {
		return nil, errors.Join(renderErrs...)
	}

	// Listing is only needed to find stale resources. Older (self-hosted)
	// servers don't support it, so a failure skips cleanup rather than failing
	// the deployment.
	existing, err := svc.ListResources(ctx)
	if err != nil {
		existing = nil
	}

	result := &Result{}

	nameFormat := input.NameFormat
	if nameFormat == "" {
		nameFormat = DefaultNameFormat
	}
	for _, tmpl := range templates {
		name := strings.NewReplacer("{field}", tmpl.Field, "{instance}", input.InstanceID).Replace(nameFormat)
		log.Info().Msgf("Publishing resource \"%s.%s\"...", input.InstanceID, tmpl.Field)
		if pubErr := Publish(ctx, svc, payloads[tmpl.Field], tmpl.Field, name); pubErr != nil {
			return result, fmt.Errorf("unable to publish resource %q: %w", tmpl.Field, pubErr)
		}
		result.Published = append(result.Published, tmpl.Field)
	}

	var errs []error
	for _, res := range existing {
		if res.Field == "" || declared[res.Field] {
			continue
		}
		id := resourceID(input.InstanceID, res)
		log.Info().Msgf("Deleting resource \"%s\": field %q is no longer declared in massdriver.yaml", id, res.Field)
		if delErr := Delete(ctx, svc, id); delErr != nil {
			errs = append(errs, fmt.Errorf("unable to delete stale resource %q: %w", id, delErr))
			continue
		}
		result.Deleted = append(result.Deleted, res.Field)
	}

	return result, errors.Join(errs...)
}
