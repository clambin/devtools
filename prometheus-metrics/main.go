package main

import (
	"flag"
	"fmt"
	"github.com/clambin/devtools/prometheus-metrics/filter"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
	"io"
	"os"
	"slices"
	"strings"
)

var (
	addr    = flag.String("addr", "http://localhost:9091/metrics", "Prometheus metrics URL")
	labels  = flag.Bool("labels", false, "include labels")
	filters = flag.String("filter", "", "comma-separated regexp filters for metric names")
	output  = flag.String("output", "text", "output mode ("+strings.Join(reporters.Modes, ", ")+")")
)

func main() {
	flag.Parse()
	if err := Main(os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Main(w io.Writer) error {
	enc, err := reporters.NewReporter(w, *output, *labels)
	if err != nil {
		return fmt.Errorf("%w (valid modes: %s)", err, strings.Join(reporters.Modes, ", "))
	}

	m, err := metrics.Scrape(*addr, *labels)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	m, err = filter.Filter(m, strings.Split(*filters, ","))
	if err != nil {
		return fmt.Errorf("failed to filter metrics: %w", err)
	}

	slices.SortFunc(m, func(a, b metrics.Metric) int {
		return strings.Compare(a.Name, b.Name)
	})

	return enc.Encode(m)
}
