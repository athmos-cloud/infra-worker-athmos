package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

type UpdateSecretHttpRequest struct {
	Data        string `json:"data"`
	Description string `json:"description"`
}

func (server *Server) WithSecretRouter() *Server {
	server.Router.GET("/secrets/:projectId/:name", func(c *gin.Context) {
		err := errors.OK
		projectID := c.Param("projectId")
		name := c.Param("name")
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()
		resp := server.SecretService.GetSecret(c, secret.GetSecretRequest{
			Name:      name,
			ProjectID: projectID,
		})

		c.JSON(err.Code, gin.H{
			"payload": resp,
		})
	})
	server.Router.GET("/secrets/:projectId", func(c *gin.Context) {
		err := errors.OK
		projectID := c.Param("projectId")
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()

		resp := server.SecretService.ListSecret(c, secret.ListSecretRequest{ProjectID: projectID})
		c.JSON(err.Code, gin.H{
			"payload": resp,
		})
	})
	server.Router.POST("/secrets", func(c *gin.Context) {
		var request secret.CreateSecretRequest
		err := errors.Created
		errRequestBody := c.BindJSON(&request)
		if errRequestBody != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("Wrong request body: %s", errRequestBody),
			})
		}
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()
		server.SecretService.CreateSecret(c, request)

		c.JSON(err.Code, gin.H{
			"message": "Secret created",
		})
	})

	server.Router.PUT("/secrets/:projectId/:name", func(c *gin.Context) {
		var request UpdateSecretHttpRequest
		err := errors.Created
		errRequestBody := c.BindJSON(&request)
		if errRequestBody != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("Wrong request body: %s", errRequestBody),
			})
		}
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()

		server.SecretService.UpdateSecret(c, secret.UpdateSecretRequest{
			ProjectID:   c.Param("projectId"),
			Name:        c.Param("name"),
			Description: request.Description,
			Data:        request.Data,
		})

		c.Status(err.Code)
	})
	server.Router.DELETE("/secrets/:projectId/:name", func(c *gin.Context) {
		err := errors.NoContent
		defer func() {
			if r := recover(); r != nil {
				handleError(c, r)
			}
		}()

		server.SecretService.DeleteSecret(c, secret.DeleteSecretRequest{
			ProjectID: c.Param("projectId"),
			Name:      c.Param("name"),
		})
		c.Status(err.Code)
	})
	return server
}
