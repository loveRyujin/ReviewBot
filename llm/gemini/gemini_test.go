package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Struct(t *testing.T) {
	t.Run("config fields are correctly set", func(t *testing.T) {
		cfg := Config{
			BaseURL:     "https://custom.gemini.api",
			ApiKey:      "test-api-key",
			Model:       "gemini-pro",
			MaxTokens:   2048,
			Temperature: 0.9,
			TopP:        0.95,
		}

		assert.Equal(t, "https://custom.gemini.api", cfg.BaseURL)
		assert.Equal(t, "test-api-key", cfg.ApiKey)
		assert.Equal(t, "gemini-pro", cfg.Model)
		assert.Equal(t, 2048, cfg.MaxTokens)
		assert.Equal(t, float32(0.9), cfg.Temperature)
		assert.Equal(t, float32(0.95), cfg.TopP)
	})
}

func TestClient_Fields(t *testing.T) {
	t.Run("client structure", func(t *testing.T) {
		// Just verify the Client type can be created
		// Full integration tests would require actual API connection
		client := &Client{
			model:       "gemini-pro",
			maxTokens:   2048,
			temperature: 0.7,
			topP:        0.9,
		}

		assert.Equal(t, "gemini-pro", client.model)
		assert.Equal(t, 2048, client.maxTokens)
		assert.Equal(t, float32(0.7), client.temperature)
		assert.Equal(t, float32(0.9), client.topP)
	})
}
