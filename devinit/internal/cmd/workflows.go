package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type sourceFile struct {
	source      string
	destination string
	isTemplate  bool
}

var (
	workFlowsCmd = &cobra.Command{
		Use:   "workflows",
		Short: "create github workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleType, _ := cmd.Flags().GetString("type")
			output, _ := cmd.Flags().GetString("output")
			gomod := cmd.Flag("gomod").Value.String()
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			info, err := readGoModFile(gomod)
			if err != nil {
				return fmt.Errorf("invalid go.mod: %w", err)
			}

			fmt.Printf("Creating GitHub workflows for %s module\n", moduleType)
			if !dryRun {
				err = writeWorkflows(output, moduleType, info)
			}
			return err
		},
	}

	sources = map[string][]sourceFile{
		"library": {
			{"workflows/libtest.yaml", ".github/workflows/test.yaml", false},
			{"workflows/release.yaml", ".github/workflows/release.yaml", false},
			{"workflows/vulnerabilities.yaml", ".github/workflows/vulnerabilities.yaml", false},
			{"dependabot.yaml", ".github/dependabot.yaml", false},
		},
		"program": {
			{"workflows/test.yaml", ".github/workflows/test.yaml", false},
			{"workflows/release.yaml", ".github/workflows/release.yaml", false},
			{"workflows/vulnerabilities.yaml", ".github/workflows/vulnerabilities.yaml", false},
			{"dependabot.yaml", ".github/dependabot.yaml", false},
		},
		"container": {
			{"workflows/test.yaml", ".github/workflows/test.yaml", false},
			{"workflows/build.yaml.tmpl", ".github/workflows/build.yaml", true},
			{"workflows/release-container.yaml.tmpl", ".github/workflows/release.yaml", true},
			{"workflows/vulnerabilities.yaml", ".github/workflows/vulnerabilities.yaml", false},
			{"dependabot.yaml", ".github/dependabot.yaml", false},
			{"Dockerfile.tmpl", "Dockerfile", true},
		},
	}
)

func writeWorkflows(output string, moduleType string, modInfo Module) error {
	toWrite, ok := sources[moduleType]
	if !ok {
		return fmt.Errorf("invalid module type: %s", moduleType)
	}
	args := Arguments{
		Module: modInfo,
	}
	for _, source := range toWrite {
		directory := filepath.Dir(source.destination)
		err := os.MkdirAll(filepath.Join(output, directory), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %q: %w", directory, err)
		}
		f, err := os.OpenFile(filepath.Join(output, source.destination), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("could not open file %q: %w", source.destination, err)
		}
		if err = writeFileFromTemplate(f, source.source, args); err != nil {
			return fmt.Errorf("could not write target file %q: %w", source.destination, err)
		}
		if err = f.Close(); err != nil {
			return fmt.Errorf("failed to close file %q: %w", source.destination, err)
		}
	}
	return nil
}
