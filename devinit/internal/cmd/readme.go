package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var (
	readmeCmd = &cobra.Command{
		Use:   "readme",
		Short: "create README file",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleType, _ := cmd.Flags().GetString("type")
			output, _ := cmd.Flags().GetString("output")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			goModPath, _ := cmd.Flags().GetString("gomod")
			author, _ := cmd.Flags().GetString("author")

			fmt.Printf("Creating basic README.md file for %s module\n", moduleType)

			f, err := os.Open(goModPath)
			if err != nil {
				return fmt.Errorf("failed to open go.mod: %w", err)
			}
			defer func() { _ = f.Close() }()
			info, err := readGoMod(f)
			if err != nil {
				return fmt.Errorf("failed to parse go.mod: %w", err)
			}

			if !dryRun {
				w, err := os.Create(filepath.Join(output, "README.md"))
				if err != nil {
					return fmt.Errorf("could not create README.md: %w", err)
				}
				if err = writeREADME(w, info, moduleType, author, "MIT"); err != nil {
					return fmt.Errorf("failed to write to README.md: %w", err)
				}
				if err = w.Close(); err != nil {
					return fmt.Errorf("failed to save README.md: %w", err)
				}
			}

			return err
		},
	}
)

func writeREADME(w io.Writer, info Module, modType string, author string, licenseType string) error {
	args := Arguments{
		Module:     info,
		ModuleType: modType,
		Author:     author,
		License:    licenseType,
	}
	return writeFileFromTemplate(w, "README.md.tmpl", args)
}
