package cmd

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ReviewBot configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
