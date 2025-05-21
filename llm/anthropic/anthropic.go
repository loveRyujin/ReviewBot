package anthropic

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client      *anthropic.Client
	model       string
	maxTokens   int
	temperature float32
	topP        float32
}

func (c *Client) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	resp, err := c.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:       anthropic.Model(c.model),
		MaxTokens:   c.maxTokens,
		Temperature: &c.temperature,
		TopP:        &c.topP,
		System:      "You're a helpful assistant.",
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewTextMessageContent(text),
				},
			},
		},
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}

		return nil, err
	}

	tokenUsage := ai.TokenUsage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		PromptTokensDetails: &openai.PromptTokensDetails{
			CachedTokens: resp.Usage.CacheReadInputTokens + resp.Usage.CacheCreationInputTokens,
		},
	}

	return &ai.Response{
		Text:       resp.Content[0].GetText(),
		TokenUsage: tokenUsage,
	}, nil
}

func (c *Client) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	color.Yellow("================Review Summary====================\n\n")

	resp, err := c.client.CreateMessagesStream(ctx, anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model:       anthropic.Model(c.model),
			MaxTokens:   c.maxTokens,
			Temperature: &c.temperature,
			TopP:        &c.topP,
			System:      "You're a helpful assistant.",
			Messages: []anthropic.Message{
				{
					Role: anthropic.RoleUser,
					Content: []anthropic.MessageContent{
						anthropic.NewTextMessageContent(text),
					},
				},
			},
		},
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			_ = handler(*data.Delta.Text)
		},
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages stream error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages stream error: %v\n", err)
		}

		return err
	}

	color.Yellow("\n==================================================")

	tokenUsage := ai.TokenUsage{
		PromptTokens:     resp.Usage.InputTokens,
		CompletionTokens: resp.Usage.OutputTokens,
		TotalTokens:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		PromptTokensDetails: &openai.PromptTokensDetails{
			CachedTokens: resp.Usage.CacheReadInputTokens + resp.Usage.CacheCreationInputTokens,
		},
	}
	color.Magenta(tokenUsage.String())

	return nil
}

type Config struct {
	BaseURL     string
	ApiKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	TopP        float32
}

func (cfg *Config) New(proxyCfg *proxy.Config) (*Client, error) {
	httpClient, _ := proxyCfg.New()
	options := []anthropic.ClientOption{
		anthropic.WithBaseURL(cfg.BaseURL),
		anthropic.WithHTTPClient(httpClient),
	}

	client := anthropic.NewClient(cfg.ApiKey, options...)

	return &Client{
		client:      client,
		model:       cfg.Model,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		topP:        cfg.TopP,
	}, nil
}
