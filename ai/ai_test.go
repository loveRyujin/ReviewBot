package ai

import (
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestTokenUsage_String(t *testing.T) {
	tests := []struct {
		name     string
		usage    TokenUsage
		contains []string
	}{
		{
			name: "basic token usage",
			usage: TokenUsage{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
			},
			contains: []string{
				"Prompt tokens: 100",
				"Completion tokens: 50",
				"Total tokens: 150",
			},
		},
		{
			name: "with cached tokens",
			usage: TokenUsage{
				PromptTokens:     200,
				CompletionTokens: 100,
				TotalTokens:      300,
				PromptTokensDetails: &openai.PromptTokensDetails{
					CachedTokens: 50,
				},
			},
			contains: []string{
				"Prompt tokens: 200",
				"CachedTokens: 50",
				"Completion tokens: 100",
				"Total tokens: 300",
			},
		},
		{
			name: "with reasoning tokens",
			usage: TokenUsage{
				PromptTokens:     150,
				CompletionTokens: 75,
				TotalTokens:      225,
				CompletionTokensDetails: &openai.CompletionTokensDetails{
					ReasoningTokens: 25,
				},
			},
			contains: []string{
				"Prompt tokens: 150",
				"Completion tokens: 75",
				"ReasoningTokens: 25",
				"Total tokens: 225",
			},
		},
		{
			name: "with both cached and reasoning tokens",
			usage: TokenUsage{
				PromptTokens:     300,
				CompletionTokens: 150,
				TotalTokens:      450,
				PromptTokensDetails: &openai.PromptTokensDetails{
					CachedTokens: 100,
				},
				CompletionTokensDetails: &openai.CompletionTokensDetails{
					ReasoningTokens: 50,
				},
			},
			contains: []string{
				"Prompt tokens: 300",
				"CachedTokens: 100",
				"Completion tokens: 150",
				"ReasoningTokens: 50",
				"Total tokens: 450",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.usage.String()
			for _, substr := range tt.contains {
				assert.Contains(t, result, substr)
			}
		})
	}
}

func TestResponse(t *testing.T) {
	t.Run("create response with text and usage", func(t *testing.T) {
		resp := &Response{
			Text: "This is a test response",
			TokenUsage: TokenUsage{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
			},
		}

		assert.Equal(t, "This is a test response", resp.Text)
		assert.Equal(t, 100, resp.TokenUsage.PromptTokens)
		assert.Equal(t, 50, resp.TokenUsage.CompletionTokens)
		assert.Equal(t, 150, resp.TokenUsage.TotalTokens)
	})
}

func TestChunkHandler(t *testing.T) {
	t.Run("chunk handler can be called", func(t *testing.T) {
		var receivedChunks []string
		handler := func(chunk string) error {
			receivedChunks = append(receivedChunks, chunk)
			return nil
		}

		// Simulate streaming
		chunks := []string{"Hello", " ", "world", "!"}
		for _, chunk := range chunks {
			err := handler(chunk)
			assert.NoError(t, err)
		}

		assert.Equal(t, chunks, receivedChunks)
	})
}
