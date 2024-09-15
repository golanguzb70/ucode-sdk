package ucodesdk

import (
	"time"
)

type Config struct {
	AppId          string
	BaseURL        string
	FunctionName   string
	RequestTimeout time.Duration
}

func (cfg *Config) SetBaseUrl(url string) {
	cfg.BaseURL = url
}