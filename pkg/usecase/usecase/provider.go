package usecase

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	repository "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Provider interface {
	List(context.Context, *resource.ProviderCollection) errors.Error
	Get(context.Context, *resource.Provider) errors.Error
	Create(context.Context, *resource.Provider) errors.Error
	Update(context.Context, *resource.Provider) errors.Error
	Delete(context.Context, *resource.Provider) errors.Error
}

type providerUseCase struct {
	gcpRepo repository.Resource
}

func NewProviderUseCase(gcpRepo repository.Resource) Provider {
	return &providerUseCase{gcpRepo: gcpRepo}
}

func (puc *providerUseCase) getRepo(ctx context.Context) repository.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return puc.gcpRepo
	}
	return nil
}

func (puc *providerUseCase) List(ctx context.Context, providers *resource.ProviderCollection) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Get(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Create(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Update(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Delete(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}
