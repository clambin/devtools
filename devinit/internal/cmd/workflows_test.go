package cmd

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

var debug = flag.Bool("update", false, "update golden files")

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

			tt.wantErr(t, createFiles(workflows, tt.mode, modInfo{}, tmpDir, false))

			for _, want := range tt.want {
				_, err = os.Stat(filepath.Join(tmpDir, want))
				assert.NoError(t, err)
			}
		})
	}
}

func Test_workflows_templates(t *testing.T) {
	in := make(targetFiles, 2)
	var err error
	in["dockerfile"], err = containerFiles.ReadFile("container/Dockerfile")
	require.NoError(t, err)
	in["buildfile"], err = workflowFiles.ReadFile("workflows/build.yaml")
	require.NoError(t, err)

	info := modInfo{
		fullPath:     "example.com/foo/bar",
		strippedPath: "foo/bar",
	}
	out, err := templateFiles(in, info)
	require.NoError(t, err)

	for file, body := range out {
		gp := filepath.Join("testdata", t.Name()+"-"+file+".golden")
		if *debug {
			require.NoError(t, os.WriteFile(gp, []byte(body), 0o644))
		}
		golden, err := os.ReadFile(gp)
		require.NoError(t, err)
		assert.Equal(t, string(golden), string(body))
	}
}
