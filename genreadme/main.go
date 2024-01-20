package main

import (
	"flag"
	"fmt"
	"github.com/clambin/devtools/genreadme/generate"
	"os"
)

var input = flag.String("input", "go.mod", "go.mod path")

func main() {
	flag.Parse()

	f, err := os.Open(*input)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open go.mod: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() { _ = f.Close() }()

	if err = generate.Write(os.Stdout, f); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse go.mod: %s\n", err.Error())
	}
}
