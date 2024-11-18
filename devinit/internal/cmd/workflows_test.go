package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_workflows(t *testing.T) {
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

			tt.wantErr(t, createFiles(workflows, tt.mode, tmpDir, false))

			for _, want := range tt.want {
				_, err = os.Stat(filepath.Join(tmpDir, want))
				assert.NoError(t, err)
			}
		})
	}
}
