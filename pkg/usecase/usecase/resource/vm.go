package resourceUc

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	resourceModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type VM interface {
	Get(context.Context, *resourceModel.VM) errors.Error
	Create(context.Context, *resourceModel.VM) errors.Error
	Update(context.Context, *resourceModel.VM) errors.Error
	Delete(context.Context, *resourceModel.VM) errors.Error
}

type vmUseCase struct {
	projectRepo repository.Project
	sshKeysRepo repository.SSHKeys
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewVMUseCase(projectRepo repository.Project, sshKeysRepo repository.SSHKeys, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) VM {
	return &vmUseCase{projectRepo: projectRepo, sshKeysRepo: sshKeysRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (vuc *vmUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return vuc.gcpRepo
	case types.ProviderAWS:
		return vuc.awsRepo
		//case types.ProviderAZURE:
		//	return vuc.azureRepo
	}
	return nil
}

func (vuc *vmUseCase) Get(ctx context.Context, vm *resourceModel.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetResourceRequest)
	project, errRepo := vuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		return errRepo
	}
	ctx.Set(context.CurrentNamespace, project.Namespace)

	foundVM, err := repo.FindVM(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.Identifier},
	})
	if !err.IsOk() {
		return err
	}
	*vm = *foundVM
	if errKeys := vuc.sshKeysRepo.GetList(ctx, vm.Auths); !errKeys.IsOk() {
		return errKeys
	}

	return errors.OK
}

func (vuc *vmUseCase) Create(ctx context.Context, vm *resourceModel.VM) errors.Error {
	logger.Info.Println("CREAATE")
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateVMRequest)
	defaults.SetDefaults(&req)

	project, errRepo := vuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		logger.Info.Println(errRepo)
		return errRepo
	}

	subnetwork, errNet := repo.FindSubnetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.ParentID.Subnetwork},
	})

	if !errNet.IsOk() {
		logger.Info.Println(errNet)
		return errNet
	}

	var keyList model.SSHKeyList
	for _, auth := range req.Auths {
		keyList = append(keyList, &model.SSHKey{
			Username:        auth.Username,
			KeyLength:       auth.RSAKeyLength,
			SecretName:      usecase.IdFromName(fmt.Sprintf("%s-%s", req.Name, auth.Username)),
			SecretNamespace: project.Namespace,
		})
	}
	if err := vuc.sshKeysRepo.CreateList(ctx, keyList); !err.IsOk() {
		return err
	}

	toCreateVM := &resourceModel.VM{
		Metadata: metadata.Metadata{
			Managed: req.Managed,
			Tags:    req.Tags,
		},
		IdentifierID: identifier.VM{
			Provider:   subnetwork.IdentifierID.Provider,
			VPC:        subnetwork.IdentifierID.VPC,
			Network:    subnetwork.IdentifierID.Network,
			Subnetwork: subnetwork.IdentifierID.Subnetwork,
			VM:         usecase.IdFromName(req.Name),
		},
		IdentifierName: identifier.VM{
			Provider:   subnetwork.IdentifierName.Provider,
			VPC:        subnetwork.IdentifierName.VPC,
			Network:    subnetwork.IdentifierName.Network,
			Subnetwork: subnetwork.IdentifierName.Subnetwork,
			VM:         req.Name,
		},
		AssignPublicIP: req.AssignPublicIP,
		Zone:           req.Zone,
		MachineType:    req.MachineType,
		Auths:          keyList,
		Disks:          req.Disks,
		OS: resourceModel.VMOS{
			ID:   req.OS.ID,
			Name: req.OS.ID,
		},
	}
	if err := repo.CreateVM(ctx, toCreateVM); !err.IsOk() {
		logger.Info.Println(err)
		return err
	}
	logger.Info.Printf("vm created: %s", toCreateVM.IdentifierID.VM)
	*vm = *toCreateVM

	return errors.Created
}

func (vuc *vmUseCase) Update(ctx context.Context, vm *resourceModel.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateVMRequest)
	defaults.SetDefaults(&req)

	foundVM, err := repo.FindVM(ctx, option.Option{Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.VM}})
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
		if errSSH := vuc.updateSSHKeys(ctx, vm.IdentifierName.VM, &vm.Auths, req); !errSSH.IsOk() {
			return errSSH
		}
	} else if req.UpdateSSHKeys {
		if errSSH := vuc.sshKeysRepo.CreateList(ctx, vm.Auths); !errSSH.IsOk() {
			return errSSH
		}
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

func (vuc *vmUseCase) updateSSHKeys(ctx context.Context, vmName string, auths *model.SSHKeyList, req dto.UpdateVMRequest) errors.Error {
	userExists := func(username string) bool {
		for _, auth := range *auths {
			if auth.Username == username {
				return true
			}
		}
		return false
	}
	userRemoved := func(username string) bool {
		for _, auth := range *req.Auths {
			if auth.Username == username {
				return false
			}
		}
		return true
	}
	projectNamespace := (*auths)[0].SecretNamespace
	for _, auth := range *req.Auths {
		if userExists(auth.Username) {
			continue
		}
		key := &model.SSHKey{
			Username:        auth.Username,
			KeyLength:       auth.RSAKeyLength,
			SecretName:      usecase.IdFromName(fmt.Sprintf("%s-%s", vmName, auth.Username)),
			SecretNamespace: projectNamespace,
		}
		if errSSH := vuc.sshKeysRepo.Create(ctx, key); !errSSH.IsOk() {
			return errSSH
		}
		*auths = append(*auths, key)
	}
	for i, auth := range *auths {
		if userRemoved(auth.Username) {
			if errSSH := vuc.sshKeysRepo.Delete(ctx, auth); !errSSH.IsOk() {
				return errSSH
			}
			*auths = append((*auths)[:i], (*auths)[i+1:]...)
		}
	}
	return errors.NoContent
}

func (vuc *vmUseCase) Delete(ctx context.Context, vm *resourceModel.VM) errors.Error {
	repo := vuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s vm not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteVMRequest)
	defaults.SetDefaults(&req)
	project, errRepo := vuc.projectRepo.Find(ctx, option.Option{Value: repository.FindProjectByIDRequest{ID: ctx.Value(context.ProjectIDKey).(string)}})
	if !errRepo.IsOk() {
		return errRepo
	}
	ctx.Set(context.CurrentNamespace, project.Namespace)

	foundVM, err := repo.FindVM(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID},
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
