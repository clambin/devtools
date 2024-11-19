package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_licenses(t *testing.T) {
	tests := []struct {
		mode    string
		wantErr assert.ErrorAssertionFunc
		want    []string
	}{
		{
			mode:    "mit",
			wantErr: assert.NoError,
			want:    []string{"LICENSE.md"},
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

			tt.wantErr(t, createFiles(licenses, tt.mode, modInfo{}, tmpDir, false))

			for _, want := range tt.want {
				_, err = os.Stat(filepath.Join(tmpDir, want))
				assert.NoError(t, err)
			}
		})
	}
}
