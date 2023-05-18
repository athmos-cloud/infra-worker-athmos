package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	ProjectController controller.Project
}

func New(
	projectService controller.Project,
) *Server {
	return &Server{
		Engine:            gin.Default(),
		ProjectController: projectService,
	}
}

func (server *Server) Start() {
	server.WithProjectRouter().WithInternalController()
	err := server.Engine.Run(fmt.Sprintf(":%d", config.Current.Http.Port))
	if err != nil {
		panic(err)
	}
}
