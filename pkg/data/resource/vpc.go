package resource

import (
	"fmt"
	resourcePlugin "github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

type VPC struct {
	Metadata   metadata.Metadata     `bson:"metadata"`
	Identifier identifier.VPC        `bson:"identifier"`
	Status     status.ResourceStatus `bson:"status"`
	Provider   string                `bson:"provider" plugin:"providerName"`
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

func NewVPC(id identifier.VPC, provider types.ProviderType) VPC {
	return VPC{
		Metadata: metadata.Metadata{
			Name: id.ID,
		},
		Identifier: id,
		Status:     status.New(id.ID, types.VPC, provider),
		Networks:   make(NetworkCollection),
	}
}

func (vpc *VPC) New(id identifier.ID, provider types.ProviderType) IResource {
	if reflect.TypeOf(id) != reflect.TypeOf(identifier.VPC{}) {
		panic(errors.InvalidArgument.WithMessage("invalid id type"))
	}
	res := NewVPC(id.(identifier.VPC), provider)
	return &res
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

func (vpc *VPC) GetPluginReference() resourcePlugin.Reference {
	if !vpc.Status.PluginReference.ChartReference.Empty() {
		return vpc.Status.PluginReference
	}
	switch vpc.Status.PluginReference.ResourceReference.ProviderType {
	case types.GCP:
		vpc.Status.PluginReference.ChartReference = resourcePlugin.HelmChartReference{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VPC.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.VPC.Version,
		}
		return vpc.Status.PluginReference
	}
	panic(errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", vpc.Status.PluginReference.ResourceReference.ProviderType)))
}

func (vpc *VPC) FromMap(data map[string]interface{}) {
	err := resourcePlugin.InjectMapIntoStruct(data, vpc)
	if !err.IsOk() {
		panic(err)
	}
}

func (vpc *VPC) Insert(project Project, update ...bool) {
	shouldUpdate := false
	if len(update) > 0 {
		shouldUpdate = update[0]
	}
	_, ok := project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID]
	if !ok && shouldUpdate {
		errors.NotFound.WithMessage(fmt.Sprintf("provider %s not found", vpc.Identifier.ProviderID))
	}
	if ok && !shouldUpdate {
		panic(errors.Conflict.WithMessage(fmt.Sprintf("vpc %s in provider %s already exists", vpc.Identifier.ID, vpc.Identifier.ProviderID)))
	}
	project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID] = *vpc
}

func (vpc *VPC) Remove(project Project) {
	_, ok := project.Resources[vpc.Identifier.ProviderID].VPCs[vpc.Identifier.ID]
	if !ok {
		panic(errors.NotFound.WithMessage(fmt.Sprintf("vpc %s in provider %s not found", vpc.Identifier.ID, vpc.Identifier.ProviderID)))
	}
	delete(project.Resources[vpc.Identifier.ProviderID].VPCs, vpc.Identifier.ID)
}

func (vpc *VPC) Equals(other VPC) bool {
	return vpc.Metadata.Equals(other.Metadata) &&
		vpc.Identifier.Equals(other.Identifier) &&
		vpc.Status.Equals(other.Status) &&
		vpc.Provider == other.Provider &&
		vpc.Networks.Equals(other.Networks)
}
