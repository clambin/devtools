package cmd

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"strings"
)

type modInfo struct {
	fullPath     string
	strippedPath string
}

func readGoModFile(path string) (modInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return modInfo{}, err
	}
	defer func() { _ = f.Close() }()
	return readGoMod(f)
}

func readGoMod(source io.Reader) (modInfo, error) {
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

	parts := strings.Split(file.Module.Mod.Path, "/")
	if len(parts) < 2 {
		return modInfo{}, errors.New("does not look like a public module")
	}

	strippedPath := strings.Join(parts[len(parts)-2:], "/")
	return modInfo{fullPath: file.Module.Mod.Path, strippedPath: strippedPath}, err
}
