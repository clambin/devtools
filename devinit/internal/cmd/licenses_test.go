package cmd

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_licenses(t *testing.T) {
	var out bytes.Buffer
	err := writeLicense(&out, "Christophe Lambin", "MIT")
	require.NoError(t, err)

	gp := filepath.Join("testdata", strings.ToLower(t.Name())+".golden")
	if *update {
		require.NoError(t, os.WriteFile(gp, out.Bytes(), 0644))
	}

	expected, err := os.ReadFile(gp)
	require.NoError(t, err)
	require.Equal(t, string(expected), out.String())
}
