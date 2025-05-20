package openai

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/sashabaranov/go-openai"
)

var _ ai.TextGenerator = (*Client)(nil)

type Client struct {
	client           *openai.Client
	model            string
	maxTokens        int
	temperature      float32
	topP             float32
	PresencePenalty  float32
	FrequencyPenalty float32
}

type Response struct {
	Text       string
	TokenUsage openai.Usage
}

func (c *Client) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	resp, err := c.chatCompletion(ctx, text)
	if err != nil {
		return nil, err
	}

	return &ai.Response{
		Text: resp.Text,
		TokenUsage: ai.TokenUsage{
			PromptTokens:            resp.TokenUsage.PromptTokens,
			CompletionTokens:        resp.TokenUsage.CompletionTokens,
			TotalTokens:             resp.TokenUsage.TotalTokens,
			CompletionTokensDetails: resp.TokenUsage.CompletionTokensDetails,
		},
	}, nil
}

func (c *Client) chatCompletion(ctx context.Context, text string) (*Response, error) {
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		PresencePenalty:  c.PresencePenalty,
		FrequencyPenalty: c.FrequencyPenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:       resp.Choices[0].Message.Content,
		TokenUsage: resp.Usage,
	}, nil
}

// StreamChatCompletion streams the chat completion response.
func (c *Client) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	req := openai.ChatCompletionRequest{
		Model:            c.model,
		MaxTokens:        c.maxTokens,
		Temperature:      c.temperature,
		TopP:             c.topP,
		PresencePenalty:  c.PresencePenalty,
		FrequencyPenalty: c.FrequencyPenalty,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		Stream: true,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	color.Yellow("================Review Summary====================" + "\n\n")

	tokenUsage := ai.TokenUsage{}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				color.Yellow("\n" + "==================================================")
				break
			}
			return err
		}

		chunk := resp.Choices[0].Delta.Content
		if err := handler(chunk); err != nil {
			return err
		}

		if resp.Usage != nil {
			tokenUsage.PromptTokens = resp.Usage.PromptTokens
			tokenUsage.CompletionTokens = resp.Usage.CompletionTokens
			tokenUsage.TotalTokens = resp.Usage.TotalTokens
			tokenUsage.PromptTokensDetails = resp.Usage.PromptTokensDetails
			tokenUsage.CompletionTokensDetails = resp.Usage.CompletionTokensDetails
		}
	}
	color.Magenta(tokenUsage.String())

	return nil
}

type Config struct {
	BaseURL          string
	ApiKey           string
	Model            string
	MaxTokens        int
	Temperature      float32
	TopP             float32
	PresencePenalty  float32
	FrequencyPenalty float32
}

func (cfg *Config) New(proxyCfg *proxy.Config) (*Client, error) {
	c := openai.DefaultConfig(cfg.ApiKey)
	if cfg.BaseURL != "" {
		c.BaseURL = cfg.BaseURL
	}

	httpClient, _ := proxyCfg.New()
	c.HTTPClient = httpClient

	client := openai.NewClientWithConfig(c)
	return &Client{
		client:           client,
		model:            cfg.Model,
		maxTokens:        cfg.MaxTokens,
		temperature:      cfg.Temperature,
		topP:             cfg.TopP,
		PresencePenalty:  cfg.PresencePenalty,
		FrequencyPenalty: cfg.FrequencyPenalty,
	}, nil
}
