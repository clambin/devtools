package reporters_test

import (
	"bytes"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONEncoder(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr assert.ErrorAssertionFunc
		want    string
	}{
		{
			name:    "labels",
			input:   testMetrics,
			wantErr: assert.NoError,
			want: `[
  {
    "name": "foo",
    "help": "bar",
    "type": "snafu",
    "labels": [
      "bar",
      "snafu"
    ]
  }
]
`,
		},
		{
			name:    "no labels",
			input:   []metrics.Metric{{Name: "foo", Help: "bar", Type: "snafu"}},
			wantErr: assert.NoError,
			want: `[
  {
    "name": "foo",
    "help": "bar",
    "type": "snafu"
  }
]
`,
		},
		{
			name:    "empty",
			input:   []metrics.Metric{},
			wantErr: assert.NoError,
			want: `[]
`,
		},
		{
			name:    "invalid",
			input:   []int{},
			wantErr: assert.NoError,
			want: `[]
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := reporters.NewJSONEncoder(&buf)
			tt.wantErr(t, enc.Encode(tt.input))
			assert.Equal(t, tt.want, buf.String())
		})
	}

}
