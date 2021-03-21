package app

import (
	"context"

	"github.com/JulesMike/spoty/cache"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/server"
	"github.com/JulesMike/spoty/spoty"
	"go.uber.org/fx"
)

func ProvideConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ProvideCache(cfg *config.Config) (*cache.Cache, error) {
	cache, err := cache.New(cfg.CacheMaxKeys, cfg.CacheMaxCost)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func ProvideSpoty(cfg *config.Config, cache *cache.Cache) (*spoty.Spoty, error) {
	spoty, err := spoty.New(cfg.ClientID, cfg.ClientSecret, cfg.Host, cfg.Port, cache)
	if err != nil {
		return nil, err
	}

	return spoty, nil
}

func ProvideServer(lc fx.Lifecycle, cfg *config.Config) (*server.Server, error) {
	server := server.New(cfg.Prod, cfg.Host, cfg.Port)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go server.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx)
		},
	})

	return server, nil
}
