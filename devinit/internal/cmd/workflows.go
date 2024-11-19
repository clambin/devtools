package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
)

//go:embed workflows/*
var workflowFiles embed.FS

//go:embed container/*
var containerFiles embed.FS

var (
	workFlowsCmd = &cobra.Command{
		Use:   "workflows",
		Short: "set up github workflows",
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
			return createFiles(workflows, moduleType, info, output, dryRun)
		},
	}

	workflows = map[string][]sourceFiles{
		"library": {
			{
				fs: workflowFiles,
				files: []sourceFile{
					{source: "workflows/libtest.yaml", destination: ".github/workflows/test.yaml"},
					{source: "workflows/release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "workflows/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "workflows/dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"program": {
			{
				fs: workflowFiles,
				files: []sourceFile{
					{source: "workflows/test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "workflows/release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "workflows/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "workflows/dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"container": {
			{
				fs: workflowFiles,
				files: []sourceFile{
					{source: "workflows/test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "workflows/build.yaml", destination: ".github/workflows/build.yaml"},
					{source: "workflows/release-container.yaml", destination: ".github/workflows/release.yaml"},
					{source: "workflows/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "workflows/dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
			{
				fs: containerFiles,
				files: []sourceFile{
					{source: "container/Dockerfile", destination: "Dockerfile"},
				},
			},
		},
	}
)
