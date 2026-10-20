package resource

import (
	"context"
	"errors"
	"fmt"
	"xo/src/bundle"
	"xo/src/telemetry"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type DecommissionInput struct {
	InstanceID string
	Bundle     bundle.Bundle
	StepPath   string
	// TemplateDir is searched for resource_<field>.jq and artifact_<field>.jq files.
	TemplateDir string
}

// Decommission deletes the instance's resources that the current step manages
// (it has an inline or file template for the field), plus any resource whose
// field is no longer declared in the bundle. Resources managed by other steps
// are left for those steps to delete.
//
// Every matching resource is attempted even if some deletions fail; the
// failures are returned together.
func Decommission(ctx context.Context, svc ResourceService, input DecommissionInput) (*Result, error) {
	ctx, span := otel.Tracer("xo").Start(ctx, "ResourceDecommission")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	result, err := decommission(ctx, svc, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

func decommission(ctx context.Context, svc ResourceService, input DecommissionInput) (*Result, error) {
	step, err := findStep(input.Bundle, input.StepPath)
	if err != nil {
		return nil, err
	}

	// Templates aren't compiled here so a broken template can't block cleanup.
	fields, err := TemplateFields(step.Resources, input.TemplateDir)
	if err != nil {
		return nil, err
	}
	managed := map[string]bool{}
	for _, field := range fields {
		managed[field] = true
	}
	declared := input.Bundle.ResourceFields()

	// Older self-hosted) servers don't support listing. Fall back to deleting
	// the step's managed resources by ID and skip cleanup of stale resources.
	existing, err := svc.ListResources(ctx)
	if err != nil {
		existing = make([]resources.Resource, 0, len(fields))
		for _, field := range fields {
			existing = append(existing, resources.Resource{Field: field})
		}
	}

	result := &Result{}
	var errs []error
	for _, res := range existing {
		var reason string
		switch {
		case managed[res.Field]:
			reason = fmt.Sprintf("managed by step %q", input.StepPath)
		case !declared[res.Field]:
			reason = fmt.Sprintf("field %q is no longer declared in massdriver.yaml", res.Field)
		default:
			continue
		}
		id := resourceID(input.InstanceID, res)
		log.Info().Msgf("Deleting resource \"%s\": %s", id, reason)
		delErr := Delete(ctx, svc, id)
		if errors.Is(delErr, resources.ErrNotFound) {
			log.Info().Msgf("Resource \"%s\" does not exist, nothing to delete", id)
			continue
		}
		if delErr != nil {
			errs = append(errs, fmt.Errorf("unable to delete resource %q: %w", id, delErr))
			continue
		}
		result.Deleted = append(result.Deleted, res.Field)
	}

	return result, errors.Join(errs...)
}
