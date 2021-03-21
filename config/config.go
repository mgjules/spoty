package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Prod         bool   `envconfig:"PROD" default:"false"`
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
	Host         string `envconfig:"HOST" default:"localhost"`
	Port         int    `envconfig:"PORT" default:"13337"`
	CacheMaxKeys int64  `envconfig:"CACHE_MAX_KEYS" default:"64"`
	CacheMaxCost int64  `envconfig:"CACHE_MAX_COST" default:"1000000"`
}

func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process env vars: %w", err)
	}

	return &cfg, nil
}
