package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultDiffUnified  = 3
	defaultMaxInputSize = 20 * 1024 * 1024
	defaultLang         = "en"
	defaultTimeout      = 30 * time.Second
	defaultProvider     = "openai"
	defaultModel        = "gpt-3.5-turbo"
)

var supportedLangs = map[string]struct{}{
	"en":    {},
	"zh-cn": {},
	"zh-tw": {},
	"ja":    {},
}

// Config holds all application settings grouped by domain.
type Config struct {
	Git     GitConfig     `mapstructure:"git" yaml:"git"`
	AI      AIConfig      `mapstructure:"ai" yaml:"ai"`
	Proxy   ProxyConfig   `mapstructure:"proxy" yaml:"proxy"`
	Prompt  PromptConfig  `mapstructure:"prompt" yaml:"prompt"`
	Runtime RuntimeConfig `mapstructure:"runtime" yaml:"runtime"`
}

// PromptConfig defines settings related to prompt templates.
type PromptConfig struct {
	Folder string `mapstructure:"folder" yaml:"folder"`
}

// GitConfig contains git-related options.
type GitConfig struct {
	DiffFile     string   `mapstructure:"diff_file" yaml:"diff_file"`
	MaxInputSize int      `mapstructure:"max_input_size" yaml:"max_input_size"`
	DiffUnified  int      `mapstructure:"diff_unified" yaml:"diff_unified"`
	ExcludedList []string `mapstructure:"exclude_list" yaml:"exclude_list"`
	Amend        bool     `mapstructure:"amend" yaml:"amend"`
	Lang         string   `mapstructure:"lang" yaml:"lang"`
}

// AIConfig describes AI provider settings.
type AIConfig struct {
	Provider         string  `mapstructure:"provider" yaml:"provider"`
	APIKey           string  `mapstructure:"api_key" yaml:"api_key"`
	BaseURL          string  `mapstructure:"base_url" yaml:"base_url"`
	Model            string  `mapstructure:"model" yaml:"model"`
	MaxTokens        int     `mapstructure:"max_tokens" yaml:"max_tokens"`
	Temperature      float32 `mapstructure:"temperature" yaml:"temperature"`
	TopP             float32 `mapstructure:"top_p" yaml:"top_p"`
	PresencePenalty  float32 `mapstructure:"presence_penalty" yaml:"presence_penalty"`
	FrequencyPenalty float32 `mapstructure:"frequency_penalty" yaml:"frequency_penalty"`
}

// ProxyConfig tracks proxy configuration fields.
type ProxyConfig struct {
	ProxyURL   string        `mapstructure:"proxy_url" yaml:"proxy_url"`
	SocksURL   string        `mapstructure:"socks_url" yaml:"socks_url"`
	Timeout    time.Duration `mapstructure:"timeout" yaml:"timeout"`
	Headers    []string      `mapstructure:"headers" yaml:"headers"`
	SkipVerify bool          `mapstructure:"skip_verify" yaml:"skip_verify"`
}

// RuntimeConfig stores command runtime options.
type RuntimeConfig struct {
	Review ReviewRuntime `mapstructure:"review" yaml:"review"`
	Commit CommitRuntime `mapstructure:"commit" yaml:"commit"`
}

// ReviewRuntime captures review command runtime flags.
type ReviewRuntime struct {
	Mode       string `mapstructure:"mode" yaml:"mode"`
	Stream     bool   `mapstructure:"stream" yaml:"stream"`
	DiffFile   string `mapstructure:"diff_file" yaml:"diff_file"`
	MaxInput   int    `mapstructure:"max_input_size" yaml:"max_input_size"`
	OutputLang string `mapstructure:"output_lang" yaml:"output_lang"`
}

// CommitRuntime captures commit command runtime flags.
type CommitRuntime struct {
	Preview    bool   `mapstructure:"preview" yaml:"preview"`
	OutputLang string `mapstructure:"output_lang" yaml:"output_lang"`
}

// NewDefault returns configuration populated with default values.
func NewDefault() *Config {
	return &Config{
		Git: GitConfig{
			DiffUnified:  defaultDiffUnified,
			MaxInputSize: defaultMaxInputSize,
			Lang:         defaultLang,
		},
		AI: AIConfig{
			Provider:         defaultProvider,
			Model:            defaultModel,
			MaxTokens:        1000,
			Temperature:      0.7,
			TopP:             1.0,
			PresencePenalty:  0.5,
			FrequencyPenalty: 0.5,
		},
		Prompt: PromptConfig{},
		Proxy: ProxyConfig{
			Timeout: defaultTimeout,
		},
		Runtime: RuntimeConfig{},
	}
}

// setDefaults registers configuration defaults with viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("git.diff_unified", defaultDiffUnified)
	v.SetDefault("git.max_input_size", defaultMaxInputSize)
	v.SetDefault("git.lang", defaultLang)

	v.SetDefault("ai.provider", defaultProvider)
	v.SetDefault("ai.model", defaultModel)
	v.SetDefault("ai.max_tokens", 1000)
	v.SetDefault("ai.temperature", 0.7)
	v.SetDefault("ai.top_p", 1.0)
	v.SetDefault("ai.presence_penalty", 0.5)
	v.SetDefault("ai.frequency_penalty", 0.5)

	v.SetDefault("proxy.timeout", defaultTimeout)

	v.SetDefault("prompt.folder", "")
}
