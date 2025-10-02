package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Load builds a Config from Viper sources and applies overrides.
type LoadOptions struct {
	ExplicitPath string
	SearchDirs   []string
	ConfigName   string
	ConfigType   string
	EnvPrefix    string
	Replacer     *strings.Replacer
	Overrides    Overrides
}

// Overrides captures per-domain flag overrides merged on top of config.
type Overrides struct {
	Git    GitOverrides
	Review ReviewOverrides
	Commit CommitOverrides
}

// GitOverrides holds CLI overrides for git settings.
type GitOverrides struct {
	DiffUnified  *int
	ExcludedList []string
	DiffFile     string
	Amend        *bool
	MaxInputSize *int
	Lang         string
}

// ReviewOverrides holds CLI overrides for review runtime options.
type ReviewOverrides struct {
	Mode       string
	Stream     *bool
	DiffFile   string
	MaxInput   *int
	OutputLang string
}

// CommitOverrides holds CLI overrides for commit runtime options.
type CommitOverrides struct {
	Preview    *bool
	OutputLang string
}

// Load constructs a Config using defaults, file/env values, and overrides.
func Load(opts LoadOptions) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	if opts.Replacer == nil {
		opts.Replacer = strings.NewReplacer(".", "_", "-", "_")
	}

	if opts.ConfigName == "" {
		opts.ConfigName = "reviewbot"
	}
	if opts.ConfigType == "" {
		opts.ConfigType = "yaml"
	}

	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	}
	v.SetEnvKeyReplacer(opts.Replacer)
	v.AutomaticEnv()

	if opts.ExplicitPath != "" {
		v.SetConfigFile(opts.ExplicitPath)
	} else {
		for _, dir := range opts.SearchDirs {
			v.AddConfigPath(dir)
		}
		v.SetConfigName(opts.ConfigName)
		v.SetConfigType(opts.ConfigType)
	}

	if err := readConfig(v); err != nil {
		return nil, err
	}

	cfg := NewDefault()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	applyOverrides(cfg, opts.Overrides)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// readConfig loads the config file when available, ignoring missing files.
func readConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil
		}

		if strings.Contains(err.Error(), "Not Found") {
			return nil
		}

		return err
	}

	return nil
}

// applyOverrides merges CLI overrides into the config snapshot.
func applyOverrides(cfg *Config, ov Overrides) {
	if ov.Git.DiffUnified != nil {
		cfg.Git.DiffUnified = *ov.Git.DiffUnified
	}
	if len(ov.Git.ExcludedList) > 0 {
		cfg.Git.ExcludedList = append(cfg.Git.ExcludedList, ov.Git.ExcludedList...)
	}
	if ov.Git.DiffFile != "" {
		cfg.Git.DiffFile = ov.Git.DiffFile
	}
	if ov.Git.Amend != nil {
		cfg.Git.Amend = *ov.Git.Amend
	}
	if ov.Git.MaxInputSize != nil {
		cfg.Git.MaxInputSize = *ov.Git.MaxInputSize
	}
	if ov.Git.Lang != "" {
		cfg.Git.Lang = ov.Git.Lang
	}

	if ov.Review.Mode != "" {
		cfg.Runtime.Review.Mode = ov.Review.Mode
	}
	if ov.Review.Stream != nil {
		cfg.Runtime.Review.Stream = *ov.Review.Stream
	}
	if ov.Review.DiffFile != "" {
		cfg.Runtime.Review.DiffFile = ov.Review.DiffFile
	}
	if ov.Review.MaxInput != nil {
		cfg.Runtime.Review.MaxInput = *ov.Review.MaxInput
	}
	if ov.Review.OutputLang != "" {
		cfg.Runtime.Review.OutputLang = ov.Review.OutputLang
	}

	if ov.Commit.Preview != nil {
		cfg.Runtime.Commit.Preview = *ov.Commit.Preview
	}
	if ov.Commit.OutputLang != "" {
		cfg.Runtime.Commit.OutputLang = ov.Commit.OutputLang
	}
}

// ResolveConfigPath returns an absolute path for the provided config file.
func ResolveConfigPath(file string) (string, error) {
	if file == "" {
		return "", nil
	}

	if filepath.IsAbs(file) {
		return file, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve config path: %w", err)
	}

	return filepath.Join(cwd, file), nil
}
