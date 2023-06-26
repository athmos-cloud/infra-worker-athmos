package controller

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Resource interface {
	GetResource(context.Context)
	ListResources(context.Context)
	GetResourceStack(context.Context)
	CreateResource(context.Context)
	UpdateResource(context.Context)
	DeleteResource(context.Context)
}

type resourceController struct {
	providerController   resourceCtrl.Provider
	networkController    resourceCtrl.Network
	subnetworkController resourceCtrl.Subnetwork
	firewallController   resourceCtrl.Firewall
	vmController         resourceCtrl.VM
}

func NewResourceController(
	providerController resourceCtrl.Provider,
	networkController resourceCtrl.Network,
	subnetworkController resourceCtrl.Subnetwork,
	firewallController resourceCtrl.Firewall,
	vmController resourceCtrl.VM,
) Resource {
	return &resourceController{
		providerController:   providerController,
		networkController:    networkController,
		subnetworkController: subnetworkController,
		firewallController:   firewallController,
		vmController:         vmController,
	}
}

func (rc *resourceController) GetResourceStack(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	rc.providerController.GetStack(ctx)
}

func (rc *resourceController) GetResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.GetProvider(ctx)
	case types.NetworkResource:
		rc.networkController.GetNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.GetSubnetwork(ctx)
	case types.FirewallResource:
		rc.firewallController.GetFirewall(ctx)
	case types.VMResource:
		rc.vmController.GetVM(ctx)
	default:
		errorCtrl.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("getting %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) ListResources(ctx context.Context) {
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.ListProviders(ctx)
	default:
		errorCtrl.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("listing %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) CreateResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.CreateProvider(ctx)
	case types.NetworkResource:
		rc.networkController.CreateNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.CreateSubnetwork(ctx)
	case types.FirewallResource:
		rc.firewallController.CreateFirewall(ctx)
	case types.VMResource:
		rc.vmController.CreateVM(ctx)
	default:
		errorCtrl.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("creating %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) UpdateResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.UpdateProvider(ctx)
	case types.NetworkResource:
		rc.networkController.UpdateNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.UpdateSubnetwork(ctx)
	case types.FirewallResource:
		rc.firewallController.UpdateFirewall(ctx)
	case types.VMResource:
		rc.vmController.UpdateVM(ctx)
	default:
		errorCtrl.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("updating %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) DeleteResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.DeleteProvider(ctx)
	case types.NetworkResource:
		rc.networkController.DeleteNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.DeleteSubnetwork(ctx)
	case types.FirewallResource:
		rc.firewallController.DeleteFirewall(ctx)
	case types.VMResource:
		rc.vmController.DeleteVM(ctx)
	default:
		errorCtrl.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("deleting %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}
