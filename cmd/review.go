package cmd

import (
	"strings"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/prompt"
	"github.com/spf13/cobra"
)

// reviewCmd represents the "review" command which automates the process of
// reviewing code changes in the git staging area. It initializes the server
// configuration, generates a diff of staged files, and uses an AI provider
// to analyze the changes and provide a review summary. The review summary
// is then displayed to the user in a formatted output.
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes in git stage",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ServerOption.Initialize(); err != nil {
			return err
		}

		// generate diff info
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
		// get file diff summary prompt for code review
		instruction, err := prompt.GetPromptTmpl(prompt.CodeReviewFileDiffTmpl, map[string]any{prompt.FileDiff: diff})
		if err != nil {
			return err
		}

		color.Cyan("We are trying to review code changes")
		resp, err := client.ChatCompletion(cmd.Context(), instruction)
		if err != nil {
			return err
		}
		summary := resp.Text
		color.Magenta(resp.TokenUsage.String())

		// Output core review summary
		color.Yellow("================Review Summary====================")
		color.Yellow("\n" + strings.TrimSpace(summary) + "\n\n")
		color.Yellow("==================================================")

		return nil
	},
}
