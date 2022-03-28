package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JulesMike/spoty/build"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/spoty"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

const (
	_readTimeout       = 2 * time.Second
	_writeTimeout      = 2 * time.Second
	_idleTimeout       = 30 * time.Second
	_readHeaderTimeout = 2 * time.Second
)

// Module exported for initialising a new Server.
var Module = fx.Options(
	fx.Provide(New),
)

// Server is the main HTTP server.
type Server struct {
	router *gin.Engine
	http   *http.Server
	spoty  *spoty.Spoty
	build  *build.Info
	addr   string
}

// New creates a new Server.
func New(cfg *config.Config, spoty *spoty.Spoty, build *build.Info) *Server {
	if cfg.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	s := Server{
		router: gin.Default(),
		addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		spoty:  spoty,
		build:  build,
	}

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
	log.Println("Listening on http://" + s.addr + " ...")
	if err := s.http.ListenAndServe(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}

// Stop stops the server.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server ...")
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
