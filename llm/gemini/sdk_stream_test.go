package gemini

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

type streamTransport func(*http.Request) (*http.Response, error)

func (f streamTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type streamBody struct {
	*strings.Reader
	onEnd  func() error
	closed bool
}

func (b *streamBody) Read(p []byte) (int, error) {
	n, err := b.Reader.Read(p)
	if err == io.EOF && b.onEnd != nil {
		return n, b.onEnd()
	}
	return n, err
}

func (b *streamBody) Close() error {
	b.closed = true
	return nil
}

func TestStreamChatCompletionWithSDK(t *testing.T) {
	const events = "data: {}\n\n" +
		"data: {\"candidates\":[]}\n\n" +
		"data: {\"candidates\":[{\"finishReason\":\"STOP\"}]}\n\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[]}}]}\n\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"internal\",\"thought\":true},{\"text\":\"first \"},{\"text\":\"second\"}]}},{\"content\":{\"parts\":[{\"text\":\"alternative\"}]}}]}\n\n" +
		"data: {\"usageMetadata\":{\"promptTokenCount\":4,\"candidatesTokenCount\":2,\"totalTokenCount\":6}}\n\n"

	for _, tt := range []struct {
		name       string
		tail       string
		cancelRead bool
	}{
		{name: "normal completion"},
		{name: "upstream error", tail: "{\"error\":{\"code\":503,\"message\":\"unavailable\",\"status\":\"UNAVAILABLE\"}}\n\n"},
		{name: "canceled while reading", cancelRead: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			body := &streamBody{Reader: strings.NewReader(events + tt.tail)}
			transport := streamTransport(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, http.MethodPost, req.Method)
				require.Equal(t, "/v1beta/models/test-model:streamGenerateContent", req.URL.Path)
				require.Equal(t, "sse", req.URL.Query().Get("alt"))
				require.Same(t, ctx, req.Context())
				if tt.cancelRead {
					body.onEnd = func() error {
						cancel()
						return req.Context().Err()
					}
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"text/event-stream"}},
					Body:       body,
				}, nil
			})
			sdk, err := genai.NewClient(ctx, &genai.ClientConfig{
				APIKey:     "test-key",
				Backend:    genai.BackendGeminiAPI,
				HTTPClient: &http.Client{Transport: transport},
				HTTPOptions: genai.HTTPOptions{
					BaseURL:    "https://gemini.test",
					APIVersion: "v1beta",
				},
			})
			require.NoError(t, err)
			client := &Client{models: sdk.Models, model: "test-model"}
			var received []string
			err = client.StreamChatCompletion(ctx, "review", func(s string) error {
				received = append(received, s)
				return nil
			})

			switch {
			case tt.cancelRead:
				require.ErrorIs(t, err, context.Canceled)
			case tt.tail != "":
				var apiErr genai.APIError
				require.ErrorAs(t, err, &apiErr)
				require.Equal(t, http.StatusServiceUnavailable, apiErr.Code)
			default:
				require.NoError(t, err)
			}
			require.Equal(t, []string{"first ", "second"}, received)
			require.True(t, body.closed)
		})
	}
}
