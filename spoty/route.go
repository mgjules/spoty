package spoty

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Spoty) RegisterRoutes(router *gin.RouterGroup) {
	// Guest routes
	guest := router.Group("/")
	guest.Use(s.unauthenticatedOnlyMiddleware())
	{
		guest.GET("/authenticate", s.authenticateHandler)
		guest.GET("/callback", s.callbackHandler)
	}

	// Authenticated routes
	authenticated := router.Group("/")
	authenticated.Use(s.authenticatedOnlyMiddleware())
	{
		authenticated.GET("/current", s.currentHandler)
		authenticated.GET("/current/images", s.currentImagesHandler)
	}
}

func (s *Spoty) unauthenticatedOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.client != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you are already authenticated"})
			return
		}

		c.Next()
	}
}

func (s *Spoty) authenticatedOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.client == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be authenticated to access this endpoint"})
			return
		}

		c.Next()
	}
}

// currentHandler godoc
// @Summary Current Playing Track
// @Description returns information about the current playing track
// @Tags spoty
// @Produce json
// @Success 200 {object} spotify.FullTrack "returns full track information"
// @Failure 401 {object} server.Error "not authenticated"
// @Failure 404 {object} server.Error "no current playing track found"
// @Router /api/current [get]
func (s *Spoty) currentHandler(c *gin.Context) {
	track, err := s.trackCurrentlyPlaying()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "could not retrieve currently playing track"})
		return
	}

	c.JSON(http.StatusOK, track)
}

// currentImagesHandler godoc
// @Summary Album Images of Current Playing Track
// @Description returns the album images of the current playing track
// @Tags spoty
// @Produce json
// @Success 200 {array} Image "returns album images"
// @Failure 401 {object} server.Error "not authenticated"
// @Failure 404 {object} server.Error "no current playing track found"
// @Failure 500 {object} server.Error "album images could not be processed"
// @Router /api/current/images [get]
func (s *Spoty) currentImagesHandler(c *gin.Context) {
	track, err := s.trackCurrentlyPlaying()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "could not retrieve currently playing track"})
		return
	}

	images, err := s.trackImages(track)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not process images for currently playing track"})
		return
	}

	c.JSON(http.StatusOK, images)
}

// authenticateHandler godoc
// @Summary Authentication
// @Description redirects user to spotify for authentication
// @Tags spoty
// @Produce json
// @Success 302 {string} string "redirection to spotify"
// @Failure 403 {object} server.Error "already authenticated"
// @Router /api/authenticate [get]
func (s *Spoty) authenticateHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, s.auth.AuthURL(s.state))
}

// callbackHandler godoc
// @Summary Callback
// @Description spotify redirects to the this endpoint on success
// @Tags spoty
// @Produce json
// @Param code query string true "code from spotify"
// @Param state query string true "state from spotify"
// @Success 200 {object} server.Success "authenticated successfully"
// @Failure 403 {object} server.Error "already authenticated"
// @Failure 403 {object} server.Error "could not retrieve token"
// @Failure 404 {object} server.Error "could not retrieve current user"
// @Router /api/callback [get]
func (s *Spoty) callbackHandler(c *gin.Context) {
	tok, err := s.auth.Token(s.state, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "could not retrieve token"})
		return
	}

	client := s.auth.NewClient(tok)
	if _, err := client.CurrentUser(); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "could not retrieve current user"})
		return
	}

	client.AutoRetry = true

	s.client = &client

	c.JSON(http.StatusOK, gin.H{"success": "welcome, you are now authenticated!"})
}
