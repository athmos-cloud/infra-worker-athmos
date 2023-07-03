package resourceUc

import (
	"fmt"

	"github.com/samber/lo"
	"gopkg.in/mcuadros/go-defaults.v1"

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
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
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
	case types.ProviderAWS:
		return puc.awsRepo
		//case types.ProviderAZURE:
		//	return puc.azureRepo
	}
	return nil
}

func (puc *providerUseCase) List(ctx context.Context, providers *resourceModel.ProviderCollection) errors.Error {
	searchOption := option.Option{Value: resourceRepo.FindAllResourceOption{
		Labels: map[string]string{model.ProjectIDLabelKey: ctx.Value(context.ProjectIDKey).(string)}},
	}
	// gcpProviders, err := puc.gcpRepo.FindAllProviders(ctx, searchOption)
	// if !err.IsOk() {
	// 	return err
	// }
	awsProviders, err := puc.awsRepo.FindAllProviders(ctx, searchOption)
	if !err.IsOk() {
		return err
	}
	if /*len(*gcpProviders) == 0 &&*/ len(*awsProviders) == 0 {
		return errors.NoContent.WithMessage("no providers found")
	}
	*providers = lo.Assign( /**gcpProviders,*/ *awsProviders)

	return errors.OK
}

func (puc *providerUseCase) GetStack(ctx context.Context, provider *resourceModel.Provider) errors.Error {
	repo := puc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s provider not supported", ctx.Value(context.ProviderTypeKey)))
	}
	project, errRepo := puc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		return errRepo
	}
	ctx.Set(context.CurrentNamespace, project.Namespace)
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
	req := ctx.Value(context.RequestKey).(dto.GetResourceRequest)
	defaults.SetDefaults(&req)

	foundProvider, err := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.Identifier}})
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
			Provider: usecase.IdFromName(req.Name),
			VPC:      req.VPC,
		},
		IdentifierName: identifier.Provider{
			Provider: req.Name,
			VPC:      req.VPC,
		},
		Auth: resourceModel.ProviderAuth{
			Name:             secret.Kubernetes.SecretName,
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

	provider, errFind := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID}})
	if !errFind.IsOk() {
		return errFind
	}
	if req.Cascade {
		project, errRepo := puc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
		if !errRepo.IsOk() {
			return errRepo
		}
		ctx.Set(context.CurrentNamespace, project.Namespace)
		if errDelete := repo.DeleteProviderCascade(ctx, provider); !errDelete.IsOk() {
			return errDelete
		}
	} else {
		if errDelete := repo.DeleteProvider(ctx, provider); !errDelete.IsOk() {
			return errDelete
		}
	}

	return errors.NoContent
}
