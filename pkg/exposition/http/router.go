package http

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/application/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router          *gin.Engine
	ProjectService  *application.ProjectService
	PluginService   *application.PluginService
	ResourceService *application.ResourceService
	SecretService   *secret.Service
}

func New(
	projectService *application.ProjectService,
	pluginService *application.PluginService,
	resourceService *application.ResourceService,
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
