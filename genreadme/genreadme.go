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

var input = flag.String("input", "go.mod", "path of the go.mod file to read")

func main() {
	flag.Parse()
	if err := createREADME(os.Stdout, *input); err != nil {
		panic(err)
	}
}

func createREADME(w io.Writer, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	info, err := getModFile(f)
	if err != nil {
		return fmt.Errorf("modfile: %w", err)
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
	_, _ = w.Write([]byte("# " + filepath.Base(info.fullPath) + "\n"))
}

func writeTag(w io.Writer, info modInfo) {
	writeLink(w,
		"GitHub tag (latest by date)",
		"https://img.shields.io/github/v/tag/"+info.strippedPath+"?color=green&label=Release&style=plastic",
		"https://"+info.fullPath+"/releases",
	)
}

func writeCodeCov(w io.Writer, info modInfo) {
	writeLink(w,
		"Codecov",
		"https://img.shields.io/codecov/c/gh/"+info.strippedPath+"?style=plastic",
		"https://app.codecov.io/gh/"+info.strippedPath,
	)
}

func writeTest(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "Test")
}

func writeBuild(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "Build")
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
		"Go Report Card",
		"https://goreportcard.com/badge/"+info.fullPath,
		"https://goreportcard.com/report/"+info.fullPath,
	)
}

func writeLicense(w io.Writer, info modInfo) {
	writeLink(w,
		"License",
		"https://img.shields.io/github/license/"+info.strippedPath+"?style=plastic",
		"LICENSE.md",
	)
}

func writeGoDoc(w io.Writer, info modInfo) {
	writeLink(w,
		"GoDoc",
		"https://pkg.go.dev/badge/"+info.fullPath+"?utm_source=godoc",
		"https://pkg.go.dev/"+info.fullPath,
	)
}

func writeLink(w io.Writer, label, image, link string) {
	output := "![" + label + "](" + image + ")"
	if link != "" {
		output = "[" + output + "](" + link + ")"
	}
	_, _ = w.Write([]byte(output + "\n"))
}
