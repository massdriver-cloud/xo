package resource_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"xo/src/bundle"
	"xo/src/resource"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
	"github.com/stretchr/testify/require"
)

func decommissionBundle() bundle.Bundle {
	return bundle.Bundle{
		Steps: []bundle.BundleStep{
			{Path: "src", Resources: map[string]bundle.StepResource{
				"reader": {Template: `{}`},
				"writer": {Template: `{ broken: }`}, // broken templates must not block cleanup
			}},
			{Path: "app", Resources: map[string]bundle.StepResource{
				"app": {Template: `{}`},
			}},
		},
		Resources: map[string]bundle.Resource{
			"reader": {}, "writer": {}, "billing": {}, "app": {},
		},
		Artifacts: map[string]interface{}{},
	}
}

func TestDecommission(t *testing.T) {
	templateDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "resource_billing.jq"), []byte(`{}`), 0644))

	resSvc := &recordingResourceService{Existing: []resources.Resource{
		{ID: "inst.app", Field: "app"},           // declared, managed by step "app" -> kept
		{ID: "inst.billing", Field: "billing"},   // jq file in this step's dir -> deleted
		{ID: "inst.database", Field: "database"}, // no longer declared -> deleted
		{ID: "inst.reader", Field: "reader"},     // inline template in this step -> deleted
		{Field: "writer"},                        // inline template, no ID -> <instance>.<field>
	}}

	result, err := resource.Decommission(context.Background(), resSvc, resource.DecommissionInput{
		InstanceID:  "inst",
		Bundle:      decommissionBundle(),
		StepPath:    "src",
		TemplateDir: templateDir,
	})
	require.NoError(t, err)
	require.Equal(t, []string{"inst.billing", "inst.database", "inst.reader", "inst.writer"}, resSvc.Deleted)
	require.Equal(t, []string{"billing", "database", "reader", "writer"}, result.Deleted)
}

func TestDecommissionContinuesPastFailures(t *testing.T) {
	resSvc := &recordingResourceService{
		DeleteErr: map[string]error{"inst.reader": fmt.Errorf("boom")},
		Existing: []resources.Resource{
			{ID: "inst.reader", Field: "reader"},
			{ID: "inst.writer", Field: "writer"},
		},
	}

	result, err := resource.Decommission(context.Background(), resSvc, resource.DecommissionInput{
		InstanceID:  "inst",
		Bundle:      decommissionBundle(),
		StepPath:    "src",
		TemplateDir: t.TempDir(),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), `unable to delete resource "inst.reader": boom`)
	require.Equal(t, []string{"inst.writer"}, resSvc.Deleted)
	require.Equal(t, []string{"writer"}, result.Deleted)
}

func TestDecommissionFailures(t *testing.T) {
	tests := []struct {
		name     string
		stepPath string
		listErr  error
		wantErr  string
	}{
		{name: "unknown step", stepPath: "nope", wantErr: `step "nope" not found in massdriver.yaml (steps: src, app)`},
		{name: "list failure", stepPath: "src", listErr: fmt.Errorf("nope"), wantErr: `unable to list resources for instance "inst": nope`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resSvc := &recordingResourceService{ListErr: tc.listErr}
			_, err := resource.Decommission(context.Background(), resSvc, resource.DecommissionInput{
				InstanceID:  "inst",
				Bundle:      decommissionBundle(),
				StepPath:    tc.stepPath,
				TemplateDir: t.TempDir(),
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
			require.Empty(t, resSvc.Deleted)
		})
	}
}
