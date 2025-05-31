package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configListCmd)
}

// availableKeys is a map of configuration keys and their descriptions
var availableKeys = map[string]string{
	"git.diff_file":        "Path to the diff file to be reviewed",
	"git.max_input_size":   "Maximum git diff input size (default: 20MB, units: bytes)",
	"git.diff_unified":     "Number of context lines in git diff output (default: 3)",
	"git.exclude_list":     "Files to exclude from git diff command",
	"git.lang":             "Language for summarization output (default: English)",
	"git.template_file":    "Path to template file for commit messages",
	"git.template_string":  "Template string for formatting commit messages",
	"ai.socks":             "SOCKS proxy URL for API connections",
	"ai.api_key":           "Authentication key for OpenAI API access",
	"ai.model":             "AI model identifier to use for requests",
	"ai.proxy":             "HTTP proxy URL for API connections",
	"ai.base_url":          "Custom base URL for API requests",
	"ai.timeout":           "Maximum duration to wait for API response",
	"ai.max_tokens":        "Maximum token limit for generated completions",
	"ai.temperature":       "Randomness control parameter (0-1): lower values for focused results, higher for creative variety",
	"ai.provider":          "Service provider selection ('openai' or 'azure')",
	"ai.skip_verify":       "Option to bypass TLS certificate verification",
	"ai.headers":           "Additional custom HTTP headers for API requests",
	"ai.top_p":             "Nucleus sampling parameter: controls diversity by limiting to top percentage of probability mass",
	"ai.frequency_penalty": "Parameter to reduce repetition by penalizing tokens based on their frequency",
	"ai.presence_penalty":  "Parameter to encourage topic diversity by penalizing previously used tokens",
	"prompt.folder":        "Directory path for custom prompt templates",
}

// configListCmd represents the "list" command which lists all configuration settings.
// It displays the settings in a formatted table with the configuration name and value.
// Sensitive information like "openai.api_key" is masked for security purposes.
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings",
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		// Create a new table with the header "Config Name" and "Value"
		tbl := table.New("Config Name", "Value", "Description")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		// Sort the keys
		keys := make([]string, 0, len(availableKeys))
		for key := range availableKeys {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		// Add the key and value to the table
		for _, v := range keys {
			// Hide the api key
			if v == "openai.api_key" {
				tbl.AddRow(v, "****************", availableKeys[v])
				continue
			}
			tbl.AddRow(v, viper.Get(v), availableKeys[v])
		}

		// Print the table
		tbl.Print()
	},
}
