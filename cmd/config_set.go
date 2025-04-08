package cmd

import "github.com/spf13/cobra"

func init() {
	configCmd.AddCommand(configSetCmd)
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
