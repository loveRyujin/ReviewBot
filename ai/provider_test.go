package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_String(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		expected string
	}{
		{
			name:     "OpenAI provider",
			provider: OpenAI,
			expected: "openai",
		},
		{
			name:     "Anthropic provider",
			provider: Anthropic,
			expected: "anthropic",
		},
		{
			name:     "DeepSeek provider",
			provider: DeepSeek,
			expected: "deepseek",
		},
		{
			name:     "Gemini provider",
			provider: Gemini,
			expected: "gemini",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.provider.String())
		})
	}
}

func TestProviderConstants(t *testing.T) {
	// Verify that provider constants are defined correctly
	assert.Equal(t, Provider("openai"), OpenAI)
	assert.Equal(t, Provider("anthropic"), Anthropic)
	assert.Equal(t, Provider("deepseek"), DeepSeek)
	assert.Equal(t, Provider("gemini"), Gemini)
}
