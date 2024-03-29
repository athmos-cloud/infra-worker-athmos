package registry

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/presenter"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

func (r *registry) NewProjectController() controller.Project {
	return controller.NewProjectController(
		usecase.NewProjectUseCase(repository.NewProjectRepository(), gcp.NewRepository(), aws.NewRepository()),
		presenter.NewProjectPresenter(),
	)
}
