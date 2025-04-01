package cmd

import "github.com/spf13/cobra"

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes in git stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
