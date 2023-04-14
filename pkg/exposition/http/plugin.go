package http

import "github.com/gin-gonic/gin"

func (server *Server) WithPluginController() *Server {
	server.Router.GET("/plugins/:providerType/:resourceType", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return server
}
