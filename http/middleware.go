package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) unauthenticatedOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.spoty.IsAuth() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you are already authenticated"})

			return
		}

		c.Next()
	}
}

func (s *Server) authenticatedOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.spoty.IsAuth() {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "you must be authenticated to access this endpoint"},
			)

			return
		}

		c.Next()
	}
}
