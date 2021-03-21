package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	http   *http.Server
	addr   string
}

func New(prod bool, host string, port int) *Server {
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}

	s := Server{
		router: gin.Default(),
		addr:   fmt.Sprintf("%s:%d", host, port),
	}

	s.http = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	s.registerHealthCheck()
	s.registerSwagger()

	return &s
}

func (s *Server) APIRoute() *gin.RouterGroup {
	return s.router.Group("/api")
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
