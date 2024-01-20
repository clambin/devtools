package reporters

import (
	"fmt"
	"io"
	"os"
)

type Reporter interface {
	Encode(any) error
}

var Modes = []string{
	"json",
	"text",
	"markdown",
}

func NewReporter(w io.Writer, mode string, labels bool) (Reporter, error) {
	switch mode {
	case "json":
		if labels {
			_, _ = fmt.Fprintln(os.Stderr, "WARNING: labels ignored in json mode")
		}
		return NewJSONEncoder(w), nil
	case "text":
		return NewTextReporter(w, labels), nil
	case "markdown":
		return NewMarkdownReporter(w, labels), nil
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}
}
