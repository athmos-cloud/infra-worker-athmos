package resourceUc

import (
	"fmt"

	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
)

type Network interface {
	Get(context.Context, *model.Network) errors.Error
	Create(context.Context, *model.Network) errors.Error
	Update(context.Context, *model.Network) errors.Error
	Delete(context.Context, *model.Network) errors.Error
}

type networkUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewNetworkUseCase(projectRepo repository.Project, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) Network {
	return &networkUseCase{projectRepo: projectRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (nuc *networkUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return nuc.gcpRepo
	case types.ProviderAWS:
		return nuc.awsRepo
		//case types.ProviderAZURE:
		//	return nuc.azureRepo
	}
	return nil
}

func (nuc *networkUseCase) Get(ctx context.Context, network *model.Network) errors.Error {
	repo := nuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s network not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetResourceRequest)

	foundNetwork, err := repo.FindNetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.Identifier},
	})
	if !err.IsOk() {
		return err
	}
	*network = *foundNetwork
	return errors.OK
}

func (nuc *networkUseCase) Create(ctx context.Context, network *model.Network) errors.Error {
	repo := nuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s createdNetwork not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateNetworkRequest)
	parId := req.ParentIDProvider
	provider, errProvider := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: parId.Provider}})
	if !errProvider.IsOk() {
		return errProvider
	}

	id := identifier.Network{
		Provider: provider.IdentifierID.Provider,
		VPC:      provider.IdentifierID.VPC,
		Network:  usecase.IdFromName(req.Name),
	}
	name := identifier.Network{
		Network:  req.Name,
		VPC:      provider.IdentifierName.VPC,
		Provider: provider.IdentifierName.Provider,
	}
	createdNetwork := &model.Network{
		Metadata: metadata.Metadata{
			Managed: req.Managed,
			Tags:    req.Tags,
		},
		IdentifierID:   id,
		IdentifierName: name,
		Region:         req.Region,
	}
	if errCreate := repo.CreateNetwork(ctx, createdNetwork); !errCreate.IsOk() {
		return errCreate
	}
	*network = *createdNetwork
	return errors.Created

}

func (nuc *networkUseCase) Update(ctx context.Context, network *model.Network) errors.Error {
	repo := nuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s network not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateNetworkRequest)

	curNetwork, err := repo.FindNetwork(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network}})
	if !err.IsOk() {
		return err
	}
	*network = *curNetwork
	if req.Managed != nil {
		network.Metadata.Managed = *req.Managed
	}
	if req.Tags != nil {
		network.Metadata.Tags = *req.Tags
	}

	if errUpdate := repo.UpdateNetwork(ctx, network); !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (nuc *networkUseCase) Delete(ctx context.Context, subnetwork *model.Network) errors.Error {
	repo := nuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s network not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteNetworkRequest)

	

	foundNetwork, errNet := repo.FindNetwork(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID}})
	if !errNet.IsOk() {
		return errNet
	}
	*subnetwork = *foundNetwork
	if req.Cascade {
		project, errRepo := nuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
		if !errRepo.IsOk() {
			return errRepo
		}
		ctx.Set(context.CurrentNamespace, project.Namespace)
		if errCascade := repo.DeleteNetworkCascade(ctx, foundNetwork); !errCascade.IsOk() {
			return errCascade
		}
	} else {
		if errDel := repo.DeleteNetwork(ctx, foundNetwork); !errDel.IsOk() {
			return errDel
		}
	}

	return errors.NoContent
}
