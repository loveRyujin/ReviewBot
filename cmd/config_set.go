package cmd

import (
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configSetCmd)
}

// configSetCmd represents the "set" command which allows users to set a
// configuration value. It requires at least two arguments: a key and a value.
// The command validates the key against a predefined list of available keys
// and updates the configuration using Viper. If the key is "git.exclude_list",
// the value is split into a list using commas. The updated configuration is
// then written to the configuration file. On success, a confirmation message
// is displayed with the path to the configuration file.
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration value",
	Args:  cobra.MinimumNArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			cobra.CheckErr(err)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if the key is valid
		if _, exist := availableKeys[args[0]]; !exist {
			return errors.New("invalid config key, please use 'reviewbot config list' to see all available keys")
		}

		// set the config value in viper
		if args[0] == "git.exclude_list" {
			viper.Set(args[0], strings.Split(args[1], ","))
		} else {
			viper.Set(args[0], args[1])
		}

		// write config to file
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		color.Green("set config value successfully, you can see the config file: %s", viper.ConfigFileUsed())

		return nil
	},
}
