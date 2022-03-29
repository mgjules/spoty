package cmd

import (
	"context"

	"github.com/JulesMike/spoty/build"
	"github.com/JulesMike/spoty/cache"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/health"
	"github.com/JulesMike/spoty/http"
	"github.com/JulesMike/spoty/logger"
	"github.com/JulesMike/spoty/spoty"
	"github.com/JulesMike/spoty/tracer"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			build.Module,
			config.Module,
			logger.Module,
			tracer.Module,
			cache.Module,
			health.Module,
			http.Module,
			spoty.Module,
			fx.Invoke(serve),
		).Run()
	},
}

func serve(lc fx.Lifecycle, s *http.Server) error {
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

func init() {
	rootCmd.AddCommand(serveCmd)
}
