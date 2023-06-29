package registry

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller"
	resourceCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/resource"
	presenter "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/presenter/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	gcpRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	secretRepository "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

func (r *registry) NewResourceController() controller.Resource {
	projectRepo := repository.NewProjectRepository()
	secretRepo := secretRepository.NewSecretRepository()
	gcpRepo := gcpRepository.NewRepository()
	awsRepo := aws.NewRepository()
	sshKeysRepo := repository.NewSSHKeyRepository()

	providerCtrl := resourceCtrl.NewProviderController(
		usecase.NewProviderUseCase(projectRepo, secretRepo, gcpRepo, awsRepo, nil),
		presenter.NewProviderPresenter(),
	)
	networkCtrl := resourceCtrl.NewNetworkController(
		usecase.NewNetworkUseCase(projectRepo, gcpRepo, awsRepo, nil),
		presenter.NewNetworkPresenter(),
	)
	subnetworkCtrl := resourceCtrl.NewSubnetworkController(
		usecase.NewSubnetworkUseCase(projectRepo, gcpRepo, awsRepo, nil),
		presenter.NewSubnetworkPresenter(),
	)
	firewallCtrl := resourceCtrl.NewFirewallController(
		usecase.NewFirewallUseCase(projectRepo, gcpRepo, awsRepo, nil),
		presenter.NewFirewallPresenter(),
	)
	vmCtrl := resourceCtrl.NewVMController(
		usecase.NewVMUseCase(projectRepo, sshKeysRepo, gcpRepo, awsRepo, nil),
		presenter.NewVMPresenter(),
	)
	sqlCtrl := resourceCtrl.NewSqlDBController(
		usecase.NewSqlDBUseCase(projectRepo, gcpRepo, awsRepo, nil),
		presenter.NewSqlDBPresenter(),
	)
	return controller.NewResourceController(
		providerCtrl,
		networkCtrl,
		subnetworkCtrl,
		firewallCtrl,
		vmCtrl,
		sqlCtrl,
	)
}
