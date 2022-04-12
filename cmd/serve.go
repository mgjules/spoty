package cmd

import (
	"context"

	"github.com/mgjules/spoty/build"
	"github.com/mgjules/spoty/cache"
	"github.com/mgjules/spoty/config"
	"github.com/mgjules/spoty/health"
	"github.com/mgjules/spoty/logger"
	"github.com/mgjules/spoty/spoty"
	"github.com/mgjules/spoty/tracer"
	"github.com/mgjules/spoty/transport/http"
	"github.com/mgjules/spoty/transport/messenger"
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
			messenger.Module,
			http.Module,
			spoty.Module,
			fx.Invoke(serve),
		).Run()
	},
}

func serve(lc fx.Lifecycle, s *http.Server) error {
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
