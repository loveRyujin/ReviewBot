package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	Short:        "help code review when merging code",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
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
	return []string{filepath.Join(homeDir, defaultConfigDir), ".", "./config/"}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
