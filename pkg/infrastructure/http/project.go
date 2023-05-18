package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/share"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithProjectRouter() *Server {
	server.Engine.GET("/projects/:id", func(c *gin.Context) {
		c.Set(share.ProjectIDKey, c.Param("id"))
		server.ProjectController.GetProject(c)
		if err := c.Value(share.ErrorContextKey); err != nil {
			handleError(c, err)
		} else {
			c.JSON(200, gin.H{
				"payload": c.Value(share.ResponseContextKey).(dto.GetProjectResponse),
			})
		}
	})

	server.Engine.GET("/projects/owner/:id", func(c *gin.Context) {
		c.Set(share.OwnerIDKey, c.Param("id"))
		server.ProjectController.ListProjectByOwner(c)
		if err := c.Value(share.ErrorContextKey); err != nil {
			handleError(c, err)
		} else {
			c.JSON(200, gin.H{
				"payload": c.Value(share.ResponseContextKey).(dto.ListProjectResponse),
			})
		}
	})

	server.Engine.POST("/projects", func(c *gin.Context) {
		server.ProjectController.CreateProject(c)
	})

	server.Engine.PUT("/projects/:id", func(c *gin.Context) {
		type Request struct {
			Name string `json:"name"`
		}
		var request Request
		if err := c.BindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprintf("Wrong request body: %s", err),
			})
		}
		c.Set(share.RequestContextKey, dto.UpdateProjectRequest{
			ID:   c.Param("id"),
			Name: request.Name,
		})
		server.ProjectController.ListProjectByOwner(c)
		if err := c.Value(share.ErrorContextKey); err != nil {
			handleError(c, err)
		} else {
			c.JSON(204, gin.H{
				"message": fmt.Sprintf("UpdatedProject %s updated", c.Param("id")),
			})
		}
	})

	server.Engine.DELETE("/projects/:id", func(c *gin.Context) {
		c.Set(share.ProjectIDKey, c.Param("id"))
		server.ProjectController.DeleteProject(c)
		if err := c.Value(share.ErrorContextKey); err != nil {
			handleError(c, err)
		} else {
			c.JSON(204, gin.H{
				"message": fmt.Sprintf("Project %s deleted", c.Param("id")),
			})
		}
	})

	return server
}
