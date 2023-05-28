package registry

import "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"

type registry struct{}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		Project:  r.NewProjectController(),
		Secret:   r.NewSecretController(),
		Resource: r.NewResourceController(),
	}
}
