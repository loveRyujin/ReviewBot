package gemini

import (
	"context"

	"google.golang.org/genai"
)

type Client struct {
	client      *genai.Client
	model       string
	maxTokens   int
	temperature float32
	topP        float32
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

func (cfg *Config) New() (*Client, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.ApiKey,
		Backend: genai.BackendGeminiAPI,
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
