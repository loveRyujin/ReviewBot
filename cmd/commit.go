package cmd

import "github.com/spf13/cobra"

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

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Automically generate commit message",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
