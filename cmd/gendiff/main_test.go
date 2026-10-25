package main

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelp(t *testing.T) {
	for _, flag := range []string{"--help", "-h"} {
		t.Run(flag, func(t *testing.T) {
			var output, errorOutput bytes.Buffer
			err := newCommand(&output, &errorOutput).Run(context.Background(), []string{"gendiff", flag})
			require.NoError(t, err)
			require.Contains(t, output.String(), "Compares two configuration files and shows a difference.")
			require.Contains(t, output.String(), "--format")
			require.Contains(t, output.String(), `default: "stylish"`)
			require.Empty(t, errorOutput.String())
		})
	}
}

func TestCLI(t *testing.T) {
	first := filepath.Join("..", "..", "testdata", "fixture", "file1.json")
	second := filepath.Join("..", "..", "testdata", "fixture", "file2.yml")
	for _, arguments := range [][]string{
		{"gendiff", "--format", "plain", first, second},
		{"gendiff", first, second, "-f", "plain"},
	} {
		var output, errorOutput bytes.Buffer
		err := newCommand(&output, &errorOutput).Run(context.Background(), arguments)
		require.NoError(t, err)
		require.Contains(t, output.String(), "Property 'common.follow' was added with value: false\n")
		require.Empty(t, errorOutput.String())
	}
}

func TestCLIErrors(t *testing.T) {
	for _, arguments := range [][]string{
		{"gendiff"}, {"gendiff", "one"}, {"gendiff", "one", "two", "three"},
		{"gendiff", "missing.json", "missing.yml"}, {"gendiff", "--unknown"},
	} {
		var output, errorOutput bytes.Buffer
		err := newCommand(&output, &errorOutput).Run(context.Background(), arguments)
		require.Error(t, err)
	}
}

func TestCLIEmptyDiff(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "fixture", "file1.json")
	var output bytes.Buffer
	err := newCommand(&output, &output).Run(context.Background(), []string{"gendiff", path, path, "-f", "plain"})
	require.NoError(t, err)
	require.Empty(t, output.String())
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("output unavailable") }

func TestCLIOutputError(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "fixture", "file1.json")
	err := newCommand(failingWriter{}, failingWriter{}).Run(context.Background(), []string{"gendiff", path, path})
	require.ErrorContains(t, err, "output unavailable")
}
