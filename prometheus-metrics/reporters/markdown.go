package reporters

import (
	"errors"
	"fmt"
	"github.com/clambin/devtools/prometheus-metrics/metrics"
	"io"
	"strings"
)

type MarkdownReporter struct {
	out    io.Writer
	labels bool
}

func NewMarkdownReporter(w io.Writer, labels bool) *MarkdownReporter {
	return &MarkdownReporter{out: w, labels: labels}
}

func (r MarkdownReporter) Encode(v any) error {
	m, ok := v.([]metrics.Metric)
	if !ok {
		return errors.New("invalid type passed to Encode")
	}

	if len(m) == 0 {
		return nil
	}

	_, _ = r.out.Write([]byte("| metric | type | "))
	if r.labels {
		_, _ = r.out.Write([]byte(" labels | "))
	}
	_, _ = r.out.Write([]byte("help |\n"))

	_, _ = r.out.Write([]byte("| --- | --- | "))
	if r.labels {
		_, _ = r.out.Write([]byte(" --- | "))
	}
	_, _ = r.out.Write([]byte("--- |\n"))

	for i := range m {
		_, _ = r.out.Write([]byte(fmt.Sprintf("| %s | %s | ", m[i].Name, m[i].Type)))
		if r.labels {
			_, _ = r.out.Write([]byte(strings.Join(m[i].Labels, ", ") + "|"))
		}
		_, _ = r.out.Write([]byte(m[i].Help + " |\n"))
	}

	return nil
}
