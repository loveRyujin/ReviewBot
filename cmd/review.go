package cmd

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
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
		var diff string
		var err error

		switch mode {
		case ModeLocal:
			g := ServerOption.GitConfig().New()
			diff, err = g.DiffFiles()
			if err != nil {
				return err
			}
		case ModeExternal:
			if len(args) != 0 { // if args is not empty, use the first argument as the diff content
				if len(args[0]) >= ServerOption.GitOptions.MaxInputSize {
					return errors.New("git diff input size exceeds limit")
				}

				diff = args[0]
			} else if diffFile != "" { // if diffFile is provided, read the content from the file
				diff, err = processFileInput(ServerOption.GitOptions.DiffFile)
				if err != nil {
					return err
				}
			} else { // if no args or diffFile is provided, read from stdin
				diff, err = processStdinInput()
				if err != nil {
					return err
				}
			}

			if len(diff) == 0 {
				return errors.New("please provide the diff content to review")
			}
		default:
			return errors.New("invalid input mode, please use 'local' or 'external'")
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

		lang := prompt.GetLanguage(ServerOption.GitOptions.Lang)
		yellow := color.New(color.FgYellow).PrintfFunc()
		if stream { // streaming mode
			if lang != prompt.DefaultLanguage {
				resp, err := client.ChatCompletion(cmd.Context(), instruction)
				if err != nil {
					return err
				}
				summary := resp.Text
				color.Magenta(resp.TokenUsage.String())

				instruction, err = prompt.GetPromptTmpl(prompt.TranslationTmpl, map[string]any{
					prompt.OutputLang:    lang,
					prompt.OutputMessage: summary,
				})
				if err != nil {
					return err
				}

				color.Cyan("We are trying to translate the code review summary to " + lang + " in streaming mode")
				stream, err := client.StreamChatCompletion(cmd.Context(), instruction)
				if err != nil {
					return err
				}
				defer stream.Close()

				color.Yellow("================Review Summary====================" + "\n\n")

				tokenUsage := ai.TokenUsage{}
				for {
					resp, err := stream.Recv()
					if err != nil {
						if errors.Is(err, io.EOF) {
							color.Yellow("\n" + "==================================================")
							break
						}
						return err
					}

					yellow(resp.Choices[0].Delta.Content)

					if resp.Usage != nil {
						tokenUsage.PromptTokens = resp.Usage.PromptTokens
						tokenUsage.CompletionTokens = resp.Usage.CompletionTokens
						tokenUsage.TotalTokens = resp.Usage.TotalTokens
						tokenUsage.PromptTokensDetails = resp.Usage.PromptTokensDetails
						tokenUsage.CompletionTokensDetails = resp.Usage.CompletionTokensDetails
					}
				}
				color.Magenta(tokenUsage.String())
			} else {
				stream, err := client.StreamChatCompletion(cmd.Context(), instruction)
				if err != nil {
					return err
				}
				defer stream.Close()

				color.Yellow("================Review Summary====================" + "\n\n")

				tokenUsage := ai.TokenUsage{}
				for {
					resp, err := stream.Recv()
					if err != nil {
						if errors.Is(err, io.EOF) {
							color.Yellow("\n" + "==================================================")
							break
						}
						return err
					}

					yellow("%s", resp.Choices[0].Delta.Content)

					if resp.Usage != nil {
						tokenUsage.PromptTokens = resp.Usage.PromptTokens
						tokenUsage.CompletionTokens = resp.Usage.CompletionTokens
						tokenUsage.TotalTokens = resp.Usage.TotalTokens
						tokenUsage.PromptTokensDetails = resp.Usage.PromptTokensDetails
						tokenUsage.CompletionTokensDetails = resp.Usage.CompletionTokensDetails
					}
				}
				color.Magenta(tokenUsage.String())
			}
		} else { // non-streaming mode
			resp, err := client.ChatCompletion(cmd.Context(), instruction)
			if err != nil {
				return err
			}
			summary := resp.Text
			color.Magenta(resp.TokenUsage.String())

			if lang != prompt.DefaultLanguage {
				// get translation prompt
				instruction, err := prompt.GetPromptTmpl(prompt.TranslationTmpl, map[string]any{
					prompt.OutputLang:    ServerOption.GitOptions.Lang,
					prompt.OutputMessage: summary,
				})
				if err != nil {
					return err
				}

				color.Cyan("We are trying to translate the code review summary to " + lang)
				resp, err = client.ChatCompletion(cmd.Context(), instruction)
				if err != nil {
					return err
				}
				summary = resp.Text
				color.Magenta(resp.TokenUsage.String())
			}

			// Output core review summary
			color.Yellow("================Review Summary====================")
			color.Yellow("\n" + strings.TrimSpace(summary) + "\n\n")
			color.Yellow("==================================================")
		}

		return nil
	},
}

// processFileInput reads the diff content from a file
func processFileInput(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

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

	if stat.Size() >= int64(ServerOption.GitOptions.MaxInputSize) {
		return false, nil
	}

	return true, nil
}
