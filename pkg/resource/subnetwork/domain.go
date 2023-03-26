package subnetwork

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
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

func (subnetwork *Subnetwork) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (subnetwork *Subnetwork) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
