package registry

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	resourceCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/resource"
	presenter "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/presenter/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	gcpRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	secretRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

func (r *registry) NewResourceController() controller.Resource {
	projectRepo := repository.NewProjectRepository()
	secretRepo := secretRepository.NewSecretRepository()
	gcpRepo := gcpRepository.NewRepository()

	providerCtrl := resourceCtrl.NewProviderController(
		usecase.NewProviderUseCase(projectRepo, secretRepo, gcpRepo, gcpRepo, gcpRepo),
		presenter.NewProviderPresenter(),
	)
	networkCtrl := resourceCtrl.NewNetworkController(
		usecase.NewNetworkUseCase(projectRepo, gcpRepo, gcpRepo, gcpRepo),
		presenter.NewNetworkPresenter(),
	)
	subnetworkCtrl := resourceCtrl.NewSubnetworkController(
		usecase.NewSubnetworkUseCase(projectRepo, gcpRepo, gcpRepo, gcpRepo),
		presenter.NewSubnetworkPresenter(),
	)
	firewallCtrl := resourceCtrl.NewFirewallController(
		usecase.NewFirewallUseCase(projectRepo, gcpRepo, gcpRepo, gcpRepo),
		presenter.NewFirewallPresenter(),
	)
	vmCtrl := resourceCtrl.NewVMController(
		usecase.NewVMUseCase(projectRepo, gcpRepo, gcpRepo, gcpRepo),
		presenter.NewVMPresenter(),
	)

	return controller.NewResourceController(
		providerCtrl,
		networkCtrl,
		subnetworkCtrl,
		firewallCtrl,
		vmCtrl,
	)
}
