package http

import (
	"errors"
	"net/http"

	"github.com/JulesMike/spoty/docs"
	ahealth "github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	errRetrieveCurrentTrack = errors.New("failed to retrieve current playing track")
	errProcessCurrentTrack  = errors.New("failed to process images for currently playing track")
)

// Success defines the structure for a successful response.
type Success struct {
	Success string `json:"success"`
}

// Error defines the structure for a failed response.
type Error struct {
	Error string `json:"error"`
}

// handleHealthCheck godoc
// @Summary Health Check
// @Description checks if server is running
// @Tags core
// @Produce json
// @Success 200 {object} ahealth.CheckerResult
// @Success 503 {object} ahealth.CheckerResult
// @Router / [get]
func (s *Server) handleHealthCheck() gin.HandlerFunc {
	opts := s.health.CompileHealthCheckerOption()
	checker := ahealth.NewChecker(opts...)

	return gin.WrapF(
		ahealth.NewHandler(
			checker,
		),
	)
}

// handleVersion godoc
// @Summary Health Check
// @Description checks the server's version
// @Tags core
// @Produce json
// @Success 200 {object} build.Info
// @Router /api/version [get]
func (s *Server) handleVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.build)
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
// @Failure 401 {object} http.Error "not authenticated"
// @Failure 404 {object} http.Error "no current playing track found"
// @Router /api/current [get]
func (s *Server) handleCurrentTrack(c *gin.Context) {
	ctx := c.Request.Context()

	track, err := s.spoty.TrackCurrentlyPlaying(ctx)
	if err != nil {
		s.logger.ErrorwContext(ctx, errRetrieveCurrentTrack.Error(), "error", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, Error{Error: errRetrieveCurrentTrack.Error()})

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
// @Failure 401 {object} http.Error "not authenticated"
// @Failure 404 {object} http.Error "no current playing track found"
// @Failure 500 {object} http.Error "album images could not be processed"
// @Router /api/current/images [get]
func (s *Server) handleCurrentTrackImages(c *gin.Context) {
	ctx := c.Request.Context()

	track, err := s.spoty.TrackCurrentlyPlaying(ctx)
	if err != nil {
		s.logger.ErrorwContext(ctx, errRetrieveCurrentTrack.Error(), "error", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, Error{Error: errRetrieveCurrentTrack.Error()})

		return
	}

	images, err := s.spoty.TrackImages(ctx, track)
	if err != nil {
		s.logger.ErrorwContext(ctx, errProcessCurrentTrack.Error(), "error", err.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			Error{Error: errProcessCurrentTrack.Error()},
		)

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
// @Failure 403 {object} http.Error "already authenticated"
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
// @Success 200 {object} http.Success "authenticated successfully"
// @Failure 403 {object} http.Error "already authenticated"
// @Failure 403 {object} http.Error "could not retrieve token"
// @Failure 404 {object} http.Error "could not retrieve current user"
// @Router /api/callback [get]
func (s *Server) handleCallback(c *gin.Context) {
	if err := s.spoty.SetupNewClient(c.Request); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, Error{Error: "could not retrieve token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "welcome, you are now authenticated!"})
}
