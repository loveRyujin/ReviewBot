package config

import (
	"github.com/loveRyujin/ReviewBot/git"
	"github.com/loveRyujin/ReviewBot/llm/gemini"
	"github.com/loveRyujin/ReviewBot/llm/openai"
	"github.com/loveRyujin/ReviewBot/proxy"
)

// GitCommandConfig converts git settings into a git command configuration.
func (c *Config) GitCommandConfig() *git.Config {
	return &git.Config{
		DiffUnified:  c.Git.DiffUnified,
		ExcludedList: c.Git.ExcludedList,
		IsAmend:      c.Git.Amend,
	}
}

// OpenAIConfig returns the OpenAI configuration derived from settings.
func (c *Config) OpenAIConfig() *openai.Config {
	return &openai.Config{
		BaseURL:          c.AI.BaseURL,
		ApiKey:           c.AI.APIKey,
		Model:            c.AI.Model,
		MaxTokens:        c.AI.MaxTokens,
		Temperature:      c.AI.Temperature,
		TopP:             c.AI.TopP,
		PresencePenalty:  c.AI.PresencePenalty,
		FrequencyPenalty: c.AI.FrequencyPenalty,
	}
}

// DeepSeekConfig mirrors the OpenAI configuration for DeepSeek.
func (c *Config) DeepSeekConfig() *openai.Config {
	return &openai.Config{
		BaseURL:          c.AI.BaseURL,
		ApiKey:           c.AI.APIKey,
		Model:            c.AI.Model,
		MaxTokens:        c.AI.MaxTokens,
		Temperature:      c.AI.Temperature,
		TopP:             c.AI.TopP,
		PresencePenalty:  c.AI.PresencePenalty,
		FrequencyPenalty: c.AI.FrequencyPenalty,
	}
}

// GeminiConfig exposes Gemini-specific configuration values.
func (c *Config) GeminiConfig() *gemini.Config {
	return &gemini.Config{
		BaseURL:     c.AI.BaseURL,
		ApiKey:      c.AI.APIKey,
		Model:       c.AI.Model,
		MaxTokens:   c.AI.MaxTokens,
		Temperature: c.AI.Temperature,
		TopP:        c.AI.TopP,
	}
}

// ProxyConfig returns proxy settings for downstream clients.
func (c *Config) ProxyConfig() *proxy.Config {
	return &proxy.Config{
		ProxyURL:   c.Proxy.ProxyURL,
		SocksURL:   c.Proxy.SocksURL,
		Timeout:    c.Proxy.Timeout,
		Headers:    c.Proxy.Headers,
		SkipVerify: c.Proxy.SkipVerify,
	}
}
