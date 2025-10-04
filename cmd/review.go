package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/pkg/progress"
	"github.com/loveRyujin/ReviewBot/prompt"
	"github.com/spf13/cobra"
)

const (
	ModeLocal    = "local"
	ModeExternal = "external"
)

var (
	mode         string
	diffFile     string
	maxInputSize int
	outputLang   string
	stream       bool
)

func init() {
	reviewCmd.PersistentFlags().IntVar(&diffUnifiedLines, "diff_unified", 3, "number of context lines to show in diff")
	reviewCmd.PersistentFlags().StringArrayVar(&excludedList, "exclude_list", []string{}, "list of files to exclude from review")
	reviewCmd.PersistentFlags().BoolVar(&amend, "amend", false, "amend the commit message")
	reviewCmd.PersistentFlags().StringVar(&mode, "mode", ModeLocal, "mode of fetch git diff information (local or external)")
	reviewCmd.PersistentFlags().StringVar(&diffFile, "diff_file", "", "path of the diff file to be reviewed")
	reviewCmd.PersistentFlags().IntVar(&maxInputSize, "max_input_size", 20*1024*1024, "maximum git diff input size(default: 20MB, units: bytes)")
	reviewCmd.PersistentFlags().StringVar(&outputLang, "output_lang", "en", "output language of the review summary(default: English)")
	reviewCmd.PersistentFlags().BoolVar(&stream, "stream", false, "enable streaming mode for AI provider")
}

// reviewCmd defines the "review" command for auto-reviewing staged git code changes using AI.
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Auto review code changes in git stage",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			cobra.CheckErr(err)
		}
		applyReviewOverrides()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// generate diff info
		diff, err := getDiffContent(args)
		if err != nil {
			return err
		}

		aiProvider := ai.Provider(globalConfig.AI.Provider)
		aiModelClient, err := GetModelClient(aiProvider)
		if err != nil {
			return err
		}

		// get file diff summary prompt for code review
		reviewPrompt, err := prompt.GetPromptTmpl(prompt.CodeReviewFileDiffTmpl, map[string]any{prompt.FileDiff: diff})
		if err != nil {
			return err
		}

		// get the language of output summary
		lang := prompt.GetLanguage(globalConfig.Git.Lang)

		color.Cyan("We are trying to review code changes")

		return executeReview(cmd.Context(), aiModelClient, reviewPrompt, lang)
	},
}

// getDiffContent get git diff content
func getDiffContent(args []string) (string, error) {
	var gitDiffContent string
	var err error

	switch mode {
	case ModeLocal:
		g := globalConfig.GitCommandConfig().New()
		gitDiffContent, err = g.DiffFiles()
	case ModeExternal:
		if len(args) != 0 { // get git diff content from arguments
			if len(args[0]) >= globalConfig.Git.MaxInputSize {
				return "", errors.New("git diff input size exceeds limit")
			}
			gitDiffContent = args[0]
		} else if diffFile != "" { // get git diff content from file
			gitDiffContent, err = processFileInput(globalConfig.Git.DiffFile)
		} else { // get git diff content from stdin
			gitDiffContent, err = processStdinInput()
		}

		if len(gitDiffContent) == 0 {
			return "", errors.New("please provide the diff content to review")
		}
	default:
		return "", errors.New("invalid input mode, please use 'local' or 'external'")
	}

	return gitDiffContent, err
}

// executeReview perform code review and process output
func executeReview(ctx context.Context, client ai.TextGenerator, reviewPrompt string, lang string) error {
	var summary string
	if stream { // streaming mode
		yellow := color.New(color.FgYellow).PrintfFunc()

		if lang != prompt.DefaultLanguage {
			resp, err := client.ChatCompletion(ctx, reviewPrompt)
			if err != nil {
				return err
			}
			summary = resp.Text
			color.Magenta(resp.TokenUsage.String())

			return streamTranslation(ctx, client, summary, lang, yellow)
		} else {
			return streamOutput(ctx, client, reviewPrompt, yellow)
		}
	} else { // non-streaming mode
		var resp *ai.Response
		var err error

		// Use spinner for AI call in non-streaming mode
		spinnerErr := progress.WithSpinnerAndCustomMessages(
			"🤖 Analyzing code changes...",
			"Code analysis completed",
			"Failed to analyze code changes",
			func() error {
				resp, err = client.ChatCompletion(ctx, reviewPrompt)
				return err
			},
		)
		if spinnerErr != nil {
			return spinnerErr
		}

		summary = resp.Text
		color.Magenta(resp.TokenUsage.String())

		if lang != prompt.DefaultLanguage {
			summary, err = translateContent(ctx, client, summary, lang)
			if err != nil {
				return err
			}
		}

		// output the review summary
		color.Yellow("================Review Summary====================")
		color.Yellow("\n" + strings.TrimSpace(summary) + "\n\n")
		color.Yellow("==================================================")
	}

	return nil
}

