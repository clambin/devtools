package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_writeWorkflows(t *testing.T) {
	tests := []struct {
		mode    string
		wantErr assert.ErrorAssertionFunc
		want    []string
	}{
		{
			mode:    "library",
			wantErr: assert.NoError,
			want:    []string{".github/dependabot.yaml", ".github/workflows/release.yaml", ".github/workflows/test.yaml", ".github/workflows/vulnerabilities.yaml"},
		},
		{
			mode:    "program",
			wantErr: assert.NoError,
			want:    []string{".github/dependabot.yaml", ".github/workflows/release.yaml", ".github/workflows/test.yaml", ".github/workflows/vulnerabilities.yaml"},
		},
		{
			mode:    "container",
			wantErr: assert.NoError,
			want:    []string{".github/dependabot.yaml", ".github/workflows/build.yaml", ".github/workflows/release.yaml", ".github/workflows/test.yaml", ".github/workflows/vulnerabilities.yaml", "Dockerfile"},
		},
		{
			mode:    "invalid",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "")
			require.NoError(t, err)
			t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

			tt.wantErr(t, writeWorkflows(tmpDir, tt.mode, Module{Path: "example.com/foo/bar", Name: "foo/bar", ShortName: "bar"}))

			for _, want := range tt.want {
				got, err := os.ReadFile(filepath.Join(tmpDir, want))
				assert.NoError(t, err)

				gp := filepath.Join("testdata", strings.ToLower(t.Name()), filepath.Base(want)+".golden")

				if *update {
					require.NoError(t, os.WriteFile(gp, got, 0o644))
				}

				wantBody, err := os.ReadFile(gp)
				require.NoError(t, err)
				assert.Equal(t, string(wantBody), string(got))
			}
		})
	}
}
