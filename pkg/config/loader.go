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
	Provider string
	Model    string
	Git      GitOverrides
	Review   ReviewOverrides
	Commit   CommitOverrides
	Prompt   PromptOverrides
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

// PromptOverrides allows overriding prompt-related settings via CLI.
type PromptOverrides struct {
	Folder string
}

// Load constructs a Config using defaults, file/env values, and overrides.
func Load(opts LoadOptions) (*Config, error) {
	v, err := ReadSettings(opts)
	if err != nil {
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

// ReadSettings reads file values, environment overrides and defaults for inspection or decoding.
func ReadSettings(opts LoadOptions) (*viper.Viper, error) {
	v, err := ReadFile(opts)
	if err != nil {
		return nil, err
	}
	setDefaults(v)
	if opts.Replacer == nil {
		opts.Replacer = strings.NewReplacer(".", "_", "-", "_")
	}
	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	}
	v.SetEnvKeyReplacer(opts.Replacer)
	v.AutomaticEnv()

	return v, nil
}

// ReadFile opens only the persisted settings so edits never save environment overrides or defaults.
func ReadFile(opts LoadOptions) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	if opts.ConfigName == "" {
		opts.ConfigName = "reviewbot"
	}
	if opts.ConfigType == "" {
		opts.ConfigType = "yaml"
	}

	if opts.ExplicitPath != "" {
		v.SetConfigFile(opts.ExplicitPath)
		if filepath.Ext(opts.ExplicitPath) == "" {
			v.SetConfigType(opts.ConfigType)
		}
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

	return v, nil
}

// SetFileValue updates the selected file without loading environment values or defaults.
func SetFileValue(opts LoadOptions, key string, value any) (string, error) {
	return SetFileValues(opts, map[string]any{key: value})
}

// SetFileValues merges the supplied fields into the document in a single write.
func SetFileValues(opts LoadOptions, values map[string]any) (string, error) {
	path := opts.ExplicitPath
	if path == "" {
		v, err := ReadFile(opts)
		if err != nil {
			return "", err
		}
		path = v.ConfigFileUsed()
		if path == "" {
			return "", errors.New("no config file selected")
		}
	}

	format := strings.TrimPrefix(filepath.Ext(path), ".")
	if format == "" {
		format = opts.ConfigType
	}
	codecs := viper.NewCodecRegistry()
	decoder, err := codecs.Decoder(format)
	if err != nil {
		return "", fmt.Errorf("config file %s: %w", path, err)
	}
	// Viper's AllSettings flattens keys and drops empty maps. Edit the
	// decoded document so unrelated settings survive a read/write cycle.
	document := make(map[string]any)
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := decoder.Decode(data, document); err != nil {
			return "", err
		}
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return "", err
		}
	default:
		return "", err
	}
	for key, value := range values {
		if err := setDocumentValue(document, key, value); err != nil {
			return "", err
		}
	}
	encoder, err := codecs.Encoder(format)
	if err != nil {
		return "", err
	}
	data, err = encoder.Encode(document)
	if err != nil {
		return "", err
	}
	return path, os.WriteFile(path, data, 0o600)
}

// Match the longest existing key first, as Viper also accepts dotted keys.
func setDocumentValue(document map[string]any, key string, value any) error {
	parent, rest, nested := strings.Cut(key, ".")
	for existing := range document {
		name := strings.ToLower(existing)
		if name == key {
			document[existing] = value
			return nil
		}
		if tail, ok := strings.CutPrefix(key, name+"."); ok && len(tail) <= len(rest) {
			parent, rest = existing, tail
		}
	}
	if !nested {
		document[key] = value
		return nil
	}
	child, ok := document[parent].(map[string]any)
	if !ok {
		if document[parent] != nil {
			return fmt.Errorf("config section %s must be a mapping", parent)
		}
		child = make(map[string]any)
		document[parent] = child
	}
	return setDocumentValue(child, rest, value)
}

// readConfig loads the config file when available, ignoring missing files.
func readConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) && v.ConfigFileUsed() == "" {
			return nil
		}

		return err
	}

	return nil
}

// applyOverrides merges CLI overrides into the config snapshot.
func applyOverrides(cfg *Config, ov Overrides) {
	if ov.Provider != "" {
		cfg.AI.Provider = ov.Provider
	}
	if ov.Model != "" {
		cfg.AI.Model = ov.Model
	}
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

	if ov.Prompt.Folder != "" {
		cfg.Prompt.Folder = ov.Prompt.Folder
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
