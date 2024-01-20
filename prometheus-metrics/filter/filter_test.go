package filter_test

import (
	"github.com/clambin/devtools/prometheus-metrics/filter"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name    string
		in      []metrics.Metric
		filters []string
		want    []metrics.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "no filters",
			in:      []metrics.Metric{{Name: "foo"}, {Name: "bar"}},
			want:    []metrics.Metric{{Name: "foo"}, {Name: "bar"}},
			wantErr: assert.NoError,
		},
		{
			name:    "one filter",
			in:      []metrics.Metric{{Name: "foo"}, {Name: "bar"}},
			filters: []string{"^foo"},
			want:    []metrics.Metric{{Name: "foo"}},
			wantErr: assert.NoError,
		},
		{
			name:    "two filters",
			in:      []metrics.Metric{{Name: "foo"}, {Name: "bar"}},
			filters: []string{"^foo", "^snafu"},
			want:    []metrics.Metric{{Name: "foo"}},
			wantErr: assert.NoError,
		},
		{
			name:    "bad filter",
			in:      []metrics.Metric{{Name: "foo"}, {Name: "bar"}},
			filters: []string{"[foo"},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out, err := filter.Filter(tt.in, tt.filters)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, out)
		})
	}
}
