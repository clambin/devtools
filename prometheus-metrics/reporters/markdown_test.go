package reporters_test

import (
	"bytes"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkdownReporter(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		labels  bool
		wantErr assert.ErrorAssertionFunc
		want    string
	}{
		{
			name:    "no labels",
			input:   testMetrics,
			labels:  false,
			wantErr: assert.NoError,
			want: `| metric | type | help |
| --- | --- | --- |
| foo | snafu | bar |
`,
		},
		{
			name:    "labels",
			input:   testMetrics,
			labels:  true,
			wantErr: assert.NoError,
			want: `| metric | type |  labels | help |
| --- | --- |  --- | --- |
| foo | snafu | bar, snafu|bar |
`,
		},
		{
			name:    "empty",
			input:   []metrics.Metric{},
			labels:  false,
			wantErr: assert.NoError,
			want:    "",
		},
		{
			name:    "invalid",
			input:   []int{},
			labels:  false,
			wantErr: assert.Error,
			want:    "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := reporters.NewMarkdownReporter(&buf, tt.labels)
			tt.wantErr(t, enc.Encode(tt.input))
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
