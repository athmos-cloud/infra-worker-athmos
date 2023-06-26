package http

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithResourceController() *Server {
	server.Engine.GET("/resources/providers/:projectId", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("projectId"))
		c.Set(context.ResourceTypeKey, types.ProviderResource)
		server.ResourceController.ListResources(c)
	})
	server.Engine.GET("/resources/stack/:projectId/:providerType/:providerId", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("projectId"))
		c.Set(context.ResourceTypeKey, types.ProviderResource)
		providerType, err := types.ProviderFromString(c.Query("providerType"))
		if !err.IsOk() {
			c.JSON(400, gin.H{"error": err.ToString()})
		}
		c.Set(context.ProviderTypeKey, providerType)
		c.Set(context.RequestKey, dto.GetProviderStackRequest{ProviderID: c.Param("providerId")})
		server.ResourceController.GetResourceStack(c)
	})
	server.Engine.GET("/resources/:projectId/:providerType/:resourceType/:resourceId", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("projectId"))
		providerType, err := types.ProviderFromString(c.Query("providerType"))
		if !err.IsOk() {
			c.JSON(400, gin.H{"error": err.ToString()})
		}
		resourceType, err := types.ResourceFromString(c.Query("resourceType"))
		if !err.IsOk() {
			c.JSON(400, gin.H{"error": err.ToString()})
		}
		c.Set(context.ProviderTypeKey, providerType)
		c.Set(context.ResourceTypeKey, resourceType)
		c.Set(context.RequestKey, dto.GetResourceRequest{Identifier: c.Param("resourceId")})
		server.ResourceController.GetResource(c)
	})
	return server
}
