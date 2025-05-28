package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

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

// toHTTPHeaders converts a slice of strings in the format "Key=Value" to http.Header.
func toHTTPHeaders(headers []string) http.Header {
	httpHeaders := make(http.Header)
	for _, header := range headers {
		parts := strings.SplitN(header, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "" || value == "" {
				continue // Skip empty keys or values
			}
			httpHeaders.Add(key, value)
		}
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
