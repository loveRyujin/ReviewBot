package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/loveRyujin/ReviewBot/pkg/version"
	"github.com/loveRyujin/ReviewBot/prompt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configPath       string
	defaultEnvPrefix = "REVIEWBOT"
	replacer         = strings.NewReplacer(".", "_", "-", "_")
	aiProviderFlag   string
	aiModelFlag      string

	defaultConfigDir  = ".config/reviewbot"
	defaultConfigFile = "reviewbot.yaml"

	// globalConfig holds the loaded application configuration
	globalConfig *config.Config
)

var rootCmd = &cobra.Command{
	Use:          "reviewbot",
	Short:        "A command-line tool that helps generate git commit messages, code reviews, etc.",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		welcome()
		version.PrintAndExitIfRequested()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(commitCmd)

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&aiProviderFlag, "ai-provider", "", "AI provider to use for requests")
	rootCmd.PersistentFlags().StringVar(&aiModelFlag, "ai-model", "", "AI model identifier to use")

	version.AddFlags(rootCmd.Flags())

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// save the original help function and set a custom one
	originalHelpFunc := rootCmd.HelpFunc()
	helpF := func(cmd *cobra.Command, args []string) {
		printProjectName()
		originalHelpFunc(cmd, args)
	}
	rootCmd.SetHelpFunc(helpF)
}

func configOptions() (config.LoadOptions, error) {
	opts := config.LoadOptions{
		ExplicitPath: configPath,
		ConfigName:   strings.TrimSuffix(defaultConfigFile, ".yaml"),
		ConfigType:   "yaml",
		EnvPrefix:    defaultEnvPrefix,
		Replacer:     replacer,
	}
	if configPath != "" {
		return opts, nil
	}
	configDir, err := resolveDefaultConfigDir()
	if err != nil {
		return opts, err
	}
	opts.SearchDirs = []string{"./config/", ".", configDir}
	return opts, ensureDefaultConfigFile()
}

func loadCommandConfig(cmd *cobra.Command, _ []string) error {
	opts, err := configOptions()
	if err != nil {
		return err
	}
	flags := cmd.Flags()
	ov := &opts.Overrides
	if flags.Changed("diff_unified") {
		ov.Git.DiffUnified = &diffUnifiedLines
	}
	if flags.Changed("exclude_list") {
		ov.Git.ExcludedList = excludedList
	}
	if flags.Changed("amend") {
		ov.Git.Amend = &amend
	}
	if flags.Changed("mode") {
		ov.Review.Mode = mode
	}
	if flags.Changed("stream") {
		ov.Review.Stream = &stream
	}
	if flags.Changed("preview") {
		ov.Commit.Preview = &preview
	}
	if flags.Changed("ai-provider") {
		ov.Provider = aiProviderFlag
	}
	if flags.Changed("ai-model") {
		ov.Model = aiModelFlag
	}

	globalConfig, err = config.Load(opts)
	if err != nil {
		return err
	}

	git := &globalConfig.Git
	if cmd.Name() == "review" {
		review := globalConfig.Runtime.Review
		if review.DiffFile != "" {
			git.DiffFile = review.DiffFile
		}
		if review.MaxInput > 0 {
			git.MaxInputSize = review.MaxInput
		}
		if review.OutputLang != "" {
			git.Lang = review.OutputLang
		}
	} else if lang := globalConfig.Runtime.Commit.OutputLang; lang != "" {
		git.Lang = lang
	}
	if flags.Changed("diff_file") {
		git.DiffFile = diffFile
	}
	if flags.Changed("max_input_size") {
		git.MaxInputSize = maxInputSize
	}
	if flags.Changed("output_lang") {
		git.Lang = outputLang
	}
	if err := git.Validate(); err != nil {
		return err
	}
	prompt.SetTemplateDir(globalConfig.Prompt.Folder)
	return nil
}

func welcome() {
	printProjectName()
	color.Blue("Welcome to ReviewBot! 🎉")
	color.Blue("Run `reviewbot --help` for quick start.")
	color.Blue("You can also run `reviewbot init` to generate a configuration file.")
	color.Blue("For usage instructions and examples, please check the README documentation: https://github.com/loveRyujin/ReviewBot/blob/master/README.md .")
}

func printProjectName() {
	reviewbotFigure := figure.NewFigure("ReviewBot", "", true)
	color.Cyan("%s\n", reviewbotFigure.String())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ensureDefaultConfigFile writes the default config to the user config dir if none exists.
func ensureDefaultConfigFile() error {
	configDir, err := resolveDefaultConfigDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, defaultConfigFile)

	if _, err := os.Stat(configFile); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	defaultCfg := config.NewDefault()
	content, err := yaml.Marshal(defaultCfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, content, 0o600)
}

func resolveDefaultConfigDir() (string, error) {
	if configDir, err := os.UserConfigDir(); err == nil && configDir != "" {
		return filepath.Join(configDir, "reviewbot"), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, defaultConfigDir), nil
}
