package gemini

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

var _ ai.TextGenerator = (*Client)(nil)

type Client struct {
	client      *genai.Client
	model       string
	maxTokens   int
	temperature float32
	topP        float32
}

func (c *Client) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	config := &genai.GenerateContentConfig{
		Temperature:       &c.temperature,
		TopP:              &c.topP,
		MaxOutputTokens:   int32(c.maxTokens),
		SystemInstruction: genai.NewContentFromText("You are a helpful assistant.", genai.RoleUser),
	}

	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(text), config)
	if err != nil {
		return nil, err
	}

	result := &ai.Response{Text: resp.Text()}
	if resp.UsageMetadata != nil {
		result.TokenUsage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		result.TokenUsage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		result.TokenUsage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
		result.TokenUsage.PromptTokensDetails = &openai.PromptTokensDetails{
			CachedTokens: int(resp.UsageMetadata.CachedContentTokenCount),
		}
	}

	return result, nil
}

func (c *Client) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	config := &genai.GenerateContentConfig{
		Temperature:       &c.temperature,
		TopP:              &c.topP,
		MaxOutputTokens:   int32(c.maxTokens),
		SystemInstruction: genai.NewContentFromText("You are a helpful assistant.", genai.RoleUser),
	}

	stream := c.client.Models.GenerateContentStream(ctx, c.model, genai.Text(text), config)

	color.Yellow("================Review Summary====================\n\n")

	tokenUsage := &ai.TokenUsage{}
	for chunk, err := range stream {
		if err != nil {
			if errors.Is(err, io.EOF) {
				color.Yellow("\n==================================================")
				break
			}
			return err
		}

		part := chunk.Candidates[0].Content.Parts[0]
		if err := handler(part.Text); err != nil {
			return err
		}

		if chunk.UsageMetadata != nil {
			tokenUsage.PromptTokens = int(chunk.UsageMetadata.PromptTokenCount)
			tokenUsage.CompletionTokens = int(chunk.UsageMetadata.CandidatesTokenCount)
			tokenUsage.TotalTokens = int(chunk.UsageMetadata.TotalTokenCount)
			tokenUsage.PromptTokensDetails = &openai.PromptTokensDetails{
				CachedTokens: int(chunk.UsageMetadata.CachedContentTokenCount),
			}
		}
	}
	color.Magenta(tokenUsage.String())

	return nil
}

type Config struct {
	ApiKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	TopP        float32
}

func (cfg *Config) New(proxyCfg *proxy.Config) (*Client, error) {
	httpClient, _ := proxyCfg.New()
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		HTTPClient: httpClient,
		APIKey:     cfg.ApiKey,
		Backend:    genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client:      client,
		model:       cfg.Model,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		topP:        cfg.TopP,
	}, nil
}
