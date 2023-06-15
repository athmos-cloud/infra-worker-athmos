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
	"gopkg.in/mcuadros/go-defaults.v1"
)

type Firewall interface {
	Get(context.Context, *model.Firewall) errors.Error
	Create(context.Context, *model.Firewall) errors.Error
	Update(context.Context, *model.Firewall) errors.Error
	Delete(context.Context, *model.Firewall) errors.Error
}

type firewallUseCase struct {
	projectRepo repository.Project
	gcpRepo     resourceRepo.Resource
	awsRepo     resourceRepo.Resource
	azureRepo   resourceRepo.Resource
}

func NewFirewallUseCase(projectRepo repository.Project, gcpRepo resourceRepo.Resource, awsRepo resourceRepo.Resource, azureRepo resourceRepo.Resource) Firewall {
	return &firewallUseCase{projectRepo: projectRepo, gcpRepo: gcpRepo, awsRepo: awsRepo, azureRepo: azureRepo}
}

func (fuc *firewallUseCase) getRepo(ctx context.Context) resourceRepo.Resource {
	switch ctx.Value(context.ProviderTypeKey).(types.Provider) {
	case types.ProviderGCP:
		return fuc.gcpRepo
		//case types.ProviderAWS:
		//	return fuc.awsRepo
		//case types.ProviderAZURE:
		//	return fuc.azureRepo
	}
	return nil
}

func (fuc *firewallUseCase) Get(ctx context.Context, firewall *model.Firewall) errors.Error {
	repo := fuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s firewall not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.GetFirewallRequest)

	foundFirewall, err := repo.FindFirewall(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Firewall},
	})
	if !err.IsOk() {
		return err
	}
	*firewall = *foundFirewall
	return errors.OK
}

func (fuc *firewallUseCase) Create(ctx context.Context, firewall *model.Firewall) errors.Error {
	repo := fuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s firewall not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.CreateFirewallRequest)
	defaults.SetDefaults(&req)

	network, errNet := repo.FindNetwork(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.ParentID.Network},
	})
	if !errNet.IsOk() {
		return errNet
	}

	toCreateFirewall := &model.Firewall{
		Metadata: metadata.Metadata{
			Managed: req.Managed,
			Tags:    req.Tags,
		},
		IdentifierID: identifier.Firewall{
			Provider: req.ParentID.Provider,
			VPC:      req.ParentID.VPC,
			Network:  req.ParentID.Network,
			Firewall: idFromName(req.Name),
		},
		IdentifierName: identifier.Firewall{
			Provider: network.IdentifierName.Provider,
			VPC:      network.IdentifierName.VPC,
			Network:  network.IdentifierName.Network,
			Firewall: req.Name,
		},
		Allow: req.AllowRules,
		Deny:  req.DenyRules,
	}
	if err := repo.CreateFirewall(ctx, toCreateFirewall); !err.IsOk() {
		return err
	}
	*firewall = *toCreateFirewall

	return errors.Created
}

func (fuc *firewallUseCase) Update(ctx context.Context, firewall *model.Firewall) errors.Error {
	repo := fuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s firewall not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.UpdateFirewallRequest)

	foundFirewall, err := repo.FindFirewall(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Firewall},
	})
	if !err.IsOk() {
		return err
	}
	*firewall = *foundFirewall
	if req.Name != nil {
		firewall.IdentifierName.Firewall = *req.Name
	}
	if req.AllowRules != nil {
		firewall.Allow = *req.AllowRules
	}
	if req.DenyRules != nil {
		firewall.Deny = *req.DenyRules
	}
	if req.Managed != nil {
		firewall.Metadata.Managed = *req.Managed
	}
	if req.Tags != nil {
		firewall.Metadata.Tags = *req.Tags
	}
	if errUpdate := repo.UpdateFirewall(ctx, firewall); !errUpdate.IsOk() {
		return errUpdate
	}

	return errors.NoContent
}

func (fuc *firewallUseCase) Delete(ctx context.Context, firewall *model.Firewall) errors.Error {
	repo := fuc.getRepo(ctx)
	if repo == nil {
		return errors.BadRequest.WithMessage(fmt.Sprintf("%s firewall not supported", ctx.Value(context.ProviderTypeKey).(types.Provider)))
	}
	req := ctx.Value(context.RequestKey).(dto.DeleteFirewallRequest)

	foundFirewall, err := repo.FindFirewall(ctx, option.Option{
		Value: resourceRepo.FindResourceOption{Name: req.IdentifierID.Firewall},
	})
	if !err.IsOk() {
		return err
	}
	*firewall = *foundFirewall
	if errDelete := repo.DeleteFirewall(ctx, firewall); !errDelete.IsOk() {
		return errDelete
	}

	return errors.NoContent
}
