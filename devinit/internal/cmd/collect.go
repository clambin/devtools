package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type sourceFiles struct {
	fs    embed.FS
	files []sourceFile
}

type sourceFile struct {
	source      string
	destination string
}

type targetFiles map[string][]byte

func createFiles(sources map[string][]sourceFiles, mode string, info modInfo, output string, dryRun bool) error {
	source, ok := sources[mode]
	if !ok {
		return fmt.Errorf("unknown output type: %s", mode)
	}

	collected, err := collectFiles(source)
	if err != nil {
		return fmt.Errorf("collecting files: %w", err)
	}
	if collected, err = templateFiles(collected, info); err != nil {
		return fmt.Errorf("processing files: %w", err)
	}
	if dryRun {
		return nil
	}
	return writeFiles(output, collected)
}

func collectFiles(sources []sourceFiles) (targetFiles, error) {
	collected := make(map[string][]byte)
	for _, source := range sources {
		for _, file := range source.files {
			var err error
			if collected[file.destination], err = source.fs.ReadFile(file.source); err != nil {
				return nil, fmt.Errorf("collect %s: %w", file.source, err)
			}
		}
	}
	return collected, nil
}

func templateFiles(sources targetFiles, info modInfo) (targetFiles, error) {
	collected := make(targetFiles, len(sources))
	for name, body := range sources {
		tmpl := template.New("devinit")
		tmpl.Delims("<<", ">>")
		tmpl, err := tmpl.Parse(string(body))
		if err != nil {
			return nil, fmt.Errorf("parse template: %w", err)
		}

		args := struct {
			App   string
			Image string
		}{
			App: filepath.Base(info.strippedPath),
			// TODO: support different registries?
			Image: "ghcr.io/" + info.strippedPath,
		}

		var out bytes.Buffer
		err = tmpl.Execute(&out, args)

		collected[name] = out.Bytes()
	}
	return collected, nil
}

func writeFiles(baseDir string, collected targetFiles) error {
	for filename, content := range collected {
		dirname := filepath.Join(baseDir, filepath.Dir(filename))
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dirname, err)
		}
		file := filepath.Join(baseDir, filename)
		if err := os.WriteFile(file, content, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", file, err)
		}
	}
	return nil
}
