package registry

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/presenter"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

func (r *registry) NewSecretController() controller.Secret {
	return controller.NewSecretController(
		usecase.NewSecretUseCase(secret.NewSecretRepository(), secret.NewKubernetesRepository()),
		presenter.NewSecretPresenter(),
	)
}
