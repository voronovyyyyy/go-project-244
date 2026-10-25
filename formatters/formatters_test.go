package formatters_test

import (
	"testing"

	"code/formatters"
	"code/internal/diff"

	"github.com/stretchr/testify/require"
)

func TestEmptyDiff(t *testing.T) {
	for _, test := range []struct{ format, expected string }{
		{"stylish", "{\n}"}, {"", "{\n}"}, {"plain", ""},
		{"json", "{\n  \"key\": \"\",\n  \"type\": \"root\",\n  \"children\": []\n}"},
	} {
		t.Run(test.format, func(t *testing.T) {
			actual, err := formatters.Format([]diff.Node{}, test.format)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestPlainValues(t *testing.T) {
	tree := []diff.Node{
		{Key: "array", Type: diff.Added, After: []any{1.0, "two"}},
		{Key: "number", Type: diff.Updated, Before: 0.0, After: 1.25},
		{Key: "object", Type: diff.Added, After: map[string]any{}},
		{Key: "removed", Type: diff.Removed, Before: nil},
		{Key: "same", Type: diff.Unchanged, Before: true, After: true},
		{Key: "text", Type: diff.Added, After: "it's fine"},
	}
	actual, err := formatters.Format(tree, "plain")
	require.NoError(t, err)
	require.Equal(t, "Property 'array' was added with value: [complex value]\n"+
		"Property 'number' was updated. From 0 to 1.25\n"+
		"Property 'object' was added with value: [complex value]\n"+
		"Property 'removed' was removed\n"+
		"Property 'text' was added with value: 'it's fine'", actual)
}

func TestStylishValues(t *testing.T) {
	tree := []diff.Node{
		{Key: "array", Type: diff.Added, After: []any{1.0, "two", nil}},
		{Key: "empty", Type: diff.Added, After: map[string]any{}},
		{Key: "nested", Type: diff.Nested, Children: []diff.Node{}},
	}
	actual, err := formatters.Format(tree, "stylish")
	require.NoError(t, err)
	require.Equal(t, "{\n  + array: [1, two, null]\n  + empty: {}\n    nested: {\n    }\n}", actual)
}

func TestUnsupportedFormat(t *testing.T) {
	result, err := formatters.Format(nil, "xml")
	require.ErrorContains(t, err, "unsupported output format")
	require.Empty(t, result)
}

func TestJSONError(t *testing.T) {
	_, err := formatters.Format([]diff.Node{{Key: "invalid", Type: diff.Added, After: make(chan int)}}, "json")
	require.Error(t, err)
}
