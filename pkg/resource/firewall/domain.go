package firewall

import (
	resource2 "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type Firewall struct {
	ID      string `bson:"_id,omitempty"`
	Network string `bson:"network"`
	Allow   []Rule `bson:"allow"`
	Deny    []Rule `bson:"deny"`
}

func (firewall *Firewall) Create(request resource2.CreateResourceRequest) (resource2.CreateResourceResponse, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Update(request resource2.UpdateResourceRequest) (resource2.UpdateResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Get(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Watch(request resource2.GetResourceRequest) (resource2.GetResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) List(request resource2.GetListResourceRequest) (resource2.GetListResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (firewall *Firewall) Delete(request resource2.DeleteResourceRequest) (resource2.DeleteResourceRequest, errors.Error) {
	//TODO implement me
	panic("implement me")
}

type Rule struct {
	Protocol string `bson:"protocol"`
	Ports    []int  `bson:"ports"`
}
