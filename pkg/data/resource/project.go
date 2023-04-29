package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	"reflect"
)

type Project struct {
	ID        string             `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Namespace string             `bson:"namespace"`
	OwnerID   string             `bson:"owner_id"`
	Resources ProviderCollection `bson:"providers"`
}

func NewProject(name string, ownerID string) Project {
	return Project{
		Name:      name,
		Namespace: fmt.Sprintf("%s-%s", name, utils.RandomString(5)),
		OwnerID:   ownerID,
		Resources: make(ProviderCollection, 10000),
	}
}

func (project *Project) Insert(resource IResource) errors.Error {
	return resource.Insert(*project)
}

func (project *Project) Update(resource IResource) errors.Error {
	return resource.Insert(*project, true)
}

func (project *Project) Delete(resource IResource) errors.Error {
	return resource.Remove(*project)
}

func (project *Project) Get(id identifier.ID) (IResource, errors.Error) {
	switch reflect.TypeOf(id) {
	case reflect.TypeOf(identifier.Provider{}):
		providerID := id.(identifier.Provider)
		provider, ok := project.Resources[providerID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", providerID.ID))
		}
		return &provider, errors.OK
	case reflect.TypeOf(identifier.VPC{}):
		vpcID := id.(identifier.VPC)
		vpc, ok := project.Resources[vpcID.ProviderID].VPCs[vpcID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vpc %s not found in provider %s", vpcID.ID, vpcID.ProviderID))
		}
		return &vpc, errors.OK
	case reflect.TypeOf(identifier.Network{}):
		networkID := id.(identifier.Network)
		network, ok := project.Resources[networkID.ProviderID].VPCs[networkID.VPCID].Networks[networkID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("network %s not found in vpc %s", networkID.ID, networkID.VPCID))
		}
		return &network, errors.OK
	case reflect.TypeOf(identifier.Subnetwork{}):
		subnetID := id.(identifier.Subnetwork)
		subnet, ok := project.Resources[subnetID.ProviderID].VPCs[subnetID.VPCID].Networks[subnetID.NetworkID].Subnetworks[subnetID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("subnet %s not found in network %s", subnetID.ID, subnetID.NetworkID))
		}
		return &subnet, errors.OK
	case reflect.TypeOf(identifier.Firewall{}):
		firewallID := id.(identifier.Firewall)
		firewall, ok := project.Resources[firewallID.ProviderID].VPCs[firewallID.VPCID].Networks[firewallID.NetworkID].Firewalls[firewallID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("firewall %s not found in network %s", firewallID.ID, firewallID.NetworkID))
		}
		return &firewall, errors.OK
	case reflect.TypeOf(identifier.VM{}):
		vmID := id.(identifier.VM)
		vm, ok := project.Resources[vmID.ProviderID].VPCs[vmID.VPCID].Networks[vmID.NetworkID].Subnetworks[vmID.ID]
		if !ok {
			return nil, errors.NotFound.WithMessage(fmt.Sprintf("vm %s not found in subnet %s", vmID.ID, vmID.SubnetID))
		}
		return &vm, errors.OK
	}
	return nil, errors.InvalidArgument.WithMessage("invalid id type")
}
