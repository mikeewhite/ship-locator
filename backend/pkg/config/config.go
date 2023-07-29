package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WebSocketAPIKey string
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("SHIPLOC", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error on reading config env vars :%w", err)
	}
	return &cfg, nil
}
