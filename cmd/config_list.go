package cmd

import "github.com/spf13/cobra"

func init() {
	configCmd.AddCommand(configListCmd)
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
