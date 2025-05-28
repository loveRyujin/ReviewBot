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

		settingsToSet := make(map[string]any)

		// viper set ai configuration
		if val := model.Provider(); val != "" {
			settingsToSet["ai.provider"] = val
		}
		if val := model.ApiKey(); val != "" {
			settingsToSet["ai.api_key"] = val
		}
		if val := model.BaseURL(); val != "" {
			settingsToSet["ai.base_url"] = val
		}
		if val := model.Model(); val != "" {
			settingsToSet["ai.model"] = val
		}

		// viper set git configuration
		if val := model.DiffFile(); val != "" {
			settingsToSet["git.diff_file"] = val
		}
		if val := model.MaxInputSize(); val != "" {
			settingsToSet["git.max_input_size"] = val
		}
		if val := model.DiffUnified(); val != "" {
			settingsToSet["git.diff_unified"] = val
		}
		if val := model.ExcludedList(); len(val) > 0 {
			settingsToSet["git.excluded_list"] = val
		}
		if val := model.Amend(); val != "" {
			settingsToSet["git.amend"] = val
		}
		if val := model.Lang(); val != "" {
			settingsToSet["git.lang"] = val
		}

		// viper set proxy configuration
		if val := model.ProxyURL(); val != "" {
			settingsToSet["proxy.proxy_url"] = val
		}
		if val := model.SocksURL(); val != "" {
			settingsToSet["proxy.socks_url"] = val
		}
		if val := model.Timeout(); val != "" {
			settingsToSet["proxy.timeout"] = val
		}
		if val := model.Headers(); len(val) > 0 {
			settingsToSet["proxy.headers"] = val
		}
		if val := model.SkipVerify(); val != "" {
			settingsToSet["proxy.skip_verify"] = val
		}

		// create config file if at least one setting was provided
		if len(settingsToSet) > 0 {
			for key, value := range settingsToSet {
				viper.Set(key, value)
			}
			if err = viper.WriteConfig(); err != nil {
				return err
			}

			color.Green("Configuration initialized successfully! You can find the config file at: %s", viper.ConfigFileUsed())
		}

		return nil
	},
}