// streamOutput streams AI-generated review output in real-time with colored formatting and token usage stats.
func streamOutput(ctx context.Context, client ai.TextGenerator, reviewPrompt string, colorF func(format string, a ...interface{})) error {
	chunkHandler := func(chunk string) error {
		colorF(chunk)
		return nil
	}

	if err := client.StreamChatCompletion(ctx, reviewPrompt, chunkHandler); err != nil {
		return err
	}

	return nil
}

// streamTranslation streams the translation of a code review summary into the specified language using the AI client and colored output.
func streamTranslation(ctx context.Context, client ai.TextGenerator, content string, lang string, colorF func(format string, a ...interface{})) error {
	instruction, err := prompt.GetPromptTmpl(prompt.TranslationTmpl, map[string]any{
		prompt.OutputLang:    lang,
		prompt.OutputMessage: content,
	})
	if err != nil {
		return err
	}

	color.Cyan("We are trying to translate the code review summary to " + lang + " in streaming mode")
	return streamOutput(ctx, client, instruction, colorF)
}

// translateContent translates the given content into the specified language using the provided AI text generator.
func translateContent(ctx context.Context, client ai.TextGenerator, content string, lang string) (string, error) {
	instruction, err := prompt.GetPromptTmpl(prompt.TranslationTmpl, map[string]any{
		prompt.OutputLang:    lang,
		prompt.OutputMessage: content,
	})
	if err != nil {
		return "", err
	}

	color.Cyan("We are trying to translate the code review summary to " + lang)

	var resp *ai.Response
	spinnerErr := progress.WithSpinnerAndCustomMessages(
		"🌐 Translating review summary...",
		"Translation completed",
		"Failed to translate review summary",
		func() error {
			resp, err = client.ChatCompletion(ctx, instruction)
			return err
		},
	)
	if spinnerErr != nil {
		return "", spinnerErr
	}

	color.Magenta(resp.TokenUsage.String())
	return resp.Text, nil
}

// processFileInput reads the diff content from a file
func processFileInput(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	allowed, err := checkFileSize(file)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", errors.New("git diff input size exceeds limit")
	}

	diff, err := processInput(file)
	if err != nil {
		return "", err
	}

	return diff, nil
}

// processStdinInput reads the diff content from stdin
func processStdinInput() (string, error) {
	if !hasStdinInput() {
		return "", nil
	}

	allowed, err := checkFileSize(os.Stdin)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", errors.New("git diff input size exceeds limit")
	}

	diff, err := processInput(os.Stdin)
	if err != nil {
		return "", err
	}

	return diff, nil
}

// hasStdinInput checks if there is input from stdin.
func hasStdinInput() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		cobra.CheckErr(err)
	}

	return (fi.Mode() & os.ModeCharDevice) == 0
}

func processInput(r io.Reader) (string, error) {
	diffContent, err := io.ReadAll(r)
	return string(diffContent), err
}

// checkFileSize checks if the file size is valid(less than 20MB).
func checkFileSize(f *os.File) (bool, error) {
	stat, err := f.Stat()
	if err != nil {
		return false, err
	}

	if stat.Size() >= int64(globalConfig.Git.MaxInputSize) {
		return false, nil
	}

	return true, nil
}

// applyReviewOverrides applies command-line flags to the global configuration
func applyReviewOverrides() {
	if diffUnifiedLines != 3 {
		globalConfig.Git.DiffUnified = diffUnifiedLines
	}
	if len(excludedList) > 0 {
		globalConfig.Git.ExcludedList = append(globalConfig.Git.ExcludedList, excludedList...)
	}
	if amend {
		globalConfig.Git.Amend = true
	}
	if mode != "" {
		globalConfig.Runtime.Review.Mode = mode
	}
	if diffFile != "" {
		globalConfig.Git.DiffFile = diffFile
	}
	if maxInputSize != 20*1024*1024 {
		globalConfig.Git.MaxInputSize = maxInputSize
	}
	if outputLang != "en" {
		globalConfig.Git.Lang = outputLang
	}
	if stream {
		globalConfig.Runtime.Review.Stream = true
	}
}
