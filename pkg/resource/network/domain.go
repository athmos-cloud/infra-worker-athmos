package network

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
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

func (network *Network) Create(request dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Update(request dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Get(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Watch(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) List(request dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (network *Network) Delete(request dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}
