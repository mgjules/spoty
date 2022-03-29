package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JulesMike/spoty/build"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/health"
	"github.com/JulesMike/spoty/logger"
	"github.com/JulesMike/spoty/spoty"
	"github.com/JulesMike/spoty/tracer"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
)

const (
	_readTimeout       = 2 * time.Second
	_writeTimeout      = 2 * time.Second
	_idleTimeout       = 30 * time.Second
	_readHeaderTimeout = 2 * time.Second
)

// Module exported for initialising a new Server and Client.
var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Provide(NewClient),
)

// Server is the main HTTP server.
type Server struct {
	router *gin.Engine
	http   *http.Server
	logger *logger.Logger
	tracer *tracer.Tracer
	spoty  *spoty.Spoty
	health *health.Checks
	build  *build.Info
	addr   string
}

// NewServer creates a new Server.
func NewServer(
	cfg *config.Config,
	logger *logger.Logger,
	tracer *tracer.Tracer,
	spoty *spoty.Spoty,
	health *health.Checks,
	build *build.Info,
) *Server {
	if cfg.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	s := Server{
		router: gin.Default(),
		addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		logger: logger,
		tracer: tracer,
		spoty:  spoty,
		health: health,
		build:  build,
	}

	desugared := logger.Desugar()

	s.router.Use(ginzap.Ginzap(desugared.Logger, time.RFC3339, true))
	s.router.Use(ginzap.RecoveryWithZap(desugared.Logger, true))
	s.router.Use(otelgin.Middleware("main"))

	s.http = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       _readTimeout,
		WriteTimeout:      _writeTimeout,
		IdleTimeout:       _idleTimeout,
		ReadHeaderTimeout: _readHeaderTimeout,
	}

	return &s
}

// RegisterRoutes registers the REST HTTP routes.
func (s *Server) RegisterRoutes() {
	// Health Check
	s.router.GET("/", s.handleHealthCheck())

	// Swagger
	s.router.GET("/swagger/*any", s.handleSwagger())

	api := s.router.Group("/api")
	{
		// Version
		api.GET("/version", s.handleVersion())

		// Guest routes
		guest := api.Group("/")
		guest.Use(s.unauthenticatedOnly())
		{
			guest.GET("/authenticate", s.handleAuthenticate)
			guest.GET("/callback", s.handleCallback)
		}

		// Authenticated routes
		authenticated := api.Group("/")
		authenticated.Use(s.authenticatedOnly())
		{
			authenticated.GET("/current", s.handleCurrentTrack)
			authenticated.GET("/current/images", s.handleCurrentTrackImages)
		}
	}
}

// Start starts the server.
// It blocks until the server stops.
func (s *Server) Start() error {
	s.logger.Infof("Listening on http://%s...", s.addr)

	if err := s.http.ListenAndServe(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}

// Stop stops the server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server ...")

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
