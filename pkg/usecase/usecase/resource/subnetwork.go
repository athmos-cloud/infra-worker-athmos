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
	"gopkg.in/mcuadros/go-defaults.v1"
)

type Subnetwork interface {
	Get(context.Context, *model.Subnetwork) errors.Error
	Create(context.Context, *model.Subnetwork) errors.Error
	Update(context.Context, *model.Subnetwork) errors.Error
	Delete(context.Context, *model.Subnetwork) errors.Error
}

type subnetworkUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewSubnetworkUseCase(projectRepo repository.Project, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) Subnetwork {
	return &subnetworkUseCase{projectRepo: projectRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (suc *subnetworkUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return suc.gcpRepo
	case types.ProviderAWS:
		return suc.awsRepo
		//case types.ProviderAZURE:
		//	return suc.azureRepo
	}
	return nil
}

func (suc *subnetworkUseCase) Get(ctx context.Context, subnetwork *model.Subnetwork) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s subnetwork not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetResourceRequest)
	defaults.SetDefaults(&req)

	findSubnetwork, err := repo.FindSubnetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.Identifier},
	})
	if !err.IsOk() {
		return err
	}
	*subnetwork = *findSubnetwork

	return errors.OK
}

func (suc *subnetworkUseCase) Create(ctx context.Context, subnetwork *model.Subnetwork) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s subnetwork not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateSubnetworkRequest)
	defaults.SetDefaults(&req)

	network, errNet := repo.FindNetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.ParentID.Network},
	})
	if !errNet.IsOk() {
		return errNet
	}

	toCreateSubnet := &model.Subnetwork{
		Metadata: metadata.Metadata{
			Managed: req.Managed,
			Tags:    req.Tags,
		},
		IdentifierID: identifier.Subnetwork{
			Provider:   network.IdentifierID.Provider,
			VPC:        network.IdentifierID.VPC,
			Network:    network.IdentifierID.Network,
			Subnetwork: usecase.IdFromName(req.Name),
		},
		IdentifierName: identifier.Subnetwork{
			Provider:   network.IdentifierName.Provider,
			VPC:        network.IdentifierName.VPC,
			Network:    network.IdentifierName.Network,
			Subnetwork: req.Name,
		},
		Region:      req.Region,
		IPCIDRRange: req.IPCIDRRange,
	}

	if errSubnet := repo.CreateSubnetwork(ctx, toCreateSubnet); !errSubnet.IsOk() {
		return errSubnet
	}
	*subnetwork = *toCreateSubnet

	return errors.Created
}

func (suc *subnetworkUseCase) Update(ctx context.Context, subnetwork *model.Subnetwork) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s subnetwork not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateSubnetworkRequest)
	defaults.SetDefaults(&req)

	foundSubnetwork, err := repo.FindSubnetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Subnetwork},
	})
	if !err.IsOk() {
		return err
	}

	*subnetwork = *foundSubnetwork
	if req.Tags != nil {
		subnetwork.Metadata.Tags = *req.Tags
	}
	if req.Region != nil {
		subnetwork.Region = *req.Region
	}
	if req.IPCIDRRange != nil {
		subnetwork.IPCIDRRange = *req.IPCIDRRange
	}

	if errUpdate := repo.UpdateSubnetwork(ctx, subnetwork); !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (suc *subnetworkUseCase) Delete(ctx context.Context, subnetwork *model.Subnetwork) errors.Error {
	repo := suc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s subnetwork not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteSubnetworkRequest)
	defaults.SetDefaults(&req)

	foundNetwork, err := repo.FindSubnetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID},
	})
	if !err.IsOk() {
		return err
	}
	*subnetwork = *foundNetwork
	if req.Cascade {
		project, errRepo := suc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
		if !errRepo.IsOk() {
			return errRepo
		}
		ctx.Set(context.CurrentNamespace, project.Namespace)
		if delErr := repo.DeleteSubnetworkCascade(ctx, subnetwork); !delErr.IsOk() {
			return delErr
		}
	} else {
		if delErr := repo.DeleteSubnetwork(ctx, subnetwork); !delErr.IsOk() {
			return delErr
		}
	}

	return errors.NoContent
}
