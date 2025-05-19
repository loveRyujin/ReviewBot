package cmd

import (
	"errors"
	"time"

	"github.com/loveRyujin/ReviewBot/git"
	"github.com/loveRyujin/ReviewBot/llm/gemini"
	"github.com/loveRyujin/ReviewBot/llm/openai"
	"github.com/loveRyujin/ReviewBot/pkg/command"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/spf13/viper"
)

var (
	ServerOption *ServerOptions
)

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
		BaseURL:          s.AiOptions.BaseURL,
		ApiKey:           s.AiOptions.ApiKey,
		Model:            s.AiOptions.Model,
		MaxTokens:        s.AiOptions.MaxTokens,
		Temperature:      s.AiOptions.Temperature,
		TopP:             s.AiOptions.TopP,
		PresencePenalty:  s.AiOptions.PresencePenalty,
		FrequencyPenalty: s.AiOptions.FrequencyPenalty,
	}
}

func (s *ServerOptions) DeepSeekConfig() *openai.Config {
	return &openai.Config{
		BaseURL:          s.AiOptions.BaseURL,
		ApiKey:           s.AiOptions.ApiKey,
		Model:            s.AiOptions.Model,
		MaxTokens:        s.AiOptions.MaxTokens,
		Temperature:      s.AiOptions.Temperature,
		TopP:             s.AiOptions.TopP,
		PresencePenalty:  s.AiOptions.PresencePenalty,
		FrequencyPenalty: s.AiOptions.FrequencyPenalty,
	}
}

func (s *ServerOptions) GeminiConfig() *gemini.Config {
	return &gemini.Config{
		BaseURL:          s.AiOptions.BaseURL,
		ApiKey:           s.AiOptions.ApiKey,
		Model:            s.AiOptions.Model,
		MaxTokens:        s.AiOptions.MaxTokens,
		Temperature:      s.AiOptions.Temperature,
		TopP:             s.AiOptions.TopP,
		PresencePenalty:  s.AiOptions.PresencePenalty,
		FrequencyPenalty: s.AiOptions.FrequencyPenalty,
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
	if !command.IsCommandAvailable("git") {
		return errors.New("git command not found in your system PATH, Please install git")
	}

	if diffUnifiedLines != 3 {
		s.GitOptions.DiffUnified = diffUnifiedLines
	}

	if len(excludedList) > 0 {
		s.GitOptions.ExcludedList = append(s.GitOptions.ExcludedList, excludedList...)
	}

	if diffFile != "" {
		s.GitOptions.DiffFile = diffFile
	}

	if maxInputSize != 20*1024*1024 {
		s.GitOptions.MaxInputSize = maxInputSize
	}

	if outputLang != "en" {
		s.GitOptions.Lang = outputLang
	}

	return nil
}

type GitOptions struct {
	DiffFile     string   `mapstructure:"diff_file"`
	MaxInputSize int      `mapstructure:"max_input_size"`
	DiffUnified  int      `mapstructure:"diff_unified"`
	ExcludedList []string `mapstructure:"exclude_list"`
	Amend        bool     `mapstructure:"amend"`
	Lang         string   `mapstructure:"lang"`
}

func NewGitOptions() *GitOptions {
	return &GitOptions{
		MaxInputSize: 20 * 1024 * 1024,
		DiffUnified:  3,
		ExcludedList: []string{},
		Amend:        false,
		Lang:         "en",
	}
}

func (g *GitOptions) validate() error {
	if g.DiffUnified < 0 {
		return errors.New("diff_unified must be a non-negative integer")
	}

	if g.MaxInputSize <= 0 {
		return errors.New("max_input_size must be a positive integer")
	}

	return nil
}

type AiOptions struct {
	Provider         string  `mapstructure:"provider"`
	ApiKey           string  `mapstructure:"api_key"`
	BaseURL          string  `mapstructure:"base_url"`
	Model            string  `mapstructure:"model"`
	MaxTokens        int     `mapstructure:"max_tokens"`
	Temperature      float32 `mapstructure:"temperature"`
	TopP             float32 `mapstructure:"top_p"`
	PresencePenalty  float32 `mapstructure:"presence_penalty"`
	FrequencyPenalty float32 `mapstructure:"frequency_penalty"`
}

func NewAiOptions() *AiOptions {
	return &AiOptions{
		Provider:         "openai",
		ApiKey:           "xxxxxx",
		Model:            "gpt-3.5-turbo",
		MaxTokens:        1000,
		Temperature:      0.7,
		TopP:             1.0,
		PresencePenalty:  0.5,
		FrequencyPenalty: 0.5,
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
