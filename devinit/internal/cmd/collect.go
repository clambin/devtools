package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
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

func createFiles(sources map[string][]sourceFiles, mode string, output string, dryRun bool) error {
	source, ok := sources[mode]
	if !ok {
		return fmt.Errorf("unknown output type: %s", mode)
	}

	collected, err := collectFiles(source)
	if err == nil && !dryRun {
		err = writeFiles(output, collected)
	}
	return err

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
