package http

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/application/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithPluginController() *Server {
	server.Router.GET("/plugins/:providerType/:resourceType", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()
		resp := server.PluginService.GetPlugin(dto.GetPluginRequest{
			ProviderType: types.ProviderType(c.Param("providerType")),
			ResourceType: types.ResourceType(c.Param("resourceType")),
		})

		c.JSON(err.Code, gin.H{
			"payload": resp.Plugin,
		})
	})
	return server
}
