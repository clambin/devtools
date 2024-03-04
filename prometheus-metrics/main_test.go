package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Main(t *testing.T) {
	tests := []struct {
		name       string
		metrics    string
		format     string
		filters    string
		wantErr    assert.ErrorAssertionFunc
		wantOutput string
	}{
		{
			name: "markdown",
			metrics: `# HELP foo_total foo
# TYPE foo_total counter
foo_total{bar="200"} 20
# HELP bar_total bar
# TYPE bar_total counter
bar_total 20
`,
			format:  "markdown",
			filters: "bar_total",
			wantErr: assert.NoError,
			wantOutput: `| metric | type | help |
| --- | --- | --- |
| bar_total | COUNTER | bar |
`,
		},
		{
			name:    "no metrics",
			format:  "markdown",
			filters: "bar_total",
			wantErr: assert.Error,
		},
		{
			name:    "invalid format",
			format:  "invalid",
			filters: "bar_total",
			wantErr: assert.Error,
		},
		{
			name: "invalid filter",
			metrics: `# HELP foo_total foo
# TYPE foo_total counter
foo_total{bar="200"} 20
# HELP bar_total bar
# TYPE bar_total counter
bar_total 20
`,
			format:  "markdown",
			filters: "[foo",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tt.metrics == "" {
					http.Error(w, "invalid", http.StatusNotFound)
					return
				}
				_, _ = w.Write([]byte(tt.metrics))
			}))
			defer s.Close()

			var out bytes.Buffer

			*addr = s.URL
			*output = tt.format
			*filters = tt.filters
			tt.wantErr(t, Main(&out))
			assert.Equal(t, tt.wantOutput, out.String())
		})
	}
}
