package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

var (
	licensesCmd = &cobra.Command{
		Use:   "licenses",
		Short: "create license file",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			author, _ := cmd.Flags().GetString("author")

			fmt.Println("Creating license file")

			if dryRun {
				return nil
			}
			w, err := os.Create(filepath.Join(output, "LICENSE.md"))
			if err != nil {
				return fmt.Errorf("could not create LICENSE file: %w", err)
			}
			defer func() { _ = w.Close() }()
			//TODO: support more licenses
			return writeLicense(w, author, "MIT")
		},
	}
)

func writeLicense(w io.Writer, author string, _ string) error {
	args := Arguments{
		Author: author,
	}
	return writeFileFromTemplate(w, "MIT.md.tmpl", args)
}
