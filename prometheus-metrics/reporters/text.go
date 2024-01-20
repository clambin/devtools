package reporters

import (
	"errors"
	"fmt"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"io"
	"strings"
)

type TextReporter struct {
	out    io.Writer
	labels bool
}

func NewTextReporter(w io.Writer, labels bool) *TextReporter {
	return &TextReporter{out: w, labels: labels}
}

func (r TextReporter) Encode(v any) error {
	m, ok := v.([]metrics.Metric)
	if !ok {
		return errors.New("invalid type passed to Encode")
	}

	for i := range m {
		_, _ = r.out.Write([]byte(fmt.Sprintf("%-40s\t%-8s\t%s", m[i].Name, m[i].Type, m[i].Help)))
		if r.labels {
			_, _ = r.out.Write([]byte("\t" + strings.Join(m[i].Labels, ", ")))
		}
		_, _ = r.out.Write([]byte("\n"))
	}
	return nil
}
