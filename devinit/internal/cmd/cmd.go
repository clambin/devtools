package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devinit",
	Short: "devinit sets up a Golang module",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("type", "t", "program", "module type (program, container, library)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug logging")
	rootCmd.PersistentFlags().Bool("dry-run", false, "do not create files")
	rootCmd.PersistentFlags().StringP("output", "o", ".", "output path")
	rootCmd.PersistentFlags().StringP("gomod", "g", "./go.mod", "path to module's go.mod file")

	rootCmd.AddCommand(workFlowsCmd, licensesCmd, allCmd)
}
