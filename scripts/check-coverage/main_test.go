package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadCoverage(t *testing.T) {
	actual, err := readCoverage(strings.NewReader("mode: set\ncode/a.go:1.1,2.1 3 1\ncode/a.go:3.1,4.1 1 0\ncode/a.go:1.1,2.1 3 0\n"))
	require.NoError(t, err)
	require.Equal(t, 75.0, actual)
}

func TestReadCoverageErrors(t *testing.T) {
	for _, input := range []string{"", "mode: set\n", "broken", "file not-a-number 1", "file 1 bad", strings.Repeat("x", 70000)} {
		_, err := readCoverage(strings.NewReader(input))
		require.Error(t, err)
	}
}

func TestCheckFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "coverage.out")
	require.NoError(t, os.WriteFile(path, []byte("mode: set\ncode/a.go:1.1,2.1 3 1\ncode/a.go:3.1,4.1 1 0\n"), 0o600))
	require.NoError(t, checkFile(path, 75))
	require.ErrorContains(t, checkFile(path, 80), "below")
	require.Error(t, checkFile(path+"missing", 0))
	require.NoError(t, os.WriteFile(path, []byte("broken"), 0o600))
	require.Error(t, checkFile(path, 0))
}
