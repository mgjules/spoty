package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/spoty"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Server struct {
	router *gin.Engine
	http   *http.Server
	spoty  *spoty.Spoty
	addr   string
}

func New(cfg *config.Config, spoty *spoty.Spoty) *Server {
	if cfg.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	s := Server{
		router: gin.Default(),
		addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		spoty:  spoty,
	}

	s.http = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return &s
}

func (s *Server) RegisterRoutes() {
	// Swagger
	s.router.GET("/swagger/*any", s.handleSwagger())

	api := s.router.Group("/api")
	{
		// Health Check
		api.GET("/", s.handleHealthCheck())

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

func (s *Server) Start() error {
	log.Println("Listening on http://" + s.addr + " ...")
	if err := s.http.ListenAndServe(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("stop: %w", err)
	}

	return nil
}
