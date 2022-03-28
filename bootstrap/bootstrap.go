package bootstrap

import (
	"context"

	"github.com/JulesMike/spoty/build"
	"github.com/JulesMike/spoty/cache"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/health"
	"github.com/JulesMike/spoty/http"
	"github.com/JulesMike/spoty/logger"
	"github.com/JulesMike/spoty/spoty"
	"go.uber.org/fx"
)

// Module exported for initialising application.
var Module = fx.Options(
	build.Module,
	config.Module,
	logger.Module,
	cache.Module,
	health.Module,
	http.Module,
	spoty.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lc fx.Lifecycle, s *http.Server) error {
	s.RegisterRoutes()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.Start() //nolint: errcheck

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})

	return nil
}
