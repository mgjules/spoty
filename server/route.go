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

// registerHealthCheck godoc
// @Summary Health Check
// @Description checks if server is running
// @Tags core
// @Produce json
// @Success 200 {object} Success
// @Router /api [get]
func (s *Server) registerHealthCheck() {
	s.APIRoute().GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, Success{Success: "i'm alright!"})
	})
}

func (s *Server) registerSwagger() {
	docs.SwaggerInfo.Host = s.addr
	docs.SwaggerInfo.BasePath = "/"

	url := ginSwagger.URL("http://" + s.addr + "/swagger/doc.json")
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
