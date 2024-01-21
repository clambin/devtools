package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var input = flag.String("input", "go.mod", "go.mod path")

func main() {
	flag.Parse()
	if err := Main(os.Stdout, *input); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Main(w io.Writer, source string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer func() { _ = f.Close() }()

	info, err := getModFile(f)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod: %w", err)
	}
	writeREADME(w, info)
	return nil
}

type modInfo struct {
	fullPath     string
	strippedPath string
}

func getModFile(source io.Reader) (modInfo, error) {
	mod, err := io.ReadAll(source)
	if err != nil {
		return modInfo{}, err
	}
	file, err := modfile.Parse("go.mod", mod, nil)
	if err != nil {
		return modInfo{}, fmt.Errorf("parse: %w", err)
	}
	if file.Module == nil {
		return modInfo{}, errors.New("invalid go.mod file")
	}
	if !strings.HasPrefix(file.Module.Mod.Path, "github.com/") {
		return modInfo{}, errors.New("only supports github-hosted repos")
	}

	strippedPath := strings.TrimPrefix(file.Module.Mod.Path, "github.com/")
	return modInfo{fullPath: file.Module.Mod.Path, strippedPath: strippedPath}, err
}

func writeREADME(w io.Writer, info modInfo) {
	writeTitle(w, info)
	writeTag(w, info)
	writeCodeCov(w, info)
	writeTest(w, info)
	writeBuild(w, info)
	writeGoReport(w, info)
	writeGoDoc(w, info)
	writeLicense(w, info)
}

func writeTitle(w io.Writer, info modInfo) {
	writeLine(w, "# "+filepath.Base(info.fullPath))
}

func writeTag(w io.Writer, info modInfo) {
	writeLink(w,
		"release",
		"https://img.shields.io/github/v/tag/"+info.strippedPath+"?color=green&label=release&style=plastic",
		"https://"+info.fullPath+"/releases",
	)
}

func writeCodeCov(w io.Writer, info modInfo) {
	writeLink(w,
		"codecov",
		"https://img.shields.io/codecov/c/gh/"+info.strippedPath+"?style=plastic",
		"https://app.codecov.io/gh/"+info.strippedPath,
	)
}

func writeTest(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "test")
}

func writeBuild(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "build")
}

func writeWorkFlowResult(w io.Writer, info modInfo, action string) {
	writeLink(w,
		action,
		"https://"+info.fullPath+"/workflows/"+action+"/badge.svg",
		"https://"+info.fullPath+"/actions",
	)
}

func writeGoReport(w io.Writer, info modInfo) {
	writeLink(w,
		"go report card",
		"https://goreportcard.com/badge/"+info.fullPath,
		"https://goreportcard.com/report/"+info.fullPath,
	)
}

func writeLicense(w io.Writer, info modInfo) {
	writeLink(w,
		"license",
		"https://img.shields.io/github/license/"+info.strippedPath+"?style=plastic",
		"LICENSE.md",
	)
}

func writeGoDoc(w io.Writer, info modInfo) {
	writeLink(w,
		"godoc",
		"https://pkg.go.dev/badge/"+info.fullPath+"?utm_source=godoc",
		"https://pkg.go.dev/"+info.fullPath,
	)
}

func writeLink(w io.Writer, label, image, link string) {
	output := "![" + label + "](" + image + ")"
	if link != "" {
		output = "[" + output + "](" + link + ")"
	}
	writeLine(w, output)
}

func writeLine(w io.Writer, line string) {
	_, _ = w.Write([]byte(line + "\n"))
}
