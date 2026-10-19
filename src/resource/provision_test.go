package resource_test

import (
	"context"
	"fmt"
	"testing"
	"xo/src/bundle"
	"xo/src/resource"

	"github.com/massdriver-cloud/massdriver-sdk-go/massdriver/provisioning/resources"
	"github.com/stretchr/testify/require"
)

func provisionBundle(templates map[string]string) bundle.Bundle {
	stepResources := map[string]bundle.StepResource{}
	for field, tmpl := range templates {
		stepResources[field] = bundle.StepResource{Template: tmpl}
	}
	return bundle.Bundle{
		Steps: []bundle.BundleStep{
			{Path: "src", Provisioner: "opentofu", Resources: stepResources},
			{Path: "app", Provisioner: "helm"},
		},
		Resources: map[string]bundle.Resource{
			"reader": {ResourceType: "postgres"},
			"writer": {ResourceType: "postgres"},
			"app":    {ResourceType: "endpoint"},
		},
		Artifacts: map[string]interface{}{},
	}
}

func provisionInputs() map[string]any {
	return map[string]any{
		"params":      map[string]any{"region": "us-east-1"},
		"connections": map[string]any{},
		"envs":        map[string]any{},
		"secrets":     map[string]any{},
		"outputs":     map[string]any{"reader_host": "r.example.com", "writer_host": "w.example.com"},
	}
}

func TestProvision(t *testing.T) {
	resSvc := &recordingResourceService{Existing: []resources.Resource{
		{ID: "inst.database", Field: "database"}, // removed from bundle -> deleted
		{ID: "inst.reader", Field: "reader"},
		{ID: "inst.app", Field: "app"}, // declared, published by another step -> kept
	}}

	result, err := resource.Provision(context.Background(), resSvc, resource.ProvisionInput{
		InstanceID: "inst",
		Bundle: provisionBundle(map[string]string{
			"reader": `{ host: .outputs.reader_host }`,
			"writer": `{ host: .outputs.writer_host, region: .params.region }`,
		}),
		StepPath:    "src",
		TemplateDir: t.TempDir(),
		Inputs:      provisionInputs(),
	})
	require.NoError(t, err)

	require.Equal(t, []string{"reader", "writer"}, result.Published)
	require.Equal(t, []string{"database"}, result.Deleted)
	require.Equal(t, []string{"inst.database"}, resSvc.Deleted)

	require.Len(t, resSvc.Created, 2)
	require.Equal(t, resources.ResourceInput{
		Field:   "writer",
		Name:    "Resource writer for inst",
		Payload: map[string]any{"host": "w.example.com", "region": "us-east-1"},
	}, resSvc.Created[1])
}

func TestProvisionContinuesPastFailedDeletes(t *testing.T) {
	resSvc := &recordingResourceService{
		DeleteErr: map[string]error{"inst.database": fmt.Errorf("boom")},
		Existing: []resources.Resource{
			{ID: "inst.database", Field: "database"},
			{ID: "inst.legacy", Field: "legacy"},
		},
	}

	result, err := resource.Provision(context.Background(), resSvc, resource.ProvisionInput{
		InstanceID:  "inst",
		Bundle:      provisionBundle(map[string]string{"reader": `{}`}),
		StepPath:    "src",
		TemplateDir: t.TempDir(),
		Inputs:      provisionInputs(),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), `unable to delete stale resource "inst.database": boom`)
	require.Equal(t, []string{"reader"}, result.Published)
	require.Equal(t, []string{"inst.legacy"}, resSvc.Deleted)
	require.Equal(t, []string{"legacy"}, result.Deleted)
}

func TestProvisionNameFormat(t *testing.T) {
	resSvc := &recordingResourceService{}
	_, err := resource.Provision(context.Background(), resSvc, resource.ProvisionInput{
		InstanceID:  "inst",
		Bundle:      provisionBundle(map[string]string{"reader": `{}`}),
		StepPath:    "src",
		TemplateDir: t.TempDir(),
		Inputs:      provisionInputs(),
		NameFormat:  "{instance} {field} db",
	})
	require.NoError(t, err)
	require.Equal(t, "inst reader db", resSvc.Created[0].Name)
}

func TestProvisionFailures(t *testing.T) {
	tests := []struct {
		name          string
		stepResources map[string]string
		stepPath      string
		createErr     error
		listErr       error
		wantErr       string
		wantListed    bool
	}{
		{
			name:     "unknown step",
			stepPath: "nope",
			wantErr:  `step "nope" not found in massdriver.yaml (steps: src, app)`,
		},
		{
			name:          "undeclared field",
			stepResources: map[string]string{"database": `{}`},
			wantErr:       `templates exist for resources not declared in massdriver.yaml: "database"`,
		},
		{
			name: "invalid template publishes nothing",
			stepResources: map[string]string{
				"reader": `{ host: .outputs.reader_host }`,
				"writer": `{ host: }`,
			},
			wantErr: `resource "writer"`,
		},
		{
			name: "render error publishes nothing",
			stepResources: map[string]string{
				"reader": `{ host: .outputs.reader_host }`,
				"writer": `.outputs.writer_host`,
			},
			wantErr: "expected a JSON object",
		},
		{
			name:          "list failure publishes nothing",
			stepResources: map[string]string{"reader": `{}`},
			listErr:       fmt.Errorf("nope"),
			wantErr:       `unable to list resources for instance "inst": nope`,
			wantListed:    true,
		},
		{
			name:          "publish failure skips pruning",
			stepResources: map[string]string{"reader": `{}`},
			createErr:     fmt.Errorf("boom"),
			wantErr:       `unable to publish resource "reader": boom`,
			wantListed:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stepPath := tc.stepPath
			if stepPath == "" {
				stepPath = "src"
			}
			resSvc := &recordingResourceService{
				CreateErr: tc.createErr,
				Existing:  []resources.Resource{{ID: "inst.database", Field: "database"}},
				ListErr:   tc.listErr,
			}

			_, err := resource.Provision(context.Background(), resSvc, resource.ProvisionInput{
				InstanceID:  "inst",
				Bundle:      provisionBundle(tc.stepResources),
				StepPath:    stepPath,
				TemplateDir: t.TempDir(),
				Inputs:      provisionInputs(),
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
			require.Empty(t, resSvc.Created)
			require.Empty(t, resSvc.Deleted)
			require.Equal(t, tc.wantListed, resSvc.ListCalled)
		})
	}
}
