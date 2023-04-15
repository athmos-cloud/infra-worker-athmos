package http

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/application"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router          *gin.Engine
	ProjectService  *application.ProjectService
	PluginService   *application.PluginService
	ResourceService *application.ResourceService
}

func New(projectService *application.ProjectService, pluginService *application.PluginService, resourceService *application.ResourceService) *Server {
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
