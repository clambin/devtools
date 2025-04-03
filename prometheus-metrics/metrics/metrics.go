package metrics

import (
	"codeberg.org/clambin/go-common/set"
	"errors"
	"fmt"
	prometheus "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"iter"
	"net/http"
)

type Metric struct {
	Name   string   `json:"name"`
	Help   string   `json:"help"`
	Type   string   `json:"type"`
	Labels []string `json:"labels,omitempty"`
}

func Scrape(target string, labels bool, shouldExport func(metric Metric) bool) ([]Metric, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get: %s", resp.Status)
	}
	var result []Metric
	for metric, err := range metrics(resp.Body, labels) {
		if err != nil {
			return nil, err
		}
		if shouldExport(metric) {
			result = append(result, metric)
		}
	}
	return result, nil
}

func metrics(r io.Reader, withLabels bool) iter.Seq2[Metric, error] {
	return func(yield func(Metric, error) bool) {
		dec := expfmt.NewDecoder(r, expfmt.TextVersion)
		for {
			var m prometheus.MetricFamily
			if err := dec.Decode(&m); err != nil {
				if !errors.Is(err, io.EOF) {
					yield(Metric{}, err)
				}
				return
			}

			var labels []string
			if withLabels {
				labels = getLabels(&m)
			}

			metric := Metric{
				Name:   m.GetName(),
				Help:   m.GetHelp(),
				Type:   m.GetType().String(),
				Labels: labels,
			}
			if !yield(metric, nil) {
				return
			}
		}
	}
}

func getLabels(m *prometheus.MetricFamily) []string {
	result := set.New[string]()

	for _, m := range m.GetMetric() {
		for _, l := range m.GetLabel() {
			result.Add(l.GetName())
		}
	}

	return result.ListOrdered()
}
