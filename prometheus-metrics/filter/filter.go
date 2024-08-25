package filter

import (
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"regexp"
)

func Filter(filters []string) (func(metric metrics.Metric) bool, error) {
	f, err := newMatcher(filters)
	if err != nil {
		return nil, err
	}

	return f.Match, nil
}

type matcher []*regexp.Regexp

func newMatcher(filters []string) (f matcher, err error) {
	f = make([]*regexp.Regexp, len(filters))
	for i := range filters {
		f[i], err = regexp.Compile(filters[i])
		if err != nil {
			break
		}
	}
	return f, err
}

func (m matcher) Match(in metrics.Metric) bool {
	for j := range m {
		if m[j].MatchString(in.Name) {
			return true
		}
	}
	return len(m) == 0
}
