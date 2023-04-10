package http

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/service/plugin"
	"github.com/PaulBarrie/infra-worker/pkg/service/project"
	"github.com/PaulBarrie/infra-worker/pkg/service/resource"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router          *gin.Engine
	ProjectService  *project.Service
	PluginService   *plugin.Service
	ResourceService *resource.Service
}

func New(projectService *project.Service, pluginService *plugin.Service, resourceService *resource.Service) *Server {
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
