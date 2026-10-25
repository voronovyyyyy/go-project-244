package parsers_test

import (
	"testing"

	"code/parsers"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	jsonInput := []byte(`{"nested":{"number":50,"empty":null,"list":[1,"two",false]}}`)
	expected := map[string]any{"nested": map[string]any{
		"number": float64(50), "empty": nil, "list": []any{float64(1), "two", false},
	}}
	for _, test := range []struct{ format, input string }{
		{"json", string(jsonInput)},
		{"yml", "nested:\n  number: 50\n  empty: null\n  list: [1, two, false]\n"},
		{"yaml", "nested:\n  number: 50.0\n  empty: ~\n  list: [1, two, false]\n"},
	} {
		t.Run(test.format, func(t *testing.T) {
			data := []byte(test.input)
			original := append([]byte(nil), data...)
			actual, err := parsers.Parse(data, test.format)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
			require.Equal(t, original, data)
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, test := range []struct{ name, format, input string }{
		{"unknown format", "xml", "<a/>"},
		{"broken json", "json", `{"a":`},
		{"json list", "json", `[1, 2]`},
		{"json null", "json", `null`},
		{"empty json", "json", ``},
		{"trailing json", "json", `{} {}`},
		{"broken yaml", "yaml", "a: ["},
		{"yaml list", "yaml", "- one\n- two"},
		{"yaml null", "yaml", "null"},
		{"empty yaml", "yaml", ""},
		{"multiple yaml documents", "yaml", "a: 1\n---\nb: 2"},
		{"broken second yaml document", "yaml", "a: 1\n---\n["},
		{"duplicate yaml keys", "yaml", "a: 1\na: 2"},
		{"non-string nested key", "yaml", "nested:\n  1: value"},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parsers.Parse([]byte(test.input), test.format)
			require.Error(t, err)
			require.Nil(t, actual)
		})
	}
}
