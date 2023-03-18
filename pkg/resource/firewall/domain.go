package firewall

import (
	"github.com/PaulBarrie/infra-worker/pkg/application/dto"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/resource"
)

type Firewall struct {
	ID                string `bson:"_id,omitempty"`
	ResourceReference resource.Reference
	Network           string `bson:"network"`
	Allow             []Rule `bson:"allow"`
	Deny              []Rule `bson:"deny"`
}

func (firewall *Firewall) Create(request dto.CreateResourceRequest) (dto.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Update(request dto.UpdateResourceRequest) (dto.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Get(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Watch(request dto.GetResourceRequest) (dto.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) List(request dto.GetListResourceRequest) (dto.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Delete(request dto.DeleteResourceRequest) (dto.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

type Rule struct {
	Protocol string `bson:"protocol"`
	Ports    []int  `bson:"ports"`
}
