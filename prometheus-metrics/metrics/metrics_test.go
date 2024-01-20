package metrics

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestScrape(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		labels  bool
		want    []Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			in:      ``,
			wantErr: assert.NoError,
		},
		{
			name: "single metric",
			in: `
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.221848e+06
`,
			want:    []Metric{{Name: "go_memstats_alloc_bytes", Help: "Number of bytes allocated and still in use.", Type: "GAUGE"}},
			wantErr: assert.NoError,
		},
		{
			name: "single metric with labels",
			in: `
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 20
`,
			labels:  true,
			want:    []Metric{{Name: "promhttp_metric_handler_requests_total", Help: "Total number of scrapes by HTTP status code.", Type: "COUNTER", Labels: []string{"code"}}},
			wantErr: assert.NoError,
		},
		{
			name: "single metric without labels",
			in: `
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 20
`,
			labels:  false,
			want:    []Metric{{Name: "promhttp_metric_handler_requests_total", Help: "Total number of scrapes by HTTP status code.", Type: "COUNTER"}},
			wantErr: assert.NoError,
		},
		{
			name: "multiple metrics",
			in: `
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 20
# HELP foo_total foo
# TYPE foo_total counter
foo_total{bar="200"} 20
# HELP bar_total bar
# TYPE bar_total counter
bar_total 20
`,
			labels: true,
			want: []Metric{
				{Name: "bar_total", Help: "bar", Type: "COUNTER"},
				{Name: "foo_total", Help: "foo", Labels: []string{"bar"}, Type: "COUNTER"},
				{Name: "promhttp_metric_handler_requests_total", Help: "Total number of scrapes by HTTP status code.", Labels: []string{"code"}, Type: "COUNTER"},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid metric",
			in:      `<not a valid metric>`,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(tt.in))
			}))
			defer server.Close()

			m, err := Scrape(server.URL, tt.labels)

			slices.SortFunc(m, func(a, b Metric) int { return strings.Compare(a.Name, b.Name) })

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, m)
		})
	}
}
