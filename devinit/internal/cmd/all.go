package cmd

import "github.com/spf13/cobra"

var (
	allCmd = &cobra.Command{
		Use:   "all",
		Short: "create all supported repository files",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, f := range []func(*cobra.Command, []string) error{
				workFlowsCmd.RunE,
				readmeCmd.RunE,
				licensesCmd.RunE,
			} {
				if err := f(cmd, args); err != nil {
					return err
				}
			}
			return nil
		},
	}
)
