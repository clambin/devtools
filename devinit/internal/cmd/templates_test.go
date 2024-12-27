package cmd

import (
	"bytes"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func Test_writeFileFromTemplate(t *testing.T) {
	tests := []struct {
		name   string
		source string
		args   Arguments
	}{
		{
			name:   "readme",
			source: "README.md.tmpl",
			args: Arguments{
				Module:  Module{"example.com/foo/bar", "foo/bar", "bar"},
				Author:  "Christophe Lambin",
				License: "MIT",
			},
		},
		{
			name:   "license",
			source: "MIT.md.tmpl",
			args: Arguments{
				Author: "Christophe Lambin",
			},
		},
		{
			name:   "Dockerfile",
			source: "Dockerfile.tmpl",
			args: Arguments{
				Module: Module{"example.com/foo/bar", "foo/bar", "bar"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := writeFileFromTemplate(&output, tt.source, tt.args)
			require.NoError(t, err)

			gp := filepath.Join("testdata", strings.ToLower(t.Name())+".golden")
			if *update {
				require.NoError(t, os.WriteFile(gp, output.Bytes(), 0o644))
			}

			expected, err := os.ReadFile(gp)
			require.NoError(t, err)

			assert.Equal(t, string(expected), output.String())
		})
	}
}
