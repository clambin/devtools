package cmd

import (
	"fmt"
	"github.com/clambin/devtools/devinit/internal/cmd/container"
	"github.com/clambin/devtools/devinit/internal/cmd/workflows"
	"github.com/spf13/cobra"
)

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
			return createFiles(sources, moduleType, info, output, dryRun)
		},
	}

	sources = map[string][]sourceFiles{
		"library": {
			{
				fs: workflows.FS,
				files: []sourceFile{
					{source: "libtest.yaml", destination: ".github/workflows/test.yaml"},
					{source: "release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"program": {
			{
				fs: workflows.FS,
				files: []sourceFile{
					{source: "test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"container": {
			{
				fs: workflows.FS,
				files: []sourceFile{
					{source: "test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "build.yaml", destination: ".github/workflows/build.yaml"},
					{source: "release-container.yaml", destination: ".github/workflows/release.yaml"},
					{source: "vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
			{
				fs: container.FS,
				files: []sourceFile{
					{source: "Dockerfile", destination: "Dockerfile"},
				},
			},
		},
	}
)
