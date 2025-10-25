package openai

import (
	"testing"

	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		proxyCfg  *proxy.Config
		wantModel string
		wantErr   bool
	}{
		{
			name: "valid config with default base URL",
			cfg: Config{
				ApiKey:      "sk-test-key",
				Model:       "gpt-4",
				MaxTokens:   2000,
				Temperature: 0.7,
				TopP:        1.0,
			},
			proxyCfg:  &proxy.Config{},
			wantModel: "gpt-4",
			wantErr:   false,
		},
		{
			name: "valid config with custom base URL",
			cfg: Config{
				BaseURL:     "https://api.custom.com/v1",
				ApiKey:      "sk-test-key",
				Model:       "gpt-3.5-turbo",
				MaxTokens:   1000,
				Temperature: 0.5,
				TopP:        0.9,
			},
			proxyCfg:  &proxy.Config{},
			wantModel: "gpt-3.5-turbo",
			wantErr:   false,
		},
		{
			name: "config with all parameters",
			cfg: Config{
				ApiKey:           "sk-test-key",
				Model:            "gpt-4-turbo",
				MaxTokens:        4000,
				Temperature:      0.8,
				TopP:             0.95,
				PresencePenalty:  0.5,
				FrequencyPenalty: 0.3,
			},
			proxyCfg:  &proxy.Config{},
			wantModel: "gpt-4-turbo",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.cfg.New(tt.proxyCfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, client)
			assert.NotNil(t, client.client)
			assert.Equal(t, tt.wantModel, client.model)
			assert.Equal(t, tt.cfg.MaxTokens, client.maxTokens)
			assert.Equal(t, tt.cfg.Temperature, client.temperature)
			assert.Equal(t, tt.cfg.TopP, client.topP)
		})
	}
}

func TestClient_Fields(t *testing.T) {
	cfg := Config{
		ApiKey:           "sk-test",
		Model:            "gpt-4",
		MaxTokens:        2000,
		Temperature:      0.7,
		TopP:             1.0,
		PresencePenalty:  0.5,
		FrequencyPenalty: 0.3,
	}

	client, err := cfg.New(&proxy.Config{})
	require.NoError(t, err)

	assert.Equal(t, "gpt-4", client.model)
	assert.Equal(t, 2000, client.maxTokens)
	assert.Equal(t, float32(0.7), client.temperature)
	assert.Equal(t, float32(1.0), client.topP)
	assert.Equal(t, float32(0.5), client.PresencePenalty)
	assert.Equal(t, float32(0.3), client.FrequencyPenalty)
}
