package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithInternalController() *Server {
	server.Router.GET("/ready", func(c *gin.Context) {
		err := errors.OK
		c.JSON(err.Code, gin.H{
			"payload": "ready",
		})
	})
	return server
}
