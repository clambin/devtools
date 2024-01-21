package reporters_test

import (
	"bytes"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextReporter(t *testing.T) {
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
			want:    "foo                                     \tsnafu   \tbar\n",
		},
		{
			name:    "labels",
			input:   testMetrics,
			labels:  true,
			wantErr: assert.NoError,
			want:    "foo                                     \tsnafu   \tbar\tbar, snafu\n",
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
			t.Parallel()

			var buf bytes.Buffer
			enc := reporters.NewTextReporter(&buf, tt.labels)
			tt.wantErr(t, enc.Encode(tt.input))
			assert.Equal(t, tt.want, buf.String())
		})
	}
}
