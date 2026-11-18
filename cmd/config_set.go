package cmd

import (
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/spf13/cobra"
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

		opts, err := configOptions()
		if err != nil {
			return err
		}
		var value any = args[1]
		switch args[0] {
		case "git.exclude_list":
			value = strings.Split(args[1], ",")
		case "proxy.headers":
			value = proxy.ParseHeaders(args[1])
		}
		path, err := config.SetFileValue(opts, args[0], value)
		if err != nil {
			return err
		}

		color.New(color.FgGreen).Fprintf(cmd.OutOrStdout(), "set config value successfully, you can see the config file: %s\n", path)

		return nil
	},
}
