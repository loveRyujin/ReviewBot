package cmd

import (
	"html"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/git"
	"github.com/loveRyujin/ReviewBot/prompt"
	"github.com/spf13/cobra"
)

var (
	preview          bool
	diffUnifiedLines int
	excludedList     []string
	amend            bool
)

func init() {
	commitCmd.PersistentFlags().BoolVar(&preview, "preview", false, "preview the commit message before committing")
	commitCmd.PersistentFlags().IntVar(&diffUnifiedLines, "diff_unified", 3, "number of context lines to show in diff")
	commitCmd.PersistentFlags().StringArrayVar(&excludedList, "exclude_list", []string{}, "list of files to exclude from review")
	commitCmd.PersistentFlags().BoolVar(&amend, "amend", false, "amend the commit message")
	commitCmd.PersistentFlags().StringVar(&outputLang, "output_lang", "en", "output language of the commit message(default: English)")
}

// commitCmd is a Cobra command that automates the generation of commit messages
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Automically generate commit message",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			cobra.CheckErr(err)
		}
		applyCommitOverrides()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// generate diff info
		g := globalConfig.GitCommandConfig().New()
		diff, err := g.DiffFiles()
		if err != nil {
			return err
		}

		currentModel := globalConfig.AI.Model
		provider := ai.Provider(globalConfig.AI.Provider)
		client, err := GetModelClient(provider)
		if err != nil {
			return err
		}
		// get file diff summary prompt for commit message
		instruction, err := prompt.GetPromptTmpl(prompt.CommitFileDiffTmpl, map[string]any{prompt.FileDiff: diff})
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
		instruction, err = prompt.GetPromptTmpl(prompt.CommitMessagePrefixTmpl, map[string]any{prompt.SummaryPoint: summary})
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
		instruction, err = prompt.GetPromptTmpl(prompt.CommitMessageTitleTmpl, map[string]any{prompt.SummaryPoint: summary})
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

		escapeCommitMsg := html.UnescapeString(commitMsg)
		lang := prompt.GetLanguage(globalConfig.Git.Lang)
		commitOutput := escapeCommitMsg

		if lang != prompt.DefaultLanguage {
			translated, err := translateContent(cmd.Context(), client, escapeCommitMsg, lang)
			if err != nil {
				return err
			}
			commitOutput = translated
		}

		// Output commit message from AI
		color.Yellow("================Commit Summary====================")
		color.Yellow("\n" + commitOutput + "\n")
		color.Yellow("==================================================")

		if preview {
			ready, err := confirmation.New("\nWhether to commit this preview message?", confirmation.Yes).RunPrompt()
			if err != nil {
				return err
			}
			if !ready {
				return nil
			}
		}

		output, err := g.Commit(commitOutput)
		if err != nil {
			return err
		}
		color.Yellow(output)

		return nil
	},
}

// applyCommitOverrides applies command-line flags to the global configuration
func applyCommitOverrides() {
	if diffUnifiedLines != 3 {
		globalConfig.Git.DiffUnified = diffUnifiedLines
	}
	if len(excludedList) > 0 {
		globalConfig.Git.ExcludedList = append(globalConfig.Git.ExcludedList, excludedList...)
	}
	if amend {
		globalConfig.Git.Amend = true
	}
	if outputLang != "en" {
		globalConfig.Git.Lang = outputLang
	}
	if preview {
		globalConfig.Runtime.Commit.Preview = true
	}
}
