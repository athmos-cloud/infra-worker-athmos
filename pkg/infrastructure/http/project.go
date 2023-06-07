package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/gin-gonic/gin"
)

func (server *Server) WithProjectRouter() *Server {
	server.Engine.GET("/projects/:id", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("id"))
		server.ProjectController.GetProject(c)
	})

	server.Engine.GET("/projects/owner/:id", func(c *gin.Context) {
		c.Set(context.OwnerIDKey, c.Param("id"))
		server.ProjectController.ListProjectByOwner(c)
	})

	server.Engine.POST("/projects", func(c *gin.Context) {
		var req dto.CreateProjectRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": errors.BadRequest.WithMessage(fmt.Sprintf("Expected : %+v", dto.CreateProjectRequest{}))})
			return
		}
		c.Set(context.RequestKey, req)
		server.ProjectController.CreateProject(c)
	})

	server.Engine.PUT("/projects/:id", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("id"))
		var req dto.UpdateProjectRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": errors.BadRequest.WithMessage(fmt.Sprintf("Expected : %+v", dto.UpdateProjectRequest{}))})
			return
		}
		c.Set(context.RequestKey, req)
		server.ProjectController.UpdateProject(c)
	})

	server.Engine.DELETE("/projects/:id", func(c *gin.Context) {
		c.Set(context.ProjectIDKey, c.Param("id"))
		server.ProjectController.DeleteProject(c)
	})

	return server
}
