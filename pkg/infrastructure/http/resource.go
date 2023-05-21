package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithResourceController() *Server {
	server.Engine.GET("/resources/:projectId", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("projectId"))
		c.Set(context.ResourceTypeKey, types.ProviderResource)
		server.ResourceController.ListResources(c)
	})
	server.Engine.GET("/resources/:projectId/:providerType/:providerId", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("projectId"))
		c.Set(context.ResourceTypeKey, types.ProviderResource)
		c.Set(context.ProviderTypeKey, c.Param("providerType"))
		c.Set(context.RequestKey, dto.GetProviderStackRequest{ProviderID: c.Param("providerId")})
		server.ResourceController.ListResources(c)
	})

	return server
}
