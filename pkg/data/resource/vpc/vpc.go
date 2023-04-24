package resources

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/kubernetes"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type VPC struct {
	Metadata            domain.Metadata              `bson:"metadata"`
	KubernetesResources kubernetes.ResourceList      `bson:"kubernetesResources"`
	Identifier          VPCIdentifier                `bson:"hierarchyLocation"`
	Provider            string                       `bson:"provider"`
	Networks            map[string]resources.Network `bson:"networks"`
}

type VPCIdentifier struct {
	ID         string
	ProviderID string
	SubnetID   string
}

func (vpc *VPC) WithMetadata(request domain.CreateMetadataRequest) {
	vpc.Metadata = domain.New(request)
}

func (vpc *VPC) GetMetadata() domain.Metadata {
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

func (vpc *VPC) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	for _, r := range project.Resources {
		for _, v := range r.VPCs {
			if v.Identifier.ID == vpc.Identifier.ID && upsert {
				v = *vpc
				return errors.OK
			} else if v.Identifier.ID == vpc.Identifier.ID && !upsert {
				return errors.AlreadyExists.WithMessage("vpc already exists")
			}
		}
	}
	return errors.OK
}

func (vpc *VPC) ToDomain() (interface{}, errors.Error) {
	//TODO implement me
	panic("implement me")
}
