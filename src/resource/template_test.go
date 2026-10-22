package resource_test

import (
	"testing"
	"xo/src/bundle"
	"xo/src/resource"

	"github.com/stretchr/testify/require"
)

func TestLoadTemplates(t *testing.T) {
	inline := map[string]bundle.StepResource{
		"billing":  {Template: `{ url: .outputs.function_url, region: .params.region }`},
		"database": {}, // no inline template -> falls back to resource_database.jq
	}

	templates, warnings, err := resource.LoadTemplates(inline, "testdata/templates")
	require.NoError(t, err)

	origins := map[string]string{}
	for _, tmpl := range templates {
		origins[tmpl.Field] = tmpl.Origin
	}
	require.Equal(t, map[string]string{
		"billing":  "massdriver.yaml steps[].resources.billing.template",
		"database": "testdata/templates/resource_database.jq",
		"legacy":   "testdata/templates/artifact_legacy.jq",
	}, origins)
	require.Len(t, warnings, 1, "file template shadowed by inline template should warn")
	require.Contains(t, warnings[0], "resource_billing.jq")
}

func TestTemplateFields(t *testing.T) {
	fields, err := resource.TemplateFields(map[string]bundle.StepResource{
		"billing": {Template: `{ not valid jq`},
		"other":   {},
	}, "testdata/templates")
	require.NoError(t, err)
	require.Equal(t, []string{"billing", "database", "legacy"}, fields)
}

func TestLoadTemplatesErrors(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		contains []string
	}{
		{
			name:     "parse error reports line and caret",
			source:   "{\n  url: .outputs.url,\n  region: .params.region,,\n}",
			contains: []string{`resource "broken"`, "massdriver.yaml steps[].resources.broken.template", "failed to parse", "line 3, column 26", "  3 |   region: .params.region,,\n    |                          ^"},
		},
		{
			name:     "unterminated object",
			source:   "{ url: .outputs.url",
			contains: []string{"failed to parse", "unexpected EOF"},
		},
		{
			name:     "undefined variable hints at input object",
			source:   "{ url: $outputs.url }",
			contains: []string{"failed to compile", "variable not defined: $outputs", "hint: inputs are available"},
		},
		{
			name:     "undefined function",
			source:   "{ url: foo(.outputs.url) }",
			contains: []string{"failed to compile", "function not defined: foo/1"},
		},
		{
			name:     "empty template",
			source:   "  \n",
			contains: []string{"template is empty"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := resource.LoadTemplates(map[string]bundle.StepResource{"broken": {Template: tc.source}}, t.TempDir())
			require.Error(t, err)
			for _, want := range tc.contains {
				require.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestLoadTemplatesReportsAllErrors(t *testing.T) {
	_, _, err := resource.LoadTemplates(map[string]bundle.StepResource{
		"one": {Template: "{ a: }"},
		"two": {Template: "{ b: $nope }"},
	}, t.TempDir())
	require.Error(t, err)
	require.Contains(t, err.Error(), `resource "one"`)
	require.Contains(t, err.Error(), `resource "two"`)
}

func TestTemplateRender(t *testing.T) {
	inputs := map[string]any{
		"params":  map[string]any{"region": "us-east-1"},
		"outputs": map[string]any{"function_url": "https://example.com", "list": []any{1.0, 2.0}},
	}

	tests := []struct {
		name     string
		source   string
		want     map[string]any
		contains []string
	}{
		{
			name:   "renders object",
			source: `{ url: .outputs.function_url, region: .params.region }`,
			want:   map[string]any{"url": "https://example.com", "region": "us-east-1"},
		},
		{
			name:     "multiple outputs",
			source:   `.outputs.list[] | {n: .}`,
			contains: []string{"produced 2 outputs", "wrap it in [...]"},
		},
		{
			name:     "non-object output",
			source:   `.outputs.function_url`,
			contains: []string{"produced a string (expected a JSON object)"},
		},
		{
			name:     "no output",
			source:   `empty`,
			contains: []string{"produced no output"},
		},
		{
			name:     "runtime error lists available inputs",
			source:   `{ a: (.outputs.function_url | keys) }`,
			contains: []string{"failed to evaluate", ".outputs keys: function_url, list", ".params keys: region"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			templates, _, err := resource.LoadTemplates(map[string]bundle.StepResource{"res": {Template: tc.source}}, t.TempDir())
			require.NoError(t, err)
			got, err := templates[0].Render(inputs)
			if tc.contains != nil {
				require.Error(t, err)
				for _, want := range tc.contains {
					require.Contains(t, err.Error(), want)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestTemplateRenderArrayOutputs(t *testing.T) {
	// Helm outputs are an array of rendered manifests rather than an object.
	inputs := map[string]any{
		"params": map[string]any{"region": "us-east-1"},
		"outputs": []any{
			map[string]any{"kind": "Service", "metadata": map[string]any{"name": "web"}},
			map[string]any{"kind": "Deployment", "metadata": map[string]any{"name": "web"}},
			map[string]any{"kind": "Service", "metadata": map[string]any{"name": "web-headless"}},
		},
	}

	t.Run("selects from manifests", func(t *testing.T) {
		templates, _, err := resource.LoadTemplates(map[string]bundle.StepResource{
			"svc": {Template: `{ name: ([.outputs[] | select(.kind == "Service")][0].metadata.name) }`},
		}, t.TempDir())
		require.NoError(t, err)
		got, err := templates[0].Render(inputs)
		require.NoError(t, err)
		require.Equal(t, map[string]any{"name": "web"}, got)
	})

	t.Run("runtime error describes array outputs", func(t *testing.T) {
		templates, _, err := resource.LoadTemplates(map[string]bundle.StepResource{
			"svc": {Template: `{ url: .outputs.url }`},
		}, t.TempDir())
		require.NoError(t, err)
		_, err = templates[0].Render(inputs)
		require.Error(t, err)
		require.Contains(t, err.Error(), ".outputs: array of 3 items, kinds: Deployment, Service")
		require.Contains(t, err.Error(), ".params keys: region")
	})

	t.Run("runtime error describes empty array outputs", func(t *testing.T) {
		templates, _, err := resource.LoadTemplates(map[string]bundle.StepResource{
			"svc": {Template: `{ url: .outputs.url }`},
		}, t.TempDir())
		require.NoError(t, err)
		_, err = templates[0].Render(map[string]any{"outputs": []any{}})
		require.Error(t, err)
		require.Contains(t, err.Error(), ".outputs: (empty array)")
	})
}

func TestNullPaths(t *testing.T) {
	got := resource.NullPaths(map[string]any{
		"a": nil,
		"b": map[string]any{"c": nil, "d": "ok"},
		"e": []any{"ok", nil},
	})
	require.Equal(t, []string{".a", ".b.c", ".e[1]"}, got)
}
