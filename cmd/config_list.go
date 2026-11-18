package cmd

import (
	"maps"
	"slices"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configListCmd)
}

// availableKeys is a map of configuration keys and their descriptions
var availableKeys = map[string]string{
	"git.diff_file":                 "Path to the diff file to be reviewed",
	"git.max_input_size":            "Maximum git diff input size (default: 20MB, units: bytes)",
	"git.diff_unified":              "Number of context lines in git diff output (default: 3)",
	"git.exclude_list":              "Files to exclude from git diff command",
	"git.lang":                      "Language for summarization output (default: English)",
	"git.amend":                     "Amend the previous commit",
	"proxy.socks_url":               "SOCKS proxy URL for API connections",
	"ai.api_key":                    "Authentication key for the selected provider",
	"ai.model":                      "AI model identifier to use for requests",
	"proxy.proxy_url":               "HTTP proxy URL for API connections",
	"ai.base_url":                   "Custom base URL for API requests",
	"proxy.timeout":                 "Maximum duration to wait for API response",
	"ai.max_tokens":                 "Maximum token limit for generated completions",
	"ai.temperature":                "Randomness control parameter (0-1): lower values for focused results, higher for creative variety",
	"ai.provider":                   "Service provider selection (openai, deepseek, gemini)",
	"proxy.skip_verify":             "Option to bypass TLS certificate verification",
	"proxy.headers":                 "Additional custom HTTP headers for API requests",
	"ai.top_p":                      "Nucleus sampling parameter: controls diversity by limiting to top percentage of probability mass",
	"ai.frequency_penalty":          "Parameter to reduce repetition by penalizing tokens based on their frequency",
	"ai.presence_penalty":           "Parameter to encourage topic diversity by penalizing previously used tokens",
	"runtime.commit.preview":        "Preview and confirm before committing",
	"runtime.commit.output_lang":    "Commit message language",
	"runtime.review.mode":           "Review input mode (local or external)",
	"runtime.review.stream":         "Stream review output",
	"runtime.review.diff_file":      "External diff file for review",
	"runtime.review.output_lang":    "Review output language",
	"runtime.review.max_input_size": "Maximum review input size in bytes",
	"prompt.folder":                 "Directory path for custom prompt templates",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts, err := configOptions()
		if err != nil {
			return err
		}
		settings, err := config.ReadSettings(opts)
		if err != nil {
			return err
		}
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create a new table with the header "Config Name" and "Value"
		tbl := table.New("Config Name", "Value", "Description")
		tbl.WithWriter(cmd.OutOrStdout()).WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		keys := slices.Sorted(maps.Keys(availableKeys))

		// Add the key and value to the table
		for _, v := range keys {
			// Hide the api key
			if v == "ai.api_key" || v == "proxy.headers" {
				tbl.AddRow(v, "****************", availableKeys[v])
				continue
			}
			tbl.AddRow(v, settings.Get(v), availableKeys[v])
		}

		// Print the table
		tbl.Print()
		return nil
	},
}
