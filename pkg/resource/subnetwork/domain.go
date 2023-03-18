package subnetwork

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource"
	"github.com/PaulBarrie/infra-worker/pkg/resource/vm"
)

type Subnetwork struct {
	ID                string `bson:"_id,omitempty"`
	ResourceReference resource.Reference
	Region            string  `bson:"region"`
	IPCIDRRange       string  `bson:"ipCidrRange"`
	VMs               []vm.VM `bson:"vmList"`
}

func (subnetwork *Subnetwork) Create(request dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Update(request dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Get(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Watch(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) List(request dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Delete(request dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
