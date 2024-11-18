package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
)

//go:embed actions/*
var actionFiles embed.FS

//go:embed container/*
var containerFiles embed.FS

var (
	workFlowsCmd = &cobra.Command{
		Use:   "workflows",
		Short: "set up github workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleType, _ := cmd.Flags().GetString("type")
			output, _ := cmd.Flags().GetString("output")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			fmt.Printf("Creating GitHub workflows for %s module\n", moduleType)
			return createFiles(workflows, moduleType, output, dryRun)
		},
	}

	workflows = map[string][]sourceFiles{
		"library": {
			{
				fs: actionFiles,
				files: []sourceFile{
					{source: "actions/libtest.yaml", destination: ".github/workflows/test.yaml"},
					{source: "actions/release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "actions/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "actions/dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"program": {
			{
				fs: actionFiles,
				files: []sourceFile{
					{source: "actions/test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "actions/release.yaml", destination: ".github/workflows/release.yaml"},
					{source: "actions/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "actions/dependabot.yaml", destination: ".github/dependabot.yaml"},
				},
			},
		},
		"container": {
			{
				fs: actionFiles,
				files: []sourceFile{
					{source: "actions/test.yaml", destination: ".github/workflows/test.yaml"},
					{source: "actions/build.yaml", destination: ".github/workflows/build.yaml"},
					{source: "actions/release-container.yaml", destination: ".github/workflows/release.yaml"},
					{source: "actions/vulnerabilities.yaml", destination: ".github/workflows/vulnerabilities.yaml"},
					{source: "actions/dependabot.yaml", destination: ".github/dependabot.yaml"},
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
