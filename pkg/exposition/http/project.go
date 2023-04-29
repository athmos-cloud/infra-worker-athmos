package http

import (
	"encoding/json"
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"

	dtoProject "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/project"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithProjectRouter() *Server {
	server.Router.GET("/projects/:id", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
			}
		}()
		resp := server.ProjectService.GetProjectByID(c, dtoProject.GetProjectByIDRequest{
			ProjectID: c.Param("id"),
		})
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"payload": resp.Payload,
		})
	})
	server.Router.GET("/projects/owner/:id", func(c *gin.Context) {
		err := errors.OK
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
			}
		}()
		resp := server.ProjectService.GetProjectByOwnerID(c, dtoProject.GetProjectByOwnerIDRequest{
			OwnerID: c.Param("id"),
		})
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
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
		var request dtoProject.CreateProjectRequest
		err := errors.Created
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
			}
		}()
		resp := server.ProjectService.CreateProject(c, request)
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"project_id": resp.ProjectID,
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
			}
		}()
		server.ProjectService.UpdateProjectName(c, dtoProject.UpdateProjectRequest{
			ProjectID:   c.Param("id"),
			ProjectName: request.Name,
		})
		if err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"message": fmt.Sprintf("Project %s updated", c.Param("id")),
		})
	})
	server.Router.DELETE("/projects/:id", func(c *gin.Context) {
		err := errors.NoContent
		defer func() {
			if r := recover(); r != nil {
				err = r.(errors.Error)
			}
		}()
		server.ProjectService.DeleteProject(c, dtoProject.DeleteRequest{
			ProjectID: c.Param("id"),
		})
		if !err.IsOk() {
			c.JSON(err.Code, gin.H{
				"message": err.ToString(),
			})
			return
		}
		c.JSON(err.Code, gin.H{
			"message": fmt.Sprintf("Project %s deleted", c.Param("id")),
		})
	})
	return server
}
