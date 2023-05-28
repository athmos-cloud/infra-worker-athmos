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
		//case types.ProviderAWS:
		//	return nuc.awsRepo
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
	req := ctx.Value(context.RequestKey).(dto.GetNetworkRequest)
	project, errProject := nuc.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !errProject.IsOk() {
		return errProject
	}
	foundNetwork, err := repo.FindNetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network, Namespace: project.Namespace},
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
	defaults.SetDefaults(&req)

	project, err := nuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}

	var id identifier.Network
	var name identifier.Network

	if req.ParentIDProvider != nil {
		parId := req.ParentIDProvider
		provider, errProvider := repo.FindProvider(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: parId.Provider}})
		if !errProvider.IsOk() {
			return errProvider
		}
		id = identifier.Network{
			Provider: provider.IdentifierID.VPC,
			VPC:      provider.IdentifierID.VPC,
			Network:  idFromName(req.Name),
		}

		name = identifier.Network{
			Network:  req.Name,
			VPC:      provider.IdentifierName.VPC,
			Provider: provider.IdentifierName.Provider,
		}
	} else if req.ParentIDVPC != nil {
		parId := req.ParentIDVPC
		id = identifier.Network{
			Provider: parId.Provider,
			VPC:      parId.VPC,
			Network:  idFromName(req.Name),
		}
		vpc, errProvider := repo.FindVPC(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: parId.VPC, Namespace: ""}})
		if !errProvider.IsOk() {
			return errProvider
		}
		name = identifier.Network{
			Network:  req.Name,
			VPC:      vpc.IdentifierName.VPC,
			Provider: vpc.IdentifierName.Provider,
		}
	} else {
		return errors.BadRequest.WithMessage(fmt.Sprintf("at least one parent id (vpc or provider) must be provided"))
	}

	createdNetwork := &model.Network{
		Metadata: metadata.Metadata{
			Namespace: project.Namespace,
			Managed:   req.Managed,
			Tags:      req.Tags,
		},
		IdentifierID:   id,
		IdentifierName: name,
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

	project, err := nuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !err.IsOk() {
		return err
	}

	curNetwork, err := repo.FindNetwork(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network, Namespace: project.Namespace}})
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
	project, errProj := nuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errProj.IsOk() {
		return errProj
	}
	foundNetwork, errNet := repo.FindNetwork(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network, Namespace: project.Namespace}})
	if !errNet.IsOk() {
		return errNet
	}
	*subnetwork = *foundNetwork
	errDel := repo.DeleteNetwork(ctx, foundNetwork)
	if !errDel.IsOk() {
		return errDel
	}

	return errors.NoContent
}