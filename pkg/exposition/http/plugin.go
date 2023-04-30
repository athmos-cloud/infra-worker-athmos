package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithPluginController() *Server {
	server.Router.GET("/plugins/:providerType/:resourceType", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		resp := server.PluginService.GetPlugin(dto.GetPluginRequest{
			ProviderType: common.ProviderType(c.Param("providerType")),
			ResourceType: common.ResourceType(c.Param("resourceType")),
		})

		c.JSON(err.Code, gin.H{
			"payload": resp.Plugin,
		})
	})
	return server
}
