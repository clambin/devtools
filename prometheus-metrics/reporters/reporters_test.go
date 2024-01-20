package reporters_test

import (
	"bytes"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testMetrics = []metrics.Metric{
	{Name: "foo", Help: "bar", Type: "snafu", Labels: []string{"bar", "snafu"}},
}

func TestNewReporter(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		labels  bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "json",
			mode:    "json",
			labels:  true,
			wantErr: assert.NoError,
		},
		{
			name:    "text (no labels)",
			mode:    "text",
			labels:  false,
			wantErr: assert.NoError,
		},
		{
			name:    "text (labels)",
			mode:    "text",
			labels:  true,
			wantErr: assert.NoError,
		},
		{
			name:    "markdown",
			mode:    "markdown",
			wantErr: assert.NoError,
		},
		{
			name:    "invalid",
			mode:    "invalid",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := reporters.NewReporter(&buf, tt.mode, tt.labels)
			tt.wantErr(t, err)
		})
	}
}
