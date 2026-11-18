package gemini

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

var _ ai.TextGenerator = (*Client)(nil)

type modelAPI interface {
	GenerateContent(context.Context, string, []*genai.Content, *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error)
	GenerateContentStream(context.Context, string, []*genai.Content, *genai.GenerateContentConfig) iter.Seq2[*genai.GenerateContentResponse, error]
}

type Client struct {
	models      modelAPI
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

	resp, err := c.models.GenerateContent(ctx, c.model, genai.Text(text), config)
	if err != nil {
		return nil, err
	}

	result := &ai.Response{Text: resp.Text()}
	result.TokenUsage = usageFromMetadata(resp.UsageMetadata)

	return result, nil
}

func usageFromMetadata(metadata *genai.GenerateContentResponseUsageMetadata) ai.TokenUsage {
	if metadata == nil {
		return ai.TokenUsage{}
	}
	return ai.TokenUsage{
		PromptTokens:        int(metadata.PromptTokenCount),
		CompletionTokens:    int(metadata.CandidatesTokenCount),
		TotalTokens:         int(metadata.TotalTokenCount),
		PromptTokensDetails: &openai.PromptTokensDetails{CachedTokens: int(metadata.CachedContentTokenCount)},
	}
}

func (c *Client) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	config := &genai.GenerateContentConfig{
		Temperature:       &c.temperature,
		TopP:              &c.topP,
		MaxOutputTokens:   int32(c.maxTokens),
		SystemInstruction: genai.NewContentFromText("You are a helpful assistant.", genai.RoleUser),
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	stream := c.models.GenerateContentStream(ctx, c.model, genai.Text(text), config)

	color.Yellow("================Review Summary====================\n\n")

	var tokenUsage ai.TokenUsage
	var wroteText bool
	var blockReason string
	for chunk, err := range stream {
		if err != nil {
			if errors.Is(err, io.EOF) {
				color.Yellow("\n==================================================")
				break
			}
			return err
		}

		if err := ctx.Err(); err != nil {
			return err
		}
		if chunk == nil {
			continue
		}
		if chunk.UsageMetadata != nil {
			tokenUsage = usageFromMetadata(chunk.UsageMetadata)
		}
		if !wroteText && blockReason == "" {
			blockReason = streamBlockReason(chunk)
		}
		if len(chunk.Candidates) == 0 || chunk.Candidates[0] == nil || chunk.Candidates[0].Content == nil {
			continue
		}
		for _, part := range chunk.Candidates[0].Content.Parts {
			if part == nil || part.Thought || part.Text == "" {
				continue
			}
			if err := handler(part.Text); err != nil {
				return err
			}
			wroteText = true
			if err := ctx.Err(); err != nil {
				return err
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if !wroteText && blockReason != "" {
		return fmt.Errorf("Gemini blocked the response: %s", blockReason)
	}
	color.Magenta(tokenUsage.String())

	return nil
}

func streamBlockReason(chunk *genai.GenerateContentResponse) string {
	if feedback := chunk.PromptFeedback; feedback != nil && feedback.BlockReason != "" {
		return string(feedback.BlockReason)
	}
	for _, candidate := range chunk.Candidates {
		if candidate == nil {
			continue
		}
		switch candidate.FinishReason {
		case genai.FinishReasonSafety, genai.FinishReasonRecitation, genai.FinishReasonBlocklist,
			genai.FinishReasonProhibitedContent, genai.FinishReasonSPII:
			return string(candidate.FinishReason)
		}
	}
	return ""
}

type Config struct {
	BaseURL     string
	ApiKey      string
	Model       string
	MaxTokens   int
	Temperature float32
	TopP        float32
}

func (cfg *Config) New(ctx context.Context, proxyCfg *proxy.Config) (*Client, error) {
	httpClient, _ := proxyCfg.New()
	httpOption := genai.HTTPOptions{}
	if cfg.BaseURL != "" {
		httpOption.BaseURL = cfg.BaseURL
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPClient:  httpClient,
		HTTPOptions: httpOption,
		APIKey:      cfg.ApiKey,
		Backend:     genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		models:      client.Models,
		model:       cfg.Model,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		topP:        cfg.TopP,
	}, nil
}
