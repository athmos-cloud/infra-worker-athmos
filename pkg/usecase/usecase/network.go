package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	gcpRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Network interface {
	List(context.Context, *resource.NetworkCollection) errors.Error
	Get(context.Context, *resource.Network) errors.Error
	Create(context.Context, *resource.Network) errors.Error
	Update(context.Context, *resource.Network) errors.Error
	Delete(context.Context, *resource.Network) errors.Error
}

type networkUseCase struct {
	gcpRepo gcpRepo.Network
}

func NewNetworkUseCase(gcpRepo gcpRepo.Network) Network {
	return &networkUseCase{gcpRepo: gcpRepo}
}

func (nuc *networkUseCase) List(ctx context.Context, subnetworks *resource.NetworkCollection) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (nuc *networkUseCase) Get(ctx context.Context, subnetwork *resource.Network) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (nuc *networkUseCase) Create(ctx context.Context, subnetwork *resource.Network) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (nuc *networkUseCase) Update(ctx context.Context, subnetwork *resource.Network) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (nuc *networkUseCase) Delete(ctx context.Context, subnetwork *resource.Network) errors.Error {
	//TODO implement me
	panic("implement me")
}
