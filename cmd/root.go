package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/loveRyujin/ReviewBot/pkg/version"
	"github.com/spf13/cobra"
)

var (
	configPath       string
	defaultEnvPrefix = "REVIEWBOT"
	replacer         = strings.NewReplacer(".", "_", "-", "_")

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

// initConfig loads configuration using the new config package.
func initConfig() error {
	var err error
	globalConfig, err = config.Load(config.LoadOptions{
		ExplicitPath: configPath,
		SearchDirs:   searchDirs(),
		ConfigName:   strings.TrimSuffix(defaultConfigFile, ".yaml"),
		ConfigType:   "yaml",
		EnvPrefix:    defaultEnvPrefix,
		Replacer:     replacer,
	})
	return err
}

// searchDirs returns the directories to search for the config file.
func searchDirs() []string {
	// get user home dir
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return []string{"./config/", ".", filepath.Join(homeDir, defaultConfigDir)}
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
