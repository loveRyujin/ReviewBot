package cmd

import (
	"context"
	"fmt"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/llm/anthropic"
	"github.com/loveRyujin/ReviewBot/llm/gemini"
	"github.com/loveRyujin/ReviewBot/llm/openai"
)

func NewOpenAIClient() (*openai.Client, error) {
	proxyCfg := globalConfig.ProxyConfig()
	return globalConfig.OpenAIConfig().New(proxyCfg)
}

func NewDeepSeekClient() (*openai.Client, error) {
	proxyCfg := globalConfig.ProxyConfig()
	return globalConfig.DeepSeekConfig().New(proxyCfg)
}

func NewAnthropicClient() (*anthropic.Client, error) {
	return globalConfig.AnthropicConfig().New(globalConfig.ProxyConfig())
}

func NewGeminiClient(ctx context.Context) (*gemini.Client, error) {
	proxyCfg := globalConfig.ProxyConfig()
	return globalConfig.GeminiConfig().New(ctx, proxyCfg)
}

func GetModelClient(ctx context.Context, provider ai.Provider) (ai.TextGenerator, error) {
	switch provider {
	case ai.OpenAI:
		return NewOpenAIClient()
	case ai.Anthropic:
		client, err := NewAnthropicClient()
		if err != nil {
			return nil, err
		}
		return client, nil
	case ai.DeepSeek:
		return NewDeepSeekClient()
	case ai.Gemini:
		return NewGeminiClient(ctx)
	default:
		return nil, fmt.Errorf("unsupported LLM provider %q (supported: openai, deepseek, gemini, anthropic)", provider)
	}
}
