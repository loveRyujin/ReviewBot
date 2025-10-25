package proxy

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToHTTPHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:     "empty headers",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:  "single header",
			input: []string{"Content-Type=application/json"},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "multiple headers",
			input: []string{
				"Content-Type=application/json",
				"Authorization=Bearer token123",
				"X-Custom-Header=value",
			},
			expected: map[string][]string{
				"Content-Type":    {"application/json"},
				"Authorization":   {"Bearer token123"},
				"X-Custom-Header": {"value"},
			},
		},
		{
			name: "headers with spaces",
			input: []string{
				" Content-Type = application/json ",
				" Authorization = Bearer token ",
			},
			expected: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer token"},
			},
		},
		{
			name: "invalid headers ignored",
			input: []string{
				"Content-Type=application/json",
				"InvalidHeader",
				"=EmptyKey",
				"EmptyValue=",
				" = ",
			},
			expected: map[string][]string{
				"Content-Type": {"application/json"},
			},
		},
		{
			name: "header with equals sign in value",
			input: []string{
				"X-Signature=abc=def=ghi",
			},
			expected: map[string][]string{
				"X-Signature": {"abc=def=ghi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toHTTPHeaders(tt.input)
			assert.Equal(t, http.Header(tt.expected), result)
		})
	}
}

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		wantErr   bool
		checkFunc func(t *testing.T, client *http.Client)
	}{
		{
			name: "default config",
			cfg: Config{
				Timeout: 30 * time.Second,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, client *http.Client) {
				assert.Equal(t, 30*time.Second, client.Timeout)
				assert.NotNil(t, client.Transport)
			},
		},
		{
			name: "config with HTTP proxy",
			cfg: Config{
				ProxyURL: "http://proxy.example.com:8080",
				Timeout:  60 * time.Second,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, client *http.Client) {
				assert.Equal(t, 60*time.Second, client.Timeout)
				assert.NotNil(t, client.Transport)
			},
		},
		{
			name: "config with custom headers",
			cfg: Config{
				Headers: []string{
					"X-Custom-Header=value1",
					"Authorization=Bearer token",
				},
				Timeout: 30 * time.Second,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, client *http.Client) {
				assert.NotNil(t, client.Transport)
				customTransport, ok := client.Transport.(*customTransport)
				assert.True(t, ok)
				assert.Equal(t, "value1", customTransport.headers.Get("X-Custom-Header"))
				assert.Equal(t, "Bearer token", customTransport.headers.Get("Authorization"))
			},
		},
		{
			name: "config with skip verify",
			cfg: Config{
				SkipVerify: true,
				Timeout:    30 * time.Second,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, client *http.Client) {
				assert.NotNil(t, client.Transport)
			},
		},
		{
			name: "invalid proxy URL",
			cfg: Config{
				ProxyURL: "://invalid-url",
				Timeout:  30 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.cfg.New()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, client)
			if tt.checkFunc != nil {
				tt.checkFunc(t, client)
			}
		})
	}
}

func TestCustomTransport_RoundTrip(t *testing.T) {
	t.Run("transport not set", func(t *testing.T) {
		ct := &customTransport{
			transport: nil,
			headers:   make(http.Header),
		}

		req, err := http.NewRequest("GET", "http://example.com", nil)
		require.NoError(t, err)

		resp, err := ct.RoundTrip(req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "origin transport is not set")
	})
}
