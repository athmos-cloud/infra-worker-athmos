package usecase

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type Provider interface {
	List(context.Context, *resource.ProviderCollection) errors.Error
	Get(context.Context, *resource.Provider) errors.Error
	Create(context.Context, *resource.Provider) errors.Error
	Update(context.Context, *resource.Provider) errors.Error
	Delete(context.Context, *resource.Provider) errors.Error
}

type providerUseCase struct {
	projectRepo repository.Project
	secretRepo  repository.Secret
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewProviderUseCase(projectRepo repository.Project, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) Provider {
	return &providerUseCase{projectRepo: projectRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (puc *providerUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return puc.gcpRepo
	case types.ProviderAWS:
		return puc.awsRepo
	case types.ProviderAZURE:
		return puc.azureRepo
	}
	return nil
}

func (puc *providerUseCase) List(ctx context.Context, providers *resource.ProviderCollection) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Get(ctx context.Context, provider *resource.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("provider %s not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetProviderRequest)
	project, err := puc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}
	foundProvider, err := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.Identifier.Provider, Namespace: project.Namespace}})
	if !err.IsOk() {
		return err
	}
	*provider = *foundProvider
	return errors.OK
}

func (puc *providerUseCase) Create(ctx context.Context, provider *resource.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("provider %s not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateProviderRequest)
	project, errProject := puc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errProject.IsOk() {
		return errProject
	}

	secret, errSecret := puc.secretRepo.Find(ctx, option.Option{Value: repository.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
	if !errSecret.IsOk() {
		return errSecret
	}
	providerToCreate := &resource.Provider{
		Metadata: metadata.Metadata{
			Namespace: project.Namespace,
			Managed:   true,
		},
		IdentifierID: identifier.Provider{
			Provider: idFromName(req.Name),
		},
		IdentifierName: identifier.Provider{
			Provider: req.Name,
		},
		Auth: *secret,
	}
	errCreate := repo.CreateProvider(ctx, providerToCreate)
	if !errCreate.IsOk() {
		return errCreate
	}
	*provider = *providerToCreate
	return errors.Created
}

func (puc *providerUseCase) Update(ctx context.Context, provider *resource.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("provider %s not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateProviderRequest)
	provider, errFind := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Provider}})
	if !errFind.IsOk() {
		return errFind
	}
	if req.SecretAuthName != "" {
		secret, err := puc.secretRepo.Find(ctx, option.Option{Value: repository.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
		if !err.IsOk() {
			return err
		}
		provider.Auth = *secret
	}
	if req.Name != "" {
		provider.IdentifierName.Provider = req.Name
	}
	errCreate := repo.UpdateProvider(ctx, provider)
	if !errCreate.IsOk() {
		return errCreate
	}
	return errors.Created
}

func (puc *providerUseCase) Delete(ctx context.Context, provider *resource.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("provider %s not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteProviderRequest)
	provider, errFind := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Provider}})
	if !errFind.IsOk() {
		return errFind
	}
	errDelete := repo.DeleteProvider(ctx, provider)
	if !errDelete.IsOk() {
		return errDelete
	}
	return errors.Created
}
