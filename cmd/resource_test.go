package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestReadJSON(t *testing.T) {
	got, err := readJSON(writeTempFile(t, `{"url": "https://example.com"}`))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"url": "https://example.com"}, got)

	got, err = readJSON(writeTempFile(t, `[{"kind": "Service"}]`))
	require.NoError(t, err)
	require.Equal(t, []any{map[string]any{"kind": "Service"}}, got)

	_, err = readJSON(writeTempFile(t, `{not json`))
	require.ErrorContains(t, err, "is not valid JSON")
}

func TestReadJSONObject(t *testing.T) {
	got, err := readJSONObject(writeTempFile(t, `{"region": "us-east-1"}`))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"region": "us-east-1"}, got)

	got, err = readJSONObject(writeTempFile(t, `null`))
	require.NoError(t, err)
	require.Equal(t, map[string]any{}, got)

	for _, content := range []string{`[{"region": "us-east-1"}]`, `"us-east-1"`, `1`} {
		_, err := readJSONObject(writeTempFile(t, content))
		require.ErrorContains(t, err, "must contain a JSON object", "content: %s", content)
	}
}
