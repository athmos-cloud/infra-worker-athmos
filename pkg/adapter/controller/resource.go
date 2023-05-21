package controller

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	error2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceCtrls "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type Resource interface {
	GetResource(context.Context)
	ListResources(context.Context)
	CreateResource(context.Context)
	UpdateResource(context.Context)
	DeleteResource(context.Context)
}

type resourceController struct {
	providerController   resourceCtrls.Provider
	networkController    resourceCtrls.Network
	subnetworkController resourceCtrls.Subnetwork
}

func NewResourceController(
	providerController resourceCtrls.Provider,
	networkController resourceCtrls.Network,
	subnetworkController resourceCtrls.Subnetwork,
) Resource {
	return &resourceController{
		providerController:   providerController,
		networkController:    networkController,
		subnetworkController: subnetworkController,
	}
}

func (rc *resourceController) GetResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.GetProvider(ctx)
	case types.NetworkResource:
		rc.networkController.GetNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.GetSubnetwork(ctx)
	default:
		error2.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("getting %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) ListResources(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.ListProviders(ctx)
	default:
		error2.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("listing %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) CreateResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.CreateProvider(ctx)
	case types.NetworkResource:
		rc.networkController.CreateNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.CreateSubnetwork(ctx)
	default:
		error2.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("creating %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) UpdateResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.UpdateProvider(ctx)
	case types.NetworkResource:
		rc.networkController.UpdateNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.UpdateSubnetwork(ctx)
	default:
		error2.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("updating %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}

func (rc *resourceController) DeleteResource(ctx context.Context) {
	if err := validator.Resource(ctx); !err.IsOk() {
		error2.RaiseError(ctx, err)
	}
	switch ctx.Value(context.ResourceTypeKey).(types.Resource) {
	case types.ProviderResource:
		rc.providerController.DeleteProvider(ctx)
	case types.NetworkResource:
		rc.networkController.DeleteNetwork(ctx)
	case types.SubnetworkResource:
		rc.subnetworkController.DeleteSubnetwork(ctx)
	default:
		error2.RaiseError(ctx, errors.BadRequest.WithMessage(fmt.Sprintf("deleting %s resource is not supported", ctx.Value(context.ResourceTypeKey).(types.Resource))))
	}
}
