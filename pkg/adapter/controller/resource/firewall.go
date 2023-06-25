package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	model "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	resource2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type Firewall interface {
	GetFirewall(context.Context)
	CreateFirewall(context.Context)
	UpdateFirewall(context.Context)
	DeleteFirewall(context.Context)
}

type firewallController struct {
	firewallUseCase resource2.Firewall
	firewallOutput  output.FirewallPort
}

func NewFirewallController(firewallUseCase resource2.Firewall, firewallOutput output.FirewallPort) Firewall {
	return &firewallController{firewallUseCase: firewallUseCase, firewallOutput: firewallOutput}
}

func (nc *firewallController) GetFirewall(ctx context.Context) {
	if err := resourceValidator.GetFirewall(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	firewall := &model.Firewall{}
	if err := nc.firewallUseCase.Get(ctx, firewall); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.firewallOutput.Render(ctx, firewall)
	}
}

func (nc *firewallController) CreateFirewall(ctx context.Context) {
	if err := resourceValidator.CreateFirewall(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	firewall := &model.Firewall{}
	if err := nc.firewallUseCase.Create(ctx, firewall); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.firewallOutput.RenderCreate(ctx, firewall)
	}
}

func (nc *firewallController) UpdateFirewall(ctx context.Context) {
	if err := resourceValidator.UpdateFirewall(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	firewall := &model.Firewall{}
	if err := nc.firewallUseCase.Update(ctx, firewall); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.firewallOutput.RenderUpdate(ctx, firewall)
	}
}

func (nc *firewallController) DeleteFirewall(ctx context.Context) {
	if err := resourceValidator.DeleteFirewall(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	firewall := &model.Firewall{}
	if err := nc.firewallUseCase.Delete(ctx, firewall); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		nc.firewallOutput.Render(ctx, firewall)
	}
}
