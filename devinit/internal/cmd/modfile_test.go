package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getModFile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Module
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "pass",
			input: `
module github.com/clambin/foo
`,
			want: Module{
				Path:      "github.com/clambin/foo",
				Name:      "clambin/foo",
				ShortName: "foo",
			},
			wantErr: assert.NoError,
		},
		{
			name: "other source",
			input: `
module example.com/clambin/foo
`,
			want: Module{
				Path:      "example.com/clambin/foo",
				Name:      "clambin/foo",
				ShortName: "foo",
			},
			wantErr: assert.NoError,
		},
		{
			name: "non-public module",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod, err := readGoMod(bytes.NewBufferString(tt.input))
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, mod)
		})
	}
}
