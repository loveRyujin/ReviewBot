package cmd

import (
	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/git"
	"github.com/loveRyujin/ReviewBot/prompt"
	"github.com/spf13/cobra"
)

var (
	diffUnifiedLines int
	excludedList     []string
	amend            bool
)

func init() {
	commitCmd.PersistentFlags().IntVar(&diffUnifiedLines, "diff_unified", 3, "number of context lines to show in diff")
	commitCmd.PersistentFlags().StringArrayVar(&excludedList, "exclude_list", []string{}, "list of files to exclude from review")
	commitCmd.PersistentFlags().BoolVar(&amend, "amend", false, "amend the commit message")
}

// commitCmd is a Cobra command that automates the generation of commit messages
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Automically generate commit message",
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

		currentModel := ServerOption.AiOptions.Model
		provider := ai.Provider(ServerOption.AiOptions.Provider)
		client, err := GetModelClient(provider)
		if err != nil {
			return err
		}
		// get file diff summary prompt for commit message
		instruction, err := prompt.GetFileDiffSummaryTmplForCommit(prompt.FileDiff, diff)
		if err != nil {
			return err
		}

		color.Green("Using %s model for commit message generation\n", currentModel)
		color.Green("We are trying to generate commit message\n")

		// generate file diff summary
		color.Cyan("Generating file diff summary...\n")
		resp, err := client.ChatCompletion(cmd.Context(), instruction)
		if err != nil {
			return err
		}
		summary := resp.Text
		color.Magenta(resp.TokenUsage.String())

		// generate commit message prefix
		color.Cyan("Generating commit message prefix...\n")
		instruction, err = prompt.GetCommitMessagePrefixTmpl(prompt.SummaryPoint, summary)
		if err != nil {
			return err
		}
		resp, err = client.ChatCompletion(cmd.Context(), instruction)
		if err != nil {
			return err
		}
		prefix := resp.Text
		color.Magenta(resp.TokenUsage.String())

		// generate commit message title
		color.Cyan("Generating commit message title...\n")
		instruction, err = prompt.GetCommitMessageTitleTmpl(prompt.SummaryPoint, summary)
		if err != nil {
			return err
		}
		resp, err = client.ChatCompletion(cmd.Context(), instruction)
		if err != nil {
			return err
		}
		title := resp.Text
		color.Magenta(resp.TokenUsage.String())

		// generate commit message
		commitMsg, err := git.GetCommitMessageTmpl(map[string]any{
			git.CommitMessagePrefix:  prefix,
			git.CommitMessageTitle:   title,
			git.CommitMessageSummary: summary,
		})
		if err != nil {
			return err
		}

		// Output commit message from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + commitMsg + "\n\n")
		color.Yellow("==================================================")

		output, err := g.Commit(commitMsg)
		if err != nil {
			return err
		}
		color.Yellow(output)

		return nil
	},
}
