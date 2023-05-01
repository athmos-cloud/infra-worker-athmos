package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router          *gin.Engine
	ProjectService  *project.Service
	PluginService   *plugin.Service
	ResourceService *resource.Service
	SecretService   *secret.Service
}

func New(
	projectService *project.Service,
	pluginService *plugin.Service,
	resourceService *resource.Service,
	secretService *secret.Service,
) *Server {
	return &Server{
		Router:          gin.Default(),
		ProjectService:  projectService,
		PluginService:   pluginService,
		ResourceService: resourceService,
		SecretService:   secretService,
	}
}

func (server *Server) Start() {
	server.WithProjectRouter().WithPluginController().WithResourceController().WithSecretRouter()
	err := server.Router.Run(fmt.Sprintf(":%d", config.Current.Http.Port))
	if err != nil {
		panic(err)
	}
}
