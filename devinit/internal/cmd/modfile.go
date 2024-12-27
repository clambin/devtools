package cmd

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readGoModFile(path string) (Module, error) {
	f, err := os.Open(path)
	if err != nil {
		return Module{}, err
	}
	defer func() { _ = f.Close() }()
	return readGoMod(f)
}

func readGoMod(source io.Reader) (Module, error) {
	mod, err := io.ReadAll(source)
	if err != nil {
		return Module{}, err
	}
	file, err := modfile.Parse("go.mod", mod, nil)
	if err != nil {
		return Module{}, fmt.Errorf("parse: %w", err)
	}
	if file.Module == nil {
		return Module{}, errors.New("invalid go.mod file")
	}

	parts := strings.Split(file.Module.Mod.Path, "/")
	if len(parts) < 2 {
		return Module{}, errors.New("does not look like a public module")
	}

	strippedPath := strings.Join(parts[len(parts)-2:], "/")
	return Module{
		Path:      file.Module.Mod.Path,
		Name:      strippedPath,
		ShortName: filepath.Base(strippedPath),
	}, nil
}
