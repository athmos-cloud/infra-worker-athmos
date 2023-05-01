package http

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithProjectRouter() *Server {
	server.Router.GET("/projects/:id", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		resp := server.ProjectService.GetProjectByID(c, project.GetProjectByIDRequest{
			ProjectID: c.Param("id"),
		})

		c.JSON(err.Code, gin.H{
			"payload": resp.Payload,
		})
	})
	server.Router.GET("/projects/owner/:id", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		resp := server.ProjectService.GetProjectByOwnerID(c, project.GetProjectByOwnerIDRequest{
			OwnerID: c.Param("id"),
		})

		jsonBytes, errMarshal := json.Marshal(resp.Payload)
		if errMarshal != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprintf("Error marshalling response: %s", errMarshal),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"payload": string(jsonBytes[:]),
		})
	})
	server.Router.POST("/projects", func(c *gin.Context) {
		var request project.CreateProjectRequest
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
		resp := server.ProjectService.CreateProject(c, request)

		c.JSON(err.Code, gin.H{
			"projectID": resp.ProjectID,
		})
	})
	server.Router.PUT("/projects/:id", func(c *gin.Context) {
		err := errors.NoContent
		type Request struct {
			Name string `json:"name"`
		}
		var request Request
		errRequestBody := c.BindJSON(&request)
		if errRequestBody != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("Wrong request body: %s", errRequestBody),
			})
			return
		}
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		server.ProjectService.UpdateProjectName(c, project.UpdateProjectRequest{
			ProjectID:   c.Param("id"),
			ProjectName: request.Name,
		})

		c.JSON(err.Code, gin.H{
			"message": fmt.Sprintf("UpdatedProject %s updated", c.Param("id")),
		})
	})
	server.Router.DELETE("/projects/:id", func(c *gin.Context) {
		err := errors.NoContent
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
				c.JSON(err.Code, gin.H{
					"message": err.ToString(),
				})
			}
		}()
		server.ProjectService.DeleteProject(c, project.DeleteRequest{
			ProjectID: c.Param("id"),
		})
		c.JSON(err.Code, gin.H{
			"message": fmt.Sprintf("UpdatedProject %s deleted", c.Param("id")),
		})
	})
	return server
}
