package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithSecretRouter() *Server {
	server.Router.GET("/secrets/:projectId/:name", func(c *gin.Context) {
		err := errors.OK
		projectID := c.Param("projectId")
		name := c.Param("name")
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
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
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		//resp := server.ProjectService.GetProjectByOwnerID(c, dtoProject.GetProjectByOwnerIDRequest{
		//	OwnerID: c.Param("id"),
		//})
		//
		//jsonBytes, errMarshal := json.Marshal(resp.Payload)
		//if errMarshal != nil {
		//	c.JSON(500, gin.H{
		//		"message": fmt.Sprintf("Error marshalling response: %s", errMarshal),
		//	})
		//	return
		//}
		//c.JSON(err.Code, gin.H{
		//	"payload": string(jsonBytes[:]),
		//})
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
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		server.SecretService.CreateSecret(c, request)

		c.JSON(err.Code, gin.H{
			"message": "Secret created",
		})
	})
	server.Router.PUT("/secrets/:projectId/:name", func(c *gin.Context) {
		//err := errors.NoContent
		//type Request struct {
		//	Name string `json:"name"`
		//}
		//var request Request
		//errRequestBody := c.BindJSON(&request)
		//if errRequestBody != nil {
		//	c.JSON(400, gin.H{
		//		"message": fmt.Sprintf("Wrong request body: %s", errRequestBody),
		//	})
		//	return
		//}
		//defer func() {
		//	if r := recover(); r != nil {
		//		err = r.(errors.Error)
		//		c.JSON(err.Code, gin.H{
		//			"message": err.ToString(),
		//		})
		//	}
		//}()
		//server.ProjectService.UpdateProjectName(c, dtoProject.UpdateProjectRequest{
		//	ProjectID:   c.Param("id"),
		//	ProjectName: request.Name,
		//})
		//
		//c.JSON(err.Code, gin.H{
		//	"message": fmt.Sprintf("UpdatedProject %s updated", c.Param("id")),
		//})
	})
	server.Router.DELETE("/secrets/:projectId/:name", func(c *gin.Context) {
		//err := errors.NoContent
		//defer func() {
		//	if r := recover(); r != nil {
		//		err = r.(errors.Error)
		//		c.JSON(err.Code, gin.H{
		//			"message": err.ToString(),
		//		})
		//	}
		//}()
		//server.ProjectService.DeleteProject(c, dtoProject.DeleteRequest{
		//	ProjectID: c.Param("id"),
		//})
		//c.JSON(err.Code, gin.H{
		//	"message": fmt.Sprintf("UpdatedProject %s deleted", c.Param("id")),
		//})
	})
	return server
}
