package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_writeREADME(t *testing.T) {
	tests := []struct {
		moduleType string
		modInfo    Module
	}{
		{"library", Module{"github.com/clambin/foo", "clambin/foo", "foo"}},
		{"program", Module{"github.com/clambin/foo", "clambin/foo", "foo"}},
		{"container", Module{"github.com/clambin/foo", "clambin/foo", "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.moduleType, func(t *testing.T) {
			var out bytes.Buffer
			err := writeREADME(&out, tt.modInfo, tt.moduleType, "MY NAME", "MIT")
			require.NoError(t, err)

			gp := filepath.Join("testdata", strings.ToLower(t.Name())+".golden")
			if *update {
				require.NoError(t, os.WriteFile(gp, out.Bytes(), 0644))
			}

			expected, err := os.ReadFile(gp)
			require.NoError(t, err)
			assert.Equal(t, string(expected), out.String())
		})
	}
}
