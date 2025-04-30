package proxy

import (
	"net/http"
	"time"
)

type Config struct {
	ProxyURL   string
	SocksURL   string
	Timeout    time.Duration
	Headers    []string
	SkipVerify bool
}

func (cfg *Config) New() http.Client {
	return http.Client{}
}
