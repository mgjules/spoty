package server

import (
	"net/http"

	"github.com/JulesMike/spoty/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Success struct {
	Success string `json:"success"`
}

type Error struct {
	Error string `json:"error"`
}

// handleHealthCheck godoc
// @Summary Health Check
// @Description checks if server is running
// @Tags core
// @Produce json
// @Success 200 {object} server.Success
// @Router /api [get]
func (s *Server) handleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Success{Success: "i'm alright!"})
	}
}

func (s *Server) handleSwagger() gin.HandlerFunc {
	docs.SwaggerInfo.Host = s.addr
	docs.SwaggerInfo.BasePath = "/"

	url := ginSwagger.URL("http://" + s.addr + "/swagger/doc.json")
	return ginSwagger.WrapHandler(swaggerFiles.Handler, url)
}

// handleCurrentTrack godoc
// @Summary Current Playing Track
// @Description returns information about the current playing track
// @Tags spoty
// @Produce json
// @Success 200 {object} spotify.FullTrack "returns full track information"
// @Failure 401 {object} server.Error "not authenticated"
// @Failure 404 {object} server.Error "no current playing track found"
// @Router /api/current [get]
func (s *Server) handleCurrentTrack(c *gin.Context) {
	track, err := s.spoty.TrackCurrentlyPlaying()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "could not retrieve currently playing track"})
		return
	}

	c.JSON(http.StatusOK, track)
}

// handleCurrentTrackImages godoc
// @Summary Album Images of Current Playing Track
// @Description returns the album images of the current playing track
// @Tags spoty
// @Produce json
// @Success 200 {array} spoty.Image "returns album images"
// @Failure 401 {object} server.Error "not authenticated"
// @Failure 404 {object} server.Error "no current playing track found"
// @Failure 500 {object} server.Error "album images could not be processed"
// @Router /api/current/images [get]
func (s *Server) handleCurrentTrackImages(c *gin.Context) {
	track, err := s.spoty.TrackCurrentlyPlaying()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "could not retrieve currently playing track"})
		return
	}

	images, err := s.spoty.TrackImages(track)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not process images for currently playing track"})
		return
	}

	c.JSON(http.StatusOK, images)
}

// handleAuthenticate godoc
// @Summary Authentication
// @Description redirects user to spotify for authentication
// @Tags spoty
// @Produce json
// @Success 302 {string} string "redirection to spotify"
// @Failure 403 {object} server.Error "already authenticated"
// @Router /api/authenticate [get]
func (s *Server) handleAuthenticate(c *gin.Context) {
	c.Redirect(http.StatusFound, s.spoty.AuthURL())
}

// handleCallback godoc
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
func (s *Server) handleCallback(c *gin.Context) {
	if err := s.spoty.SetupNewClient(c.Request); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "could not retrieve token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "welcome, you are now authenticated!"})
}
