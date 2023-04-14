package http

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router          *gin.Engine
	ProjectService  *service.ProjectService
	PluginService   *service.PluginService
	ResourceService *service.ResourceService
}

func New(projectService *service.ProjectService, pluginService *service.PluginService, resourceService *service.ResourceService) *Server {
	return &Server{
		Router:          gin.Default(),
		ProjectService:  projectService,
		PluginService:   pluginService,
		ResourceService: resourceService,
	}
}

func (server *Server) Start() {
	server.WithProjectRouter().WithPluginController().WithResourceController()
	err := server.Router.Run(fmt.Sprintf(":%d", config.Current.Http.Port))
	if err != nil {
		panic(err)
	}
}
