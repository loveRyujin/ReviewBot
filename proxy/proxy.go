package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/http/httpguts"
	"golang.org/x/net/proxy"
)

type customTransport struct {
	transport http.RoundTripper
	headers   http.Header
}

// RoundTrip adds custom headers and forwards HTTP requests
func (c *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.transport == nil {
		return nil, fmt.Errorf("origin transport is not set")
	}

	// Set custom headers
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return c.transport.RoundTrip(req)
}

// ParseHeaders splits a comma-separated list only where a new HTTP field starts.
func ParseHeaders(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	var headers []string
	start := 0
	for i := 1; i < len(parts); i++ {
		key, _, ok := strings.Cut(parts[i], "=")
		if ok && httpguts.ValidHeaderFieldName(strings.TrimSpace(key)) {
			headers = append(headers, strings.Join(parts[start:i], ","))
			start = i
			parts[i] = strings.TrimLeft(parts[i], " \t")
		}
	}
	return append(headers, strings.Join(parts[start:], ","))
}

// toHTTPHeaders converts a slice of strings in the format "Key=Value" to http.Header.
func toHTTPHeaders(headers []string) http.Header {
	httpHeaders := make(http.Header)
	for _, header := range headers {
		key, value, ok := strings.Cut(header, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue // Skip empty keys or values
		}

		httpHeaders.Add(key, value)
	}

	return httpHeaders
}

type Config struct {
	ProxyURL   string
	SocksURL   string
	Timeout    time.Duration
	Headers    []string
	SkipVerify bool
}

func (cfg *Config) New() (*http.Client, error) {
	// Create a custom HTTP transport with TLS configuration.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.SkipVerify,
		},
	}

	// Configure proxy settings if provided.
	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %s", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	} else if cfg.SocksURL != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.SocksURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("can't connect to the SOCKS5 proxy: %s", err)
		}
		transport.DialContext = dialer.(proxy.ContextDialer).DialContext
	}

	// Create the HTTP client with the custom transport and headers.
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &customTransport{
			transport: transport,
			headers:   toHTTPHeaders(cfg.Headers),
		},
	}

	return httpClient, nil
}
