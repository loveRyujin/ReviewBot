package cmd

import (
	"errors"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/llm/gemini"
	"github.com/loveRyujin/ReviewBot/llm/openai"
)

func NewOpenAIClient() (*openai.Client, error) {
	proxyCfg := ServerOption.ProxyConfig()
	return ServerOption.OpenaiConfig().New(proxyCfg)
}

func NewDeepSeekClient() (*openai.Client, error) {
	proxyCfg := ServerOption.ProxyConfig()
	return ServerOption.DeepSeekConfig().New(proxyCfg)
}

func NewGeminiClient() (*gemini.Client, error) {
	proxyCfg := ServerOption.ProxyConfig()
	return ServerOption.GeminiConfig().New(proxyCfg)
}

func GetModelClient(provider ai.Provider) (ai.TextGenerator, error) {
	switch provider {
	case ai.OpenAI:
		return NewOpenAIClient()
	case ai.Anthropic:
		return nil, nil
	case ai.DeepSeek:
		return NewDeepSeekClient()
	case ai.Gemini:
		return NewGeminiClient()
	default:
		return nil, errors.New("unsupported LLM provider")
	}
}
