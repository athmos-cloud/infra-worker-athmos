package resourceUc

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type Provider interface {
	List(context.Context, *model.ProviderCollection) errors.Error
	Get(context.Context, *model.Provider) errors.Error
	Create(context.Context, *model.Provider) errors.Error
	Update(context.Context, *model.Provider) errors.Error
	Delete(context.Context, *model.Provider) errors.Error
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
		//case types.ProviderAWS:
		//	return puc.awsRepo
		//case types.ProviderAZURE:
		//	return puc.azureRepo
	}
	return nil
}

func (puc *providerUseCase) List(ctx context.Context, providers *model.ProviderCollection) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (puc *providerUseCase) Get(ctx context.Context, provider *model.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetProviderRequest)
	defaults.SetDefaults(&req)

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

func (puc *providerUseCase) Create(ctx context.Context, provider *model.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateProviderRequest)
	defaults.SetDefaults(&req)

	project, errProject := puc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errProject.IsOk() {
		return errProject
	}

	secret, errSecret := puc.secretRepo.Find(ctx, option.Option{Value: repository.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
	if !errSecret.IsOk() {
		return errSecret
	}
	providerToCreate := &model.Provider{
		Metadata: metadata.Metadata{
			Namespace: project.Namespace,
			Managed:   true,
		},
		IdentifierID: identifier.Provider{
			Provider: idFromName(req.Name),
			VPC:      req.VPC,
		},
		IdentifierName: identifier.Provider{
			Provider: req.Name,
			VPC:      req.VPC,
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

func (puc *providerUseCase) Update(ctx context.Context, provider *model.Provider) errors.Error {
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
		secret, err := puc.secretRepo.Find(ctx, option.Option{Value: repository.GetSecretByProjectIdAndName{ProjectId: ctx.Value(context.ProjectIDKey).(string), Name: req.SecretAuthName}})
		if !err.IsOk() {
			return err
		}
		provider.Auth = *secret
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

func (puc *providerUseCase) Delete(ctx context.Context, provider *model.Provider) errors.Error {
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
