package cmd

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var FS embed.FS

type Arguments struct {
	Module     Module
	ModuleType string
	Author     string
	License    string
	Year       int
}

type Module struct {
	Path      string
	Name      string
	ShortName string
}

func writeFileFromTemplate(w io.Writer, source string, args Arguments) error {
	body, err := FS.ReadFile(filepath.Join("templates", source))
	if err != nil {
		return err
	}
	tmpl := template.New("devinit")
	tmpl.Delims("<<", ">>")
	if tmpl, err = tmpl.Parse(string(body)); err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	return tmpl.Execute(w, args)
}
