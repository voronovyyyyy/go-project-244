package code_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"code"

	"github.com/stretchr/testify/require"
)

func fixture(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func TestGenDiff(t *testing.T) {
	for _, first := range []string{"json", "yml", "yaml"} {
		for _, second := range []string{"json", "yml", "yaml"} {
			for _, format := range []string{"stylish", "plain", "json"} {
				t.Run(first+"/"+second+"/"+format, func(t *testing.T) {
					expected, err := os.ReadFile(fixture("expected." + format))
					require.NoError(t, err)
					actual, err := code.GenDiff(fixture("file1."+first), fixture("file2."+second), format)
					require.NoError(t, err)
					if format == "json" {
						require.JSONEq(t, string(expected), actual)
						return
					}
					require.Equal(t, strings.TrimSuffix(string(expected), "\n"), actual)
				})
			}
		}
	}
}

func TestGenDiffDefaultAndAbsolutePaths(t *testing.T) {
	first, err := filepath.Abs(fixture("file1.json"))
	require.NoError(t, err)
	second, err := filepath.Abs(fixture("file2.yaml"))
	require.NoError(t, err)
	expected, err := os.ReadFile(fixture("expected.stylish"))
	require.NoError(t, err)
	actual, err := code.GenDiff(first, second, "")
	require.NoError(t, err)
	require.Equal(t, strings.TrimSuffix(string(expected), "\n"), actual)
}

func TestGenDiffIdenticalFiles(t *testing.T) {
	actual, err := code.GenDiff(fixture("file1.json"), fixture("file1.json"), "plain")
	require.NoError(t, err)
	require.Empty(t, actual)
}

func TestGenDiffErrors(t *testing.T) {
	valid := fixture("file1.json")
	dir := t.TempDir()
	invalid := filepath.Join(dir, "broken.json")
	require.NoError(t, os.WriteFile(invalid, []byte(`{"a":`), 0o600))
	unsupported := filepath.Join(dir, "config.txt")
	require.NoError(t, os.WriteFile(unsupported, []byte(`{"a": 1}`), 0o600))
	tests := []struct{ name, first, second, format, message string }{
		{"missing first", filepath.Join(dir, "missing.json"), valid, "stylish", "read"},
		{"missing second", valid, filepath.Join(dir, "missing.json"), "stylish", "read"},
		{"invalid first", invalid, valid, "stylish", "parse"},
		{"invalid second", valid, invalid, "stylish", "parse"},
		{"unsupported input", unsupported, valid, "stylish", "unsupported"},
		{"unknown output", valid, valid, "xml", "unsupported"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := code.GenDiff(test.first, test.second, test.format)
			require.ErrorContains(t, err, test.message)
			require.Empty(t, actual)
		})
	}
}

func TestGenDiffRepresentsNullValuesInJSON(t *testing.T) {
	result, err := code.GenDiff(fixture("file1.json"), fixture("file2.json"), "json")
	require.NoError(t, err)
	var tree map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &tree))
	require.Equal(t, "root", tree["type"])
	common := tree["children"].([]any)[0].(map[string]any)
	children := common["children"].([]any)
	setting3 := children[3].(map[string]any)
	require.Equal(t, "changed", setting3["type"])
	require.Equal(t, true, setting3["value1"])
	require.NotContains(t, setting3, "value2")
}
