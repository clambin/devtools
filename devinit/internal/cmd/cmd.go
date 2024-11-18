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
	rootCmd.PersistentFlags().StringP("type", "t", "program", "Module type (program, container, library)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().Bool("dry-run", false, "do not create files")
	rootCmd.PersistentFlags().StringP("output", "o", ".", "Path to the directory containing the module")
	rootCmd.PersistentFlags().StringP("gomod", "g", "./go.mod", "Path to module go.mod file")

	rootCmd.AddCommand(workFlowsCmd, licensesCmd, allCmd)
}
