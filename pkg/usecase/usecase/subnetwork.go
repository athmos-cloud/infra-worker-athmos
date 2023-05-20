package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	gcpRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Subnetwork interface {
	List(ctx context.Context, subnetworks *[]resource.Subnetwork) errors.Error
	Get(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error
	Create(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error
	Update(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error
	Delete(ctx context.Context, subnetwork *resource.Subnetwork) errors.Error
}

type subnetworkUseCase struct {
	gcpRepo gcpRepo.Subnetwork
}

func (suc *subnetworkUseCase) List(ctx context.Context, subnetworks *[]resource.Subnetwork) errors.Error {
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

func NewSubnetworkUseCase(gcpRepo gcpRepo.Subnetwork) Subnetwork {
	return &subnetworkUseCase{gcpRepo: gcpRepo}
}
