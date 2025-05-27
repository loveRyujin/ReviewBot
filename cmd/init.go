package cmd

import (
	"os"
	"path/filepath"

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

		viper.SetConfigFile(filepath.Join(home, defaultConfigDir, defaultConfigFile))
		viper.Set("ai.provider", model.Provider())
		viper.Set("ai.api_key", model.ApiKey())
		viper.Set("ai.model", model.ModelName())
		viper.Set("ai.base_url", model.BaseURL())

		err = viper.WriteConfig()

		return err
	},
}
