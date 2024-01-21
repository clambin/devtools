package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		mod     string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			mod:  "module github.com/clambin/foo\n",
			want: `# foo
[![Release](https://img.shields.io/github/v/tag/clambin/foo?color=green&label=Release&style=plastic)](https://github.com/clambin/foo/releases)
[![Codecov](https://img.shields.io/codecov/c/gh/clambin/foo?style=plastic)](https://app.codecov.io/gh/clambin/foo)
[![Test](https://github.com/clambin/foo/workflows/Test/badge.svg)](https://github.com/clambin/foo/actions)
[![Build](https://github.com/clambin/foo/workflows/Build/badge.svg)](https://github.com/clambin/foo/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/clambin/foo)](https://goreportcard.com/report/github.com/clambin/foo)
[![GoDoc](https://pkg.go.dev/badge/github.com/clambin/foo?utm_source=godoc)](https://pkg.go.dev/github.com/clambin/foo)
[![License](https://img.shields.io/github/license/clambin/foo?style=plastic)](LICENSE.md)
`,
			wantErr: assert.NoError,
		},
		{
			name:    "invalid",
			mod:     "not a valid go.mod file",
			wantErr: assert.Error,
		},
		{
			name:    "missing",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var filename string
			if tt.mod != "" {
				filename = mkTemp(t, tt.mod)
				defer func() { _ = os.Remove(filename) }()
			}

			var out bytes.Buffer
			tt.wantErr(t, Main(&out, filename))
			assert.Equal(t, tt.want, out.String())
		})
	}
}

func mkTemp(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp("", "")
	require.NoError(t, err)
	_, err = f.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	return f.Name()
}

func Test_getModFile(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    modInfo
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "pass",
			input: `
module github.com/clambin/foo
`,
			want: modInfo{
				fullPath:     "github.com/clambin/foo",
				strippedPath: "clambin/foo",
			},
			wantErr: assert.NoError,
		},
		{
			name: "non-github",
			input: `
module foo
`,
			wantErr: assert.Error,
		},
		{
			name:    "invalid",
			input:   `invalid`,
			wantErr: assert.Error,
		},
		{
			name:    "empty",
			wantErr: assert.Error,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod, err := getModFile(bytes.NewBufferString(tt.input))
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, mod)
			}
		})
	}
}
