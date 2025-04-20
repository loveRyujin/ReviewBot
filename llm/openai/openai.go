package openai

import (
	"context"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/sashabaranov/go-openai"
)

var _ ai.TextGenerator = (*Client)(nil)

type Client struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
	topP        float32
}

type Response struct {
	Text       string
	TokenUsage ai.TokenUsage
}

func (cli *Client) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	// Implementation of the ChatCompletion method
	return nil, nil
}

type Config struct {
	ApiKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	TopP        float32
}

func (cfg *Config) New() (*Client, error) {
	return nil, nil
}
