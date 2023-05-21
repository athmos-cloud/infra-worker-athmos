package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
)

type subnetwork struct{}

func NewSubnetworkPresenter() output.SubnetworkPort {
	return &subnetwork{}
}

func (s *subnetwork) Render(cxt context.Context, subnetwork *resource.Subnetwork) {
	//TODO implement me
	panic("implement me")
}

func (s *subnetwork) RenderCreate(cxt context.Context, subnetwork *resource.Subnetwork) {
	//TODO implement me
	panic("implement me")
}

func (s *subnetwork) RenderUpdate(cxt context.Context, subnetwork *resource.Subnetwork) {
	//TODO implement me
	panic("implement me")
}

func (s *subnetwork) RenderDelete(cxt context.Context, subnetwork *resource.Subnetwork) {
	//TODO implement me
	panic("implement me")
}
