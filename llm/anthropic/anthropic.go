package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/sashabaranov/go-openai"
)

var _ ai.TextGenerator = (*Client)(nil)

type Config struct {
	BaseURL   string
	APIKey    string
	Model     string
	MaxTokens int
}

type Client struct {
	config     Config
	httpClient *http.Client
	endpoint   string
}

func (cfg *Config) New(proxyCfg *proxy.Config) (*Client, error) {
	client, err := proxyCfg.New()
	if err != nil {
		return nil, err
	}
	// The API key header must stay on the configured endpoint.
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = "https://api.anthropic.com"
	}
	if !strings.HasSuffix(base, "/v1") {
		base += "/v1"
	}
	return &Client{config: *cfg, httpClient: client, endpoint: base + "/messages"}, nil
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) request(ctx context.Context, text string, stream bool) (*http.Response, error) {
	payload := struct {
		Model     string    `json:"model"`
		MaxTokens int       `json:"max_tokens"`
		Messages  []message `json:"messages"`
		Stream    bool      `json:"stream"`
	}{
		Model:     c.config.Model,
		MaxTokens: c.config.MaxTokens,
		Messages:  []message{{Role: "user", Content: text}},
		Stream:    stream,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("x-api-key", c.config.APIKey)
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
		if err != nil {
			return nil, fmt.Errorf("anthropic HTTP %s: %w", resp.Status, err)
		}
		return nil, fmt.Errorf("anthropic HTTP %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return resp, nil
}

type usage struct {
	Input      int `json:"input_tokens"`
	Output     int `json:"output_tokens"`
	CacheRead  int `json:"cache_read_input_tokens"`
	CacheWrite int `json:"cache_creation_input_tokens"`
}

func (u usage) tokens() ai.TokenUsage {
	input := u.Input + u.CacheRead + u.CacheWrite
	result := ai.TokenUsage{PromptTokens: input, CompletionTokens: u.Output, TotalTokens: input + u.Output}
	if u.CacheRead > 0 {
		result.PromptTokensDetails = &openai.PromptTokensDetails{CachedTokens: u.CacheRead}
	}
	return result
}

type block struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *Client) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	resp, err := c.request(ctx, text, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var msg struct {
		Content []block `json:"content"`
		Usage   usage   `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("decode anthropic message: %w", err)
	}
	if msg.Content == nil {
		return nil, fmt.Errorf("anthropic message missing content")
	}
	var out strings.Builder
	for _, part := range msg.Content {
		if part.Type == "text" {
			out.WriteString(part.Text)
		}
	}
	return &ai.Response{Text: out.String(), TokenUsage: msg.Usage.tokens()}, nil
}
