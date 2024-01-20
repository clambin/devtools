package main

import (
	"flag"
	"fmt"
	"github.com/clambin/devtools/prometheus-metrics/filter"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"github.com/clambin/devtools/prometheus-metrics/reporters"
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

	enc, err := reporters.NewReporter(os.Stdout, *output, *labels)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s (valid modes: %s)\n", err.Error(), strings.Join(reporters.Modes, ", "))
		os.Exit(1)
	}

	m, err := metrics.Scrape(*addr, *labels)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get metrics: %s\n", err.Error())
		os.Exit(2)
	}

	m, err = filter.Filter(m, strings.Split(*filters, ","))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to filter metrics: %s\n", err.Error())
		os.Exit(3)
	}

	slices.SortFunc(m, func(a, b metrics.Metric) int {
		return strings.Compare(a.Name, b.Name)
	})

	_ = enc.Encode(m)
}
