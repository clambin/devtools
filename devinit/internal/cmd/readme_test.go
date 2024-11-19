package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_writeREADME(t *testing.T) {
	tests := []struct {
		moduleType string
		modInfo
		want string
	}{
		{
			moduleType: "library",
			modInfo: modInfo{
				fullPath:     "github.com/clambin/foo",
				strippedPath: "clambin/foo",
			},
			want: `# foo
[![release](https://img.shields.io/github/v/tag/clambin/foo?color=green&label=release&style=plastic)](https://github.com/clambin/foo/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/foo?style=plastic)](https://app.codecov.io/gh/clambin/foo)
[![test](https://github.com/clambin/foo/workflows/test/badge.svg)](https://github.com/clambin/foo/actions)
[![go report card](https://goreportcard.com/badge/github.com/clambin/foo)](https://goreportcard.com/report/github.com/clambin/foo)
[![godoc](https://pkg.go.dev/badge/github.com/clambin/foo?utm_source=godoc)](https://pkg.go.dev/github.com/clambin/foo)
[![license](https://img.shields.io/github/license/clambin/foo?style=plastic)](LICENSE.md)
`,
		},
		{
			moduleType: "program",
			modInfo: modInfo{
				fullPath:     "github.com/clambin/foo",
				strippedPath: "clambin/foo",
			},
			want: `# foo
[![release](https://img.shields.io/github/v/tag/clambin/foo?color=green&label=release&style=plastic)](https://github.com/clambin/foo/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/foo?style=plastic)](https://app.codecov.io/gh/clambin/foo)
[![test](https://github.com/clambin/foo/workflows/test/badge.svg)](https://github.com/clambin/foo/actions)
[![go report card](https://goreportcard.com/badge/github.com/clambin/foo)](https://goreportcard.com/report/github.com/clambin/foo)
[![license](https://img.shields.io/github/license/clambin/foo?style=plastic)](LICENSE.md)
`,
		},
		{
			moduleType: "container",
			modInfo: modInfo{
				fullPath:     "github.com/clambin/foo",
				strippedPath: "clambin/foo",
			},
			want: `# foo
[![release](https://img.shields.io/github/v/tag/clambin/foo?color=green&label=release&style=plastic)](https://github.com/clambin/foo/releases)
[![codecov](https://img.shields.io/codecov/c/gh/clambin/foo?style=plastic)](https://app.codecov.io/gh/clambin/foo)
[![build](https://github.com/clambin/foo/workflows/build/badge.svg)](https://github.com/clambin/foo/actions)
[![go report card](https://goreportcard.com/badge/github.com/clambin/foo)](https://goreportcard.com/report/github.com/clambin/foo)
[![license](https://img.shields.io/github/license/clambin/foo?style=plastic)](LICENSE.md)
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.moduleType, func(t *testing.T) {
			var out bytes.Buffer
			writeREADME(&out, tt.modInfo, tt.moduleType)
			assert.Equal(t, tt.want, out.String())
		})
	}
}
