package cmd

import (
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes in git stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ServerOption.Initialize(); err != nil {
			return err
		}

		g := ServerOption.GitConfig().New()
		diff, err := g.DiffFiles()
		if err != nil {
			return err
		}

		provider := ai.Provider(ServerOption.AiOptions.Provider)
		client, err := GetModelClient(provider)
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
