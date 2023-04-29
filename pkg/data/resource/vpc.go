package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type VPC struct {
	Metadata   metadata.Metadata     `bson:"metadata"`
	Identifier identifier.VPC        `bson:"identifier"`
	Status     status.ResourceStatus `bson:"status"`
	Provider   string                `bson:"provider" plugin:"provider"`
	Networks   NetworkCollection     `bson:"networks"`
}

type VPCCollection map[string]VPC

func (collection *VPCCollection) Equals(other VPCCollection) bool {
	if len(*collection) != len(other) {
		return false
	}
	for key, value := range *collection {
		if !value.Equals(other[key]) {
			return false
		}
	}
	return true
}

func NewVPC(id identifier.VPC, provider common.ProviderType) VPC {
	return VPC{
		Metadata: metadata.Metadata{
			Name: id.ID,
		},
		Identifier: id,
		Status:     status.New(id.ID, common.VPC, provider),
		Networks:   make(NetworkCollection),
	}
}

func (vpc *VPC) New(id identifier.ID, provider common.ProviderType) (IResource, errors.Error) {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.VPC{}) {
		return nil, errors.InvalidArgument.WithMessage("invalid id type")
	}
	res := NewVPC(id.(identifier.VPC), provider)
	return &res, errors.OK
}

func (vpc *VPC) SetMetadata(request metadata.CreateMetadataRequest) {
	vpc.Metadata = metadata.New(request)
}

func (vpc *VPC) GetMetadata() metadata.Metadata {
	return vpc.Metadata
}

func (vpc *VPC) SetStatus(resourceStatus status.ResourceStatus) {
	vpc.Status = resourceStatus
}

func (vpc *VPC) GetStatus() status.ResourceStatus {
	return vpc.Status
}

func (vpc *VPC) GetPluginReference() (resourcePlugin.Reference, errors.Error) {
	if !vpc.Status.PluginReference.ChartReference.Empty() {
		return vpc.Status.PluginReference, errors.OK
	}
	switch vpc.Status.PluginReference.ResourceReference.ProviderType {
	case common.GCP:
		vpc.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VPC.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.VPC.Version,
		}
		return vpc.Status.PluginReference, errors.OK
	}
	return resourcePlugin.Reference{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", vpc.Status.PluginReference.ResourceReference.ProviderType))
}

func (vpc *VPC) FromMap(data map[string]interface{}) errors.Error {
	return resourcePlugin.InjectMapIntoStruct(data, vpc)
}

func (vpc *VPC) Insert(project Project, update ...bool) errors.Error {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	_, ok := project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID]
	if !ok && shouldUpdate {
		return errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", vpc.Identifier.ProviderID))
	}
	if ok && !shouldUpdate {
		return errors.Conflict.WithMessage(fmt.Sprintf("vpc %s in provider %s already exists", vpc.Identifier.ID, vpc.Identifier.ProviderID))
	}
	project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID] = *vpc
	return errors.OK
}

func (vpc *VPC) Remove(project Project) errors.Error {
	_, ok := project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID]
	if !ok {
		return errors.NotFound.WithMessage(fmt.Sprintf("vpc %s in provider %s not found", vpc.Identifier.ID, vpc.Identifier.ProviderID))
	}
	delete(project.Resources[vpc.Identifier.ProviderID].VPCs, vpc.Identifier.ID)
	return errors.NoContent
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.Metadata.Equals(other.Metadata) &&
		vpc.Identifier.Equals(other.Identifier) &&
		vpc.Status.Equals(other.Status) &&
		vpc.Provider == other.Provider &&
		vpc.Networks.Equals(other.Networks)
}
