package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

// Module exported for returning the application Config.
var Module = fx.Options(
	fx.Provide(New),
)

// Config is the configuration for the application.
type Config struct {
	ServiceName         string `envconfig:"SERVICE_NAME" default:"spoty"`
	Prod                bool   `envconfig:"PROD" default:"false"`
	SpotifyClientID     string `envconfig:"SPOTIFY_CLIENT_ID" required:"true"`
	SpotifyClientSecret string `envconfig:"SPOTIFY_CLIENT_SECRET" required:"true"`
	HttpServerHost      string `envconfig:"HTTP_SERVER_HOST" default:"localhost"`
	HttpServerPort      int    `envconfig:"HTTP_SERVER_PORT" default:"13337"`
	CacheMaxKeys        int64  `envconfig:"CACHE_MAX_KEYS" default:"64"`
	CacheMaxCost        int64  `envconfig:"CACHE_MAX_COST" default:"1000000"`
	JaegerEndpoint      string `envconfig:"JAEGER_ENDPOINT" default:"http://localhost:14268/api/traces"`
	AMQPURI             string `envconfig:"AMQP_URI" default:"amqp://guest:guest@localhost:5672"`
}

// New processes and returns a new application Config.
func New() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}

	return &cfg, nil
}
