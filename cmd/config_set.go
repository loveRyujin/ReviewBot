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

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration value",
	Args:  cobra.MinimumNArgs(2),
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
