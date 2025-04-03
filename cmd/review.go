package cmd

import (
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/llm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes in git stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := NewServerOptions()

		// applies user-provided configuration from the command line
		if err := opts.ApplyCfg(); err != nil {
			return err
		}

		// reads in config file and ENV variables
		if err := viper.Unmarshal(opts); err != nil {
			return err
		}

		// validates the options
		if err := opts.Validate(); err != nil {
			return err
		}

		g := opts.GitConfig().New()
		diff, err := g.DiffFiles()
		if err != nil {
			return err
		}

		provider := ai.Provider("viper.provider")
		client, err := llm.GetModelClient(provider)
		if err != nil {
			return err
		}
		resp, err := client.ChatCompletion(cmd.Context(), diff)
		if err != nil {
			return err
		}
		_ = resp

		return nil
	},
}
