package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

func (gcp *gcpRepository) FindFirewall(ctx context.Context, opt option.Option) (*resource.Firewall, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) FindAllRecursiveFirewalls(ctx context.Context, opt option.Option) (*resource.FirewallCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) CreateFirewall(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) UpdateFirewall(ctx context.Context, provider *resource.Provider) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (gcp *gcpRepository) DeleteFirewall(ctx context.Context, firewall *resource.Firewall) errors.Error {
	//TODO implement me
	panic("implement me")
}
