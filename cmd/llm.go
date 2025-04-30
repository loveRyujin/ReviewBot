package cmd

import (
	"errors"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/llm/openai"
)

func NewOpenAIClient() (*openai.Client, error) {
	proxyCfg := ServerOption.ProxyConfig()
	return ServerOption.OpenaiConfig().New(proxyCfg)
}

func GetModelClient(provider ai.Provider) (ai.TextGenerator, error) {
	switch provider {
	case ai.OpenAI:
		return NewOpenAIClient()
	case ai.Anthropic:
		return nil, nil
	case ai.DeepSeek:
		return nil, nil
	default:
		return nil, errors.New("unsupported LLM provider")
	}
}
