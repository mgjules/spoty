package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) unauthenticatedOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.spoty.IsAuth() {
			rErr := NewError(
				"already-authenticated",
				"You are already authenticated.",
				http.StatusForbidden,
				"You cannot authenticate again as you are already authenticated.",
				c.Request.URL.String(),
				nil,
			)

			ctx := c.Request.Context()
			s.logger.ErrorwContext(ctx, "failed to authenticate", "error", rErr.Error())
			c.AbortWithStatusJSON(http.StatusForbidden, rErr)

			return
		}

		c.Next()
	}
}

func (s *Server) authenticatedOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.spoty.IsAuth() {
			rErr := NewError(
				"not-authenticated",
				"You do not have access.",
				http.StatusUnauthorized,
				"You cannot access this endpoint because you are not authenticated.",
				c.Request.URL.String(),
				nil,
			)

			ctx := c.Request.Context()
			s.logger.ErrorwContext(ctx, "failed to access endpoint", "error", rErr.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, rErr)

			return
		}

		c.Next()
	}
}
