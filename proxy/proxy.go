package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

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

	// Create a custom HTTP client with the specified timeout and transport.
	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
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
	return httpClient, nil
}
