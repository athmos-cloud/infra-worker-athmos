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

type VM interface {
	Get(context.Context, *model.VM) errors.Error
	Create(context.Context, *model.VM) errors.Error
	Update(context.Context, *model.VM) errors.Error
	Delete(context.Context, *model.VM) errors.Error
}

type vmUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewVMUseCase(gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) VM {
	return &vmUseCase{gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (vuc *vmUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return vuc.gcpRepo
		//case types.ProviderAWS:
		//	return vuc.awsRepo
		//case types.ProviderAZURE:
		//	return vuc.azureRepo
	}
	return nil
}

func (vuc *vmUseCase) Get(ctx context.Context, vm *model.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetVMRequest)
	project, errProject := vuc.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !errProject.IsOk() {
		return errProject
	}
	foundVM, err := repo.FindVM(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.VM, Namespace: project.Namespace},
	})
	if !err.IsOk() {
		return err
	}
	*vm = *foundVM

	return errors.OK
}

func (vuc *vmUseCase) Create(ctx context.Context, vm *model.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateVMRequest)
	defaults.SetDefaults(&req)

	project, errRepo := vuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		return errRepo
	}

	network, errNet := repo.FindSubnetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.ParentID.Subnetwork, Namespace: project.Namespace},
	})
	if !errNet.IsOk() {
		return errNet
	}

	toCreateVM := &model.VM{
		Metadata: metadata.Metadata{
			Namespace: project.Namespace,
			Managed:   *req.Managed,
			Tags:      req.Tags,
		},
		IdentifierID: identifier.VM{
			Provider:   req.ParentID.Provider,
			VPC:        req.ParentID.VPC,
			Network:    req.ParentID.Network,
			Subnetwork: req.ParentID.Subnetwork,
			VM:         idFromName(req.Name),
		},
		IdentifierName: identifier.VM{
			Provider:   network.IdentifierName.Provider,
			VPC:        network.IdentifierName.VPC,
			Network:    network.IdentifierName.Network,
			Subnetwork: network.IdentifierName.Subnetwork,
			VM:         req.Name,
		},
		AssignPublicIP: req.AssignPublicIP,
		Zone:           req.Zone,
		MachineType:    req.MachineType,
		Auths:          req.Auths,
		Disks:          req.Disks,
		OS:             req.OS,
	}
	if err := repo.CreateVM(ctx, toCreateVM); !err.IsOk() {
		return err
	}
	*vm = *toCreateVM

	return errors.Created
}

func (vuc *vmUseCase) Update(ctx context.Context, vm *model.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateVMRequest)
	defaults.SetDefaults(&req)
	project, errProject := vuc.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !errProject.IsOk() {
		return errProject
	}
	foundVM, err := repo.FindVM(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network, Namespace: project.Namespace},
	})
	if !err.IsOk() {
		return err
	}

	*vm = *foundVM
	if req.Name != nil {
		vm.IdentifierName.VM = *req.Name
	}
	if req.Tags != nil {
		vm.Metadata.Tags = *req.Tags
	}
	if req.Managed != nil {
		vm.Metadata.Managed = *req.Managed
	}
	if req.AssignPublicIP != nil {
		vm.AssignPublicIP = *req.AssignPublicIP
	}
	if req.Zone != nil {
		vm.Zone = *req.Zone
	}
	if req.MachineType != nil {
		vm.MachineType = *req.MachineType
	}
	if req.Auths != nil {
		vm.Auths = *req.Auths
	}
	if req.Disks != nil {
		vm.Disks = *req.Disks
	}
	if req.OS != nil {
		vm.OS = *req.OS
	}
	if errUpdate := repo.UpdateVM(ctx, vm); !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (vuc *vmUseCase) Delete(ctx context.Context, vm *model.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteVMRequest)
	defaults.SetDefaults(&req)
	project, errProject := vuc.projectRepo.Find(ctx, option.Option{
		Value: repository.FindProjectByIDRequest{
			ID: ctx.Value(context.ProjectIDKey).(string),
		},
	})
	if !errProject.IsOk() {
		return errProject
	}
	foundVM, err := repo.FindVM(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Network, Namespace: project.Namespace},
	})
	if !err.IsOk() {
		return err
	}
	*vm = *foundVM

	if delErr := repo.DeleteVM(ctx, vm); !delErr.IsOk() {
		return delErr
	}

	return errors.NoContent
}
