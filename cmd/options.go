package cmd

import (
	"errors"
	"sync"
	"time"

	"github.com/loveRyujin/ReviewBot/git"
	"github.com/loveRyujin/ReviewBot/llm/openai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/loveRyujin/ReviewBot/util"
	"github.com/spf13/viper"
)

var (
	ServerOption *ServerOptions
	once         sync.Once
)

func init() {
	once.Do(func() {
		ServerOption = NewServerOptions()
	})
}

type ServerOptions struct {
	GitOptions   *GitOptions   `mapstructure:"git"`
	AiOptions    *AiOptions    `mapstructure:"ai"`
	ProxyOptions *ProxyOptions `mapstructure:"proxy"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		GitOptions:   NewGitOptions(),
		AiOptions:    NewAiOptions(),
		ProxyOptions: NewProxyOptions(),
	}
}

func (s *ServerOptions) Initialize() error {
	if err := s.applyCfg(); err != nil {
		return err
	}

	if err := viper.Unmarshal(s); err != nil {
		return err
	}

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

func (s *ServerOptions) validate() error {
	if err := s.GitOptions.validate(); err != nil {
		return err
	}

	if err := s.AiOptions.validate(); err != nil {
		return err
	}

	return nil
}

func (s *ServerOptions) GitConfig() *git.Config {
	return &git.Config{
		DiffUnified:  s.GitOptions.DiffUnified,
		ExcludedList: s.GitOptions.ExcludedList,
		IsAmend:      s.GitOptions.Amend,
	}
}

func (s *ServerOptions) OpenaiConfig() *openai.Config {
	return &openai.Config{
		ApiKey:      s.AiOptions.ApiKey,
		Model:       s.AiOptions.Model,
		MaxTokens:   s.AiOptions.MaxTokens,
		Temperature: s.AiOptions.Temperature,
		TopP:        s.AiOptions.TopP,
	}
}

func (s *ServerOptions) ProxyConfig() *proxy.Config {
	return &proxy.Config{
		ProxyURL:   s.ProxyOptions.ProxyURL,
		SocksURL:   s.ProxyOptions.SocksURL,
		Timeout:    s.ProxyOptions.Timeout,
		Headers:    s.ProxyOptions.Headers,
		SkipVerify: s.ProxyOptions.SkipVerify,
	}
}

// applyCfg applies user-provided configuration from the command line.
// It ensures the "git" command is available in the system PATH and updates
// configuration settings such as the number of unified lines for diffs
// and the list of excluded items.
func (s *ServerOptions) applyCfg() error {
	if !util.IsCommandAvailable("git") {
		return errors.New("git command not found in your system PATH, Please install git")
	}

	if diffUnifiedLines != 3 {
		viper.Set("git.diff_unified", diffUnifiedLines)
	}

	if len(excludedList) > 0 {
		viper.Set("git.exclude_list", excludedList)
	}

	return nil
}

type GitOptions struct {
	DiffUnified  int      `mapstructure:"diff_unified"`
	ExcludedList []string `mapstructure:"exclude_list"`
	Amend        bool     `mapstructure:"amend"`
}

func NewGitOptions() *GitOptions {
	return &GitOptions{
		DiffUnified:  3,
		ExcludedList: []string{},
		Amend:        false,
	}
}

func (g *GitOptions) validate() error {
	if g.DiffUnified < 0 {
		return errors.New("diff_unified must be a non-negative integer")
	}

	return nil
}

type AiOptions struct {
	Provider         string  `mapstructure:"provider"`
	ApiKey           string  `mapstructure:"api_key"`
	Model            string  `mapstructure:"model"`
	MaxTokens        int     `mapstructure:"max_tokens"`
	Temperature      float32 `mapstructure:"temperature"`
	TopP             float32 `mapstructure:"top_p"`
	PresencePenalty  float32 `mapstructure:"presence_penalty"`
	FrequencyPenalty float32 `mapstructure:"frequency_penalty"`
}

func NewAiOptions() *AiOptions {
	return &AiOptions{
		Provider:    "openai",
		ApiKey:      "xxxxxx",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        1.0,
	}
}

func (a *AiOptions) validate() error {
	if a.Provider == "" {
		return errors.New("provider cannot be empty")
	}

	if a.ApiKey == "" {
		return errors.New("api_key cannot be empty")
	}

	return nil
}

type ProxyOptions struct {
	ProxyURL   string        `mapstructure:"proxy_url"`
	SocksURL   string        `mapstructure:"socks_url"`
	Timeout    time.Duration `mapstructure:"timeout"`
	Headers    []string      `mapstructure:"headers"`
	SkipVerify bool          `mapstructure:"skip_verify"`
}

func NewProxyOptions() *ProxyOptions {
	return &ProxyOptions{
		ProxyURL:   "",
		SocksURL:   "",
		Timeout:    30 * time.Second,
		Headers:    []string{},
		SkipVerify: false,
	}
}
