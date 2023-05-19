package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithSecretRouter() *Server {
	server.Engine.GET("/secrets/:projectId/:name", func(c *gin.Context) {
		c.Set(share.RequestContextKey, dto.GetSecretRequest{
			ProjectID: c.Param("projectId"),
			Name:      c.Param("name"),
		})
		server.SecretController.GetSecret(c)
	})

	server.Engine.GET("/secrets/:projectId", func(c *gin.Context) {
		c.Set(share.ProjectIDKey, c.Param("projectId"))
		server.SecretController.ListProjectSecret(c)
	})

	server.Engine.POST("/secrets", func(c *gin.Context) {
		var request dto.CreateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{"message": errors.BadRequest.WithMessage(fmt.Sprintf("Expected : %+v", dto.CreateSecretRequest{}))})
			return
		}
		c.Set(share.ProjectIDKey, request.ProjectID)
		c.Set(share.RequestContextKey, request)
		server.SecretController.CreateSecret(c)
	})

	server.Engine.PUT("/secrets/:projectId/:name", func(c *gin.Context) {
		var request dto.UpdateSecretRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{"message": errors.BadRequest.WithMessage(fmt.Sprintf("Expected : %+v", dto.UpdateSecretRequest{}))})
			return
		}
		c.Set(share.ProjectIDKey, c.Param("projectId"))
		c.Set(share.RequestContextKey, request)
		server.SecretController.UpdateSecret(c)
	})

	server.Engine.DELETE("/secrets/:projectId/:name", func(c *gin.Context) {
		c.Set(share.ProjectIDKey, c.Param("projectId"))
		c.Set(share.RequestContextKey, dto.DeleteSecretRequest{
			Name: c.Param("name"),
		})
		server.SecretController.DeleteSecret(c)
	})

	return server
}
