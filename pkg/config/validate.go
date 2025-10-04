package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	errInvalidLanguage = errors.New("invalid language")
	errMissingProvider = errors.New("provider cannot be empty")
	errMissingAPIKey   = errors.New("api_key cannot be empty")
)

// Validate runs domain-specific validation across all configuration scopes.
func (c *Config) Validate() error {
	if err := c.Git.Validate(); err != nil {
		return fmt.Errorf("git: %w", err)
	}
	if err := c.AI.Validate(); err != nil {
		return fmt.Errorf("ai: %w", err)
	}
	if err := c.Proxy.Validate(); err != nil {
		return fmt.Errorf("proxy: %w", err)
	}
	if err := c.Runtime.Validate(); err != nil {
		return fmt.Errorf("runtime: %w", err)
	}
	return nil
}

// Validate ensures git settings remain within allowed bounds.
func (g GitConfig) Validate() error {
	if g.DiffUnified < 0 {
		return fmt.Errorf("diff_unified must be >= 0")
	}
	if g.MaxInputSize <= 0 {
		return fmt.Errorf("max_input_size must be > 0")
	}
	if g.Lang != "" {
		if _, ok := supportedLangs[g.Lang]; !ok {
			return fmt.Errorf("%w: %s", errInvalidLanguage, g.Lang)
		}
	}
	return nil
}

// Validate checks AI configuration for required values and numeric ranges.
func (a AIConfig) Validate() error {
	if strings.TrimSpace(a.Provider) == "" {
		return errMissingProvider
	}
	if strings.TrimSpace(a.APIKey) == "" {
		return errMissingAPIKey
	}
	if a.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be > 0")
	}
	if a.Temperature < 0 || a.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1")
	}
	if a.TopP < 0 || a.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1")
	}
	return nil
}

// Validate ensures proxy URLs and timeout constraints are valid.
func (p ProxyConfig) Validate() error {
	if p.ProxyURL != "" {
		if _, err := url.ParseRequestURI(p.ProxyURL); err != nil {
			return fmt.Errorf("proxy_url invalid: %w", err)
		}
	}
	if p.SocksURL != "" {
		if _, err := url.ParseRequestURI(p.SocksURL); err != nil {
			return fmt.Errorf("socks_url invalid: %w", err)
		}
	}
	if p.Timeout < 0 {
		return fmt.Errorf("timeout must be >= 0")
	}
	return nil
}

// Validate runs validation for runtime sections.
func (r RuntimeConfig) Validate() error {
	if err := r.Review.Validate(); err != nil {
		return fmt.Errorf("review: %w", err)
	}
	if err := r.Commit.Validate(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// Validate ensures review runtime flags carry supported values.
func (r ReviewRuntime) Validate() error {
	if r.Mode != "" {
		switch r.Mode {
		case "local", "external":
		default:
			return fmt.Errorf("mode must be local or external")
		}
	}
	if r.MaxInput < 0 {
		return fmt.Errorf("max_input_size must be >= 0")
	}
	return nil
}

// Validate is a placeholder for future commit runtime constraints.
func (c CommitRuntime) Validate() error {
	return nil
}
