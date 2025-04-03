package llm

import (
	"errors"

	"github.com/loveRyujin/ReviewBot/ai"
)

func GetModelClient(provider ai.Provider) (ai.TextGenerator, error) {
	switch provider {
	case ai.OpenAI:
		return nil, nil
	case ai.Anthropic:
		return nil, nil
	case ai.DeepSeek:
		return nil, nil
	default:
		return nil, errors.New("unsupported LLM provider")
	}
}
