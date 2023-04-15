package http

import (
	"github.com/PaulBarrie/infra-worker/pkg/common"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/plugin"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithPluginController() *Server {
	server.Router.GET("/plugins/:providerType/:resourceType", func(c *gin.Context) {
		resp, err := server.PluginService.GetPlugin(dto.GetPluginRequest{
			ProviderType: common.ProviderType(c.Param("providerType")),
			ResourceType: common.ResourceType(c.Param("resourceType")),
		})
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"payload": resp.Plugin,
		})
	})
	return server
}
