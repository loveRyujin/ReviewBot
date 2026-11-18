package cmd

import (
	"testing"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetModelClient(t *testing.T) {
	// Setup global config for testing
	globalConfig = &config.Config{
		AI: config.AIConfig{
			Provider:    "openai",
			APIKey:      "sk-test-key",
			Model:       "gpt-4",
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        1.0,
		},
		Proxy: config.ProxyConfig{},
	}

	tests := []struct {
		name     string
		provider ai.Provider
		wantErr  bool
	}{
		{
			name:     "openai provider",
			provider: ai.OpenAI,
			wantErr:  false,
		},
		{
			name:     "deepseek provider",
			provider: ai.DeepSeek,
			wantErr:  false,
		},
		{
			name:     "gemini provider",
			provider: ai.Gemini,
			wantErr:  false,
		},
		{
			name:     "anthropic provider (not implemented)",
			provider: ai.Anthropic,
			wantErr:  true,
		},
		{
			name:     "unsupported provider",
			provider: ai.Provider("unsupported"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := GetModelClient(t.Context(), tt.provider)
			if tt.wantErr {
				assert.Nil(t, client)
				assert.ErrorContains(t, err, "unsupported LLM provider")
				assert.ErrorContains(t, err, tt.provider.String())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestNewOpenAIClient(t *testing.T) {
	globalConfig = &config.Config{
		AI: config.AIConfig{
			Provider:    "openai",
			APIKey:      "sk-test-key",
			Model:       "gpt-4",
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        1.0,
		},
		Proxy: config.ProxyConfig{},
	}

	client, err := NewOpenAIClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewDeepSeekClient(t *testing.T) {
	globalConfig = &config.Config{
		AI: config.AIConfig{
			Provider:    "deepseek",
			APIKey:      "sk-test-key",
			Model:       "deepseek-chat",
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        1.0,
		},
		Proxy: config.ProxyConfig{},
	}

	client, err := NewDeepSeekClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewGeminiClient(t *testing.T) {
	globalConfig = &config.Config{
		AI: config.AIConfig{
			Provider:    "gemini",
			APIKey:      "test-gemini-key",
			Model:       "gemini-pro",
			MaxTokens:   2048,
			Temperature: 0.9,
			TopP:        0.95,
		},
		Proxy: config.ProxyConfig{},
	}

	client, err := NewGeminiClient(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
