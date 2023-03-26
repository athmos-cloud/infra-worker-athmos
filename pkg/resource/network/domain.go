package network

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource/firewall"
	"github.com/PaulBarrie/infra-worker/pkg/resource/subnetwork"
	"sigs.k8s.io/kustomize/api/resource"
)

type Network struct {
	ID                string                  `bson:"_id,omitempty"`
	ResourceReference resource.Resource       `bson:"resourceReference"`
	Subnetworks       []subnetwork.Subnetwork `bson:"subnetworks"`
	Firewalls         []firewall.Firewall     `bson:"firewalls"`
}

func (network *Network) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
