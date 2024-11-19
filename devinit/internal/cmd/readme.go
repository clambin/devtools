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
		Short: "Generate basic README.md",
		RunE: func(cmd *cobra.Command, args []string) error {
			moduleType, _ := cmd.Flags().GetString("type")
			output, _ := cmd.Flags().GetString("output")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			goModPath, _ := cmd.Flags().GetString("gomod")

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
					return fmt.Errorf("failed to create README.md: %w", err)
				}
				writeREADME(w, info, moduleType)
				_ = w.Close()
			}

			return err
		},
	}
)

func writeREADME(w io.Writer, info modInfo, moduleType string) {
	writeTitle(w, info)
	writeTag(w, info)
	writeCodeCov(w, info)
	if moduleType != "container" {
		writeTest(w, info)
	} else {
		writeBuild(w, info)
	}
	writeGoReport(w, info)
	if moduleType == "library" {
		writeGoDoc(w, info)
	}
	writeLicense(w, info)
}

func writeTitle(w io.Writer, info modInfo) {
	writeLine(w, "# "+filepath.Base(info.fullPath))
}

func writeTag(w io.Writer, info modInfo) {
	writeLink(w,
		"release",
		"https://img.shields.io/github/v/tag/"+info.strippedPath+"?color=green&label=release&style=plastic",
		"https://"+info.fullPath+"/releases",
	)
}

func writeCodeCov(w io.Writer, info modInfo) {
	writeLink(w,
		"codecov",
		"https://img.shields.io/codecov/c/gh/"+info.strippedPath+"?style=plastic",
		"https://app.codecov.io/gh/"+info.strippedPath,
	)
}

func writeTest(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "test")
}

func writeBuild(w io.Writer, info modInfo) {
	writeWorkFlowResult(w, info, "build")
}

func writeWorkFlowResult(w io.Writer, info modInfo, action string) {
	writeLink(w,
		action,
		"https://"+info.fullPath+"/workflows/"+action+"/badge.svg",
		"https://"+info.fullPath+"/actions",
	)
}

func writeGoReport(w io.Writer, info modInfo) {
	writeLink(w,
		"go report card",
		"https://goreportcard.com/badge/"+info.fullPath,
		"https://goreportcard.com/report/"+info.fullPath,
	)
}

func writeLicense(w io.Writer, info modInfo) {
	writeLink(w,
		"license",
		"https://img.shields.io/github/license/"+info.strippedPath+"?style=plastic",
		"LICENSE.md",
	)
}

func writeGoDoc(w io.Writer, info modInfo) {
	writeLink(w,
		"godoc",
		"https://pkg.go.dev/badge/"+info.fullPath+"?utm_source=godoc",
		"https://pkg.go.dev/"+info.fullPath,
	)
}

func writeLink(w io.Writer, label, image, link string) {
	output := "![" + label + "](" + image + ")"
	if link != "" {
		output = "[" + output + "](" + link + ")"
	}
	writeLine(w, output)
}

func writeLine(w io.Writer, line string) {
	_, _ = w.Write([]byte(line + "\n"))
}
