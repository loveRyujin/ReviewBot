package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configPath       string
	defaultEnvPrefix = "REVIEWBOT"
	replacer         = strings.NewReplacer(".", "_", "-", "_")

	defaultConfigDir  = ".config/reviewbot"
	defaultConfigFile = "reviewbot.yaml"

	once sync.Once
)

var rootCmd = &cobra.Command{
	Use:          "reviewbot",
	Short:        "A command-line tool that helps generate git commit messages, code reviews, etc.",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		printProjectName()
		color.Blue("Welcome to ReviewBot! 🎉")
		color.Blue("Run `reviewbot --help` for quick start.")
		color.Blue("You can also run `reviewbot init` to generate a configuration file.")
		color.Blue("For usage instructions and examples, please check the README documentation: https://github.com/loveRyujin/ReviewBot/blob/master/README.md .")
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

	once.Do(func() {
		ServerOption = NewServerOptions()
	})
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		for _, dir := range searchDirs() {
			viper.AddConfigPath(dir)
		}

		viper.SetConfigName(defaultConfigFile)
		viper.SetConfigType("yaml")
	}

	setupEnvironmentVariables()

	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
}

// setupEnvironmentVariables sets up the environment variables for viper.
func setupEnvironmentVariables() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.SetEnvKeyReplacer(replacer)
}

// searchDirs returns the directories to search for the config file.
func searchDirs() []string {
	// get user home dir
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return []string{"./config/", ".", filepath.Join(homeDir, defaultConfigDir)}
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
