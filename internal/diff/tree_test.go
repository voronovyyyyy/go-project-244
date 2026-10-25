package diff_test

import (
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	first := map[string]any{
		"removed": map[string]any{"child": true},
		"nested":  map[string]any{"old": nil, "same": 1.0},
		"changed": map[string]any{"child": true},
		"null":    nil, "list": []any{1.0, nil}, "same": false,
	}
	second := map[string]any{
		"added":   map[string]any{"child": true},
		"nested":  map[string]any{"new": nil, "same": 1.0},
		"changed": "scalar", "list": []any{1.0, nil}, "same": false,
	}
	tree := diff.Build(first, second)
	require.Equal(t, []diff.Node{
		{Key: "added", Type: diff.Added, After: second["added"]},
		{Key: "changed", Type: diff.Updated, Before: first["changed"], After: "scalar"},
		{Key: "list", Type: diff.Unchanged, Before: first["list"]},
		{Key: "nested", Type: diff.Nested, Children: []diff.Node{
			{Key: "new", Type: diff.Added},
			{Key: "old", Type: diff.Removed},
			{Key: "same", Type: diff.Unchanged, Before: 1.0},
		}},
		{Key: "null", Type: diff.Removed},
		{Key: "removed", Type: diff.Removed, Before: first["removed"]},
		{Key: "same", Type: diff.Unchanged, Before: false},
	}, tree)
	require.Contains(t, first, "removed")
	require.NotContains(t, first, "added")
	require.Contains(t, second, "added")
	require.NotContains(t, second, "removed")
}

func TestBuildEmptyObjects(t *testing.T) {
	require.Equal(t, []diff.Node{}, diff.Build(map[string]any{}, map[string]any{}))
	tree := diff.Build(map[string]any{"object": map[string]any{}}, map[string]any{"object": map[string]any{}})
	require.Len(t, tree, 1)
	require.Equal(t, diff.Nested, tree[0].Type)
	require.Empty(t, tree[0].Children)
}
