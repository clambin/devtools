package metrics

import (
	"errors"
	"fmt"
	"github.com/clambin/go-common/set"
	prometheus "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
	"net/http"
)

type Metric struct {
	Name   string   `json:"name"`
	Help   string   `json:"help"`
	Type   string   `json:"type"`
	Labels []string `json:"labels,omitempty"`
}

func Scrape(target string, labels bool) ([]Metric, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get: %s", resp.Status)
	}
	return GetMetrics(resp.Body, labels)
}

func GetMetrics(r io.Reader, withLabels bool) ([]Metric, error) {
	dec := expfmt.NewDecoder(r, expfmt.TextVersion)

	var metrics []Metric
	var err error

	for {
		var m prometheus.MetricFamily
		if err = dec.Decode(&m); err != nil {
			break
		}

		var labels []string
		if withLabels {
			labels = getLabels(&m)
		}

		metrics = append(metrics, Metric{
			Name:   m.GetName(),
			Help:   m.GetHelp(),
			Type:   m.GetType().String(),
			Labels: labels,
		})
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	return metrics, err
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
