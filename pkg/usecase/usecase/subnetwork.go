package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	repository "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Subnetwork interface {
	List(context.Context, *resource.SubnetworkCollection) errors.Error
	Get(context.Context, *resource.Subnetwork) errors.Error
	Create(context.Context, *resource.Subnetwork) errors.Error
	Update(context.Context, *resource.Subnetwork) errors.Error
	Delete(context.Context, *resource.Subnetwork) errors.Error
}

type subnetworkUseCase struct {
	gcpRepo   repository.Resource
	awsRepo   repository.Resource
	azureRepo repository.Resource
}

func NewSubnetworkUseCase(gcpRepo repository.Resource, awsRepo repository.Resource, azureRepo repository.Resource) Network {
	return &networkUseCase{gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (suc *subnetworkUseCase) List(ctx context.Context, subnetworks *resource.SubnetworkCollection) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (suc *subnetworkUseCase) Get(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (suc *subnetworkUseCase) Create(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (suc *subnetworkUseCase) Update(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (suc *subnetworkUseCase) Delete(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error {
	//TODO implement me
	panic("implement me")
}
