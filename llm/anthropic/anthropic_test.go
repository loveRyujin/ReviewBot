package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/loveRyujin/ReviewBot/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T, serve http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(serve)
	t.Cleanup(server.Close)
	client, err := (&Config{
		BaseURL: server.URL + "/gateway/v1/", APIKey: "test-key", Model: "test-model", MaxTokens: 321,
	}).New(&proxy.Config{Headers: []string{"X-Gateway=reviewbot"}})
	require.NoError(t, err)
	return client
}

func checkRequest(t *testing.T, r *http.Request, stream bool) {
	t.Helper()
	assert.Equal(t, "POST", r.Method)
	assert.Equal(t, "/gateway/v1/messages", r.URL.Path)
	assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
	assert.Equal(t, "2023-06-01", r.Header.Get("anthropic-version"))
	assert.Equal(t, "reviewbot", r.Header.Get("X-Gateway"))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	var payload map[string]any
	assert.NoError(t, json.Unmarshal(body, &payload))
	assert.Equal(t, map[string]any{
		"model": "test-model", "max_tokens": float64(321), "stream": stream,
		"messages": []any{map[string]any{"role": "user", "content": "review me"}},
	}, payload)
}

func TestChatCompletion(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		checkRequest(t, r, false)
		io.WriteString(w, `{"content":[{"type":"text","text":"first"},{"type":"thinking","thinking":"hidden"},
			{"type":"text","text":" second"}],"usage":{"input_tokens":6,"output_tokens":8,
			"cache_read_input_tokens":3,"cache_creation_input_tokens":1}}`)
	})
	resp, err := client.ChatCompletion(t.Context(), "review me")
	require.NoError(t, err)
	assert.Equal(t, "first second", resp.Text)
	assert.Equal(t, 10, resp.TokenUsage.PromptTokens)
	assert.Equal(t, 8, resp.TokenUsage.CompletionTokens)
	assert.Equal(t, 18, resp.TokenUsage.TotalTokens)
	require.NotNil(t, resp.TokenUsage.PromptTokensDetails)
	assert.Equal(t, 3, resp.TokenUsage.PromptTokensDetails.CachedTokens)
}

func TestRequestErrors(t *testing.T) {
	for _, tc := range []struct {
		name       string
		status     int
		body, want string
	}{
		{"api error", 401, `{"error":{"type":"authentication_error","message":"invalid key"}}`, "authentication_error"},
		{"gateway error", 502, "upstream unavailable", "502 Bad Gateway"},
		{"invalid json", 200, "{", "decode anthropic message"},
		{"missing content", 200, "{}", "missing content"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				io.WriteString(w, tc.body)
			})
			resp, err := client.ChatCompletion(t.Context(), "review me")
			assert.Nil(t, resp)
			assert.ErrorContains(t, err, tc.want)
			if tc.status != 200 {
				assert.ErrorContains(t, client.StreamChatCompletion(t.Context(), "review me", nil), tc.want)
			}
		})
	}
}

func TestStreamCompletion(t *testing.T) {
	// Keep the connection open after message_stop to distinguish a stop from EOF.
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		checkRequest(t, r, true)
		assert.Equal(t, "text/event-stream", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "text/event-stream")
		io.WriteString(w, ": keepalive\r\n\r\n")
		for _, data := range []string{
			`{"type":"ping"}`,
			`{"type":"message_start","message":{"usage":{"input_tokens":12,"output_tokens":1}}}`,
			`{"type":"content_block_start","content_block":{"type":"text","text":"一"}}`,
			`{"type":"content_block_delta","delta":{"type":"thinking_delta","thinking":"hidden"}}`,
			`{"type":"content_block_delta","delta":{"type":"text_delta","text":"二"}}`,
			`{"type":"content_block_stop"}`,
			`{"type":"message_delta","usage":{"output_tokens":4}}`,
			`{"type":"message_stop"}`,
		} {
			io.WriteString(w, "data: "+data+"\r\n\r\n")
			w.(http.Flusher).Flush()
		}
		<-r.Context().Done()
	})
	var chunks []string
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	err := client.StreamChatCompletion(ctx, "review me", func(s string) error {
		chunks = append(chunks, s)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"一", "二"}, chunks)
}

func TestStreamUsage(t *testing.T) {
	body := strings.NewReader(`event: message_start
data: {"type":"message_start",
data: "message":{"usage":{"input_tokens":6,"cache_read_input_tokens":2,"output_tokens":1}}}

data: {"type":"message_delta","usage":{"output_tokens":4}}

data: {"type":"message_delta","usage":{"output_tokens":9}}

data: {"type":"message_stop"}

`)
	got, err := readStream(t.Context(), body, nil)
	require.NoError(t, err)
	assert.Equal(t, 8, got.PromptTokens)
	assert.Equal(t, 9, got.CompletionTokens)
	assert.Equal(t, 17, got.TotalTokens)
}

func TestStreamErrors(t *testing.T) {
	handlerErr := errors.New("handler stopped")
	for _, tc := range []struct{ name, event, want string }{
		{"api error", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`, "overloaded_error: busy"},
		{"early EOF", `{"type":"content_block_stop"}`, "before message_stop"},
		{"bad event", "{", "decode anthropic stream event"},
		{"callback", `{"type":"content_block_delta","delta":{"type":"text_delta","text":"x"}}`, handlerErr.Error()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "data: "+tc.event+"\n\n")
			})
			err := client.StreamChatCompletion(t.Context(), "review me", func(string) error { return handlerErr })
			assert.ErrorContains(t, err, tc.want)
			if tc.name == "callback" {
				assert.ErrorIs(t, err, handlerErr)
			}
		})
	}
}

func TestCancellation(t *testing.T) {
	for _, method := range []string{"chat", "stream"} {
		t.Run(method, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				io.Copy(io.Discard, r.Body)
				w.(http.Flusher).Flush()
				cancel()
				<-r.Context().Done()
			})
			var err error
			if method == "chat" {
				_, err = client.ChatCompletion(ctx, "review me")
			} else {
				err = client.StreamChatCompletion(ctx, "review me", nil)
			}
			assert.ErrorIs(t, err, context.Canceled)
		})
	}
}
