package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type VPC struct {
	Metadata            metadata.Metadata       `bson:"metadata"`
	Identifier          identifier.VPC          `bson:"identifier"`
	KubernetesResources kubernetes.ResourceList `bson:"kubernetesResources"`
	Provider            string                  `bson:"provider"`
	Networks            NetworkCollection       `bson:"networks"`
}

type VPCCollection map[string]VPC

func NewVPC(id identifier.VPC) VPC {
	return VPC{
		Metadata: metadata.Metadata{
			Name: id.ID,
		},
		Identifier: id,
		Networks:   make(NetworkCollection),
	}
}

func (vpc *VPC) WithMetadata(request metadata.CreateMetadataRequest) {
	vpc.Metadata = metadata.New(request)
}

func (vpc *VPC) GetMetadata() metadata.Metadata {
	return vpc.Metadata
}

func (vpc *VPC) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VPC.Chart,
			ChartVersion: config.Current.Plugins.Crossplane.GCP.VPC.Version,
		}, errors.Error{}
	}
	return dto.GetPluginReferenceResponse{}, errors.InvalidArgument.WithMessage(fmt.Sprintf("provider type %s not supported", request.ProviderType))
}

func (vpc *VPC) FromMap(data map[string]interface{}) errors.Error {
	if data["name"] == nil {
		return errors.InvalidArgument.WithMessage("name is required")
	}
	return errors.OK
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

func (vpc *VPC) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
