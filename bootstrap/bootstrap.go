package bootstrap

import (
	"context"

	"github.com/JulesMike/spoty/cache"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/server"
	"github.com/JulesMike/spoty/spoty"
	"go.uber.org/fx"
)

// Module exported for initialising application.
var Module = fx.Options(
	config.Module,
	cache.Module,
	server.Module,
	spoty.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lc fx.Lifecycle, s *server.Server) error {
	s.RegisterRoutes()

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})

	return nil
}
