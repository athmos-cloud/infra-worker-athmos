package resourceUc

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	resourceModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/secret"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type Provider interface {
	List(context.Context, *resourceModel.ProviderCollection) errors.Error
	Get(context.Context, *resourceModel.Provider) errors.Error
	GetStack(context.Context, *resourceModel.Provider) errors.Error
	Create(context.Context, *resourceModel.Provider) errors.Error
	Update(context.Context, *resourceModel.Provider) errors.Error
	Delete(context.Context, *resourceModel.Provider) errors.Error
}

type providerUseCase struct {
	projectRepo repository.Project
	secretRepo  secret.Secret
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewProviderUseCase(projectRepo repository.Project, secretRepo secret.Secret, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) Provider {
	return &providerUseCase{projectRepo: projectRepo, secretRepo: secretRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (puc *providerUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return puc.gcpRepo
		//case types.ProviderAWS:
		//	return puc.awsRepo
		//case types.ProviderAZURE:
		//	return puc.azureRepo
	}
	return nil
}

func (puc *providerUseCase) List(ctx context.Context, providers *resourceModel.ProviderCollection) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}

	searchOption := option.Option{Value: resourceRepo.FindAllResourceOption{
		Labels: map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}},
	}

	foundProviders, err := repo.FindAllProviders(ctx, searchOption)
	if !err.IsOk() {
		return err
	}
	if len(*foundProviders) == 0 {
		return errors.NotFound.WithMessage("no providers found")
	}
	*providers = *foundProviders
	return errors.OK
}

func (puc *providerUseCase) GetStack(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetProviderStackRequest)
	foundProvider, err := repo.FindProviderStack(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.ProviderID}})
	if !err.IsOk() {
		return err
	}
	*provider = *foundProvider
	return errors.OK
}

func (puc *providerUseCase) Get(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetProviderRequest)
	defaults.SetDefaults(&req)

	foundProvider, err := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Provider}})
	if !err.IsOk() {
		return err
	}
	*provider = *foundProvider
	return errors.OK
}

func (puc *providerUseCase) Create(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateProviderRequest)
	defaults.SetDefaults(&req)

	secret, errSecret := puc.secretRepo.Find(ctx, option.Option{Value: secret.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
	if !errSecret.IsOk() {
		return errSecret
	}
	providerToCreate := &resourceModel.Provider{
		IdentifierID: identifier.Provider{
			Provider: idFromName(req.Name),
			VPC:      req.VPC,
		},
		IdentifierName: identifier.Provider{
			Provider: req.Name,
			VPC:      req.VPC,
		},
		Auth: resourceModel.ProviderAuth{
			Name:             secret.Name,
			KubernetesSecret: secret.Kubernetes,
		},
	}

	errCreate := repo.CreateProvider(ctx, providerToCreate)
	if !errCreate.IsOk() {
		return errCreate
	}
	*provider = *providerToCreate

	return errors.Created
}

func (puc *providerUseCase) Update(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateProviderRequest)
	defaults.SetDefaults(&req)

	provider, errFind := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Provider}})
	if !errFind.IsOk() {
		return errFind
	}
	if req.SecretAuthName != "" {
		secret, err := puc.secretRepo.Find(ctx, option.Option{Value: secret.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
		if !err.IsOk() {
			return errors.BadRequest.WithMessage(fmt.Sprintf("secret %s does not exist", req.SecretAuthName))
		}
		provider.Auth = resourceModel.ProviderAuth{
			Name:             secret.Name,
			KubernetesSecret: secret.Kubernetes,
		}
	}
	if req.Name != "" {
		provider.IdentifierName.Provider = req.Name
	}
	errUpdate := repo.UpdateProvider(ctx, provider)
	if !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (puc *providerUseCase) Delete(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteProviderRequest)
	defaults.SetDefaults(&req)

	provider, errFind := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Provider}})
	if !errFind.IsOk() {
		return errFind
	}
	errDelete := repo.DeleteProvider(ctx, provider)
	if !errDelete.IsOk() {
		return errDelete
	}

	return errors.NoContent
}
