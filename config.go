package ucodesdk

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	appId          string
	BaseURL        string
	FunctionName   string
	RequestTimeout time.Duration
}

func (cfg *Config) SetAppId() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	// Get the APP_ID value from APP_ID variables
	appId := os.Getenv("APP_ID")
	if appId == "" {
		return fmt.Errorf("APP_ID environment variable not set")
	}

	cfg.appId = appId
	// fmt.Println(cfg.appId)
	return nil
}

func (cfg *Config) SetBaseUrl(url string) {
	cfg.BaseURL = url
}
