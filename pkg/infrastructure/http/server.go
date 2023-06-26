package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	ProjectController  controller.Project
	SecretController   controller.Secret
	ResourceController controller.Resource
}

func New(
	projectController controller.Project,
	secretController controller.Secret,
	resourceController controller.Resource,
) *Server {
	return &Server{
		Engine:             gin.Default(),
		ProjectController:  projectController,
		SecretController:   secretController,
		ResourceController: resourceController,
	}
}

func (server *Server) Start() {
	server.WithProjectRouter().WithInternalRouter().WithSecretRouter().WithResourceController()
	err := server.Engine.Run(fmt.Sprintf(":%d", config.Current.Http.Port))
	if err != nil {
		panic(err)
	}
}
