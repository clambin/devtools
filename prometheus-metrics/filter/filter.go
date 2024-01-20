package filter

import (
	"fmt"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"regexp"
)

func Filter(in []metrics.Metric, filters []string) ([]metrics.Metric, error) {
	out := make([]metrics.Metric, 0, len(in))
	re := make([]*regexp.Regexp, len(filters))

	for i := range in {
		var match bool
		for j := range filters {
			if re[j] == nil {
				var err error
				if re[j], err = regexp.Compile(filters[i]); err != nil {
					return nil, fmt.Errorf("bad regexp %s: %w", filters[i], err)
				}
			}

			if re[j].MatchString(in[i].Name) {
				match = true
				break
			}
		}

		if len(filters) == 0 || match {
			out = append(out, in[i])
		}
	}

	return out, nil
}
