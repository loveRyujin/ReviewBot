package cmd

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/form"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ReviewBot configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		model, err := form.Run()
		if err != nil {
			return err
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// set viper config file path
		viper.SetConfigFile(filepath.Join(home, defaultConfigDir, defaultConfigFile))

		// viper set ai configuration
		viper.Set("ai.provider", model.Provider())
		viper.Set("ai.api_key", model.ApiKey())
		viper.Set("ai.model", model.ModelName())
		viper.Set("ai.base_url", model.BaseURL())

		// viper set git configuration
		viper.Set("git.diff_file", model.DiffFile())
		viper.Set("git.max_input_size", model.MaxInputSize())
		viper.Set("git.diff_unified", model.DiffUnified())
		viper.Set("git.excluded_list", model.ExcludedList())
		viper.Set("git.amend", model.Amend())
		viper.Set("git.lang", model.Lang())

		// viper set proxy configuration
		viper.Set("proxy.proxy_url", model.ProxyURL())
		viper.Set("proxy.socks_url", model.SocksURL())
		viper.Set("proxy.timeout", model.Timeout())
		viper.Set("proxy.headers", model.Headers())
		viper.Set("proxy.skip_verify", model.SkipVerify())

		if err = viper.WriteConfig(); err != nil {
			return err
		}

		color.Green("Configuration initialized successfully! You can find the config file at: %s", viper.ConfigFileUsed())

		return nil
	},
}
