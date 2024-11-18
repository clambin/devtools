package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
)

//go:embed licenses/*
var licenseFiles embed.FS

var (
	licensesCmd = &cobra.Command{
		Use:   "licenses",
		Short: "set up license",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			//TODO: support more licenses
			fmt.Println("Creating license file")
			return createFiles(licenses, "mit", output, dryRun)
		},
	}

	licenses = map[string][]sourceFiles{
		"mit": {
			{
				fs:    licenseFiles,
				files: []sourceFile{{source: "licenses/MIT.md", destination: "LICENSE.md"}},
			},
		},
	}
)
