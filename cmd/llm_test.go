package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:     "anthropic provider",
			provider: ai.Anthropic,
			wantErr:  false,
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

func TestAnthropicFactoryFromConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/messages", r.URL.Path)
		assert.Equal(t, "test-anthropic-key", r.Header.Get("x-api-key"))
		assert.Equal(t, "factory", r.Header.Get("X-Test"))
		var payload struct {
			Model     string
			MaxTokens int `json:"max_tokens"`
		}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "claude-test", payload.Model)
		assert.Equal(t, 2048, payload.MaxTokens)
		io.WriteString(w, `{"content":[{"type":"text","text":"review"}]}`)
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "reviewbot.yaml")
	settings := `ai:
  provider: anthropic
  api_key: ""
  model: claude-test
  max_tokens: 2048
  base_url: ` + server.URL + `
proxy:
  headers: ["X-Test=factory"]
`
	require.NoError(t, os.WriteFile(path, []byte(settings), 0600))
	t.Setenv("REVIEWBOT_AI_API_KEY", "test-anthropic-key")
	cfg, err := config.Load(config.LoadOptions{ExplicitPath: path, EnvPrefix: "REVIEWBOT"})
	require.NoError(t, err)
	previous := globalConfig
	t.Cleanup(func() { globalConfig = previous })
	globalConfig = cfg
	client, err := GetModelClient(t.Context(), ai.Provider(cfg.AI.Provider))
	require.NoError(t, err)
	resp, err := client.ChatCompletion(t.Context(), "test")
	require.NoError(t, err)
	assert.Equal(t, "review", resp.Text)

	cfg.Proxy.ProxyURL = "://invalid"
	client, err = GetModelClient(t.Context(), ai.Anthropic)
	require.ErrorContains(t, err, "invalid proxy URL")
	assert.True(t, client == nil)
}
