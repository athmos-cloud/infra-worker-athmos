package resources

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	domain2 "github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
)

type VPC struct {
	Metadata domain.Metadata              `bson:"metadata"`
	Provider string                       `bson:"provider"`
	Networks map[string]resources.Network `bson:"networks"`
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
	*vpc = VPC{}
	if data["id"] == nil {
		vpc.Metadata.ID = utils.GenerateUUID()
	} else {
		vpc.Metadata.ID = data["id"].(string)
	}
	if data["name"] == nil {
		return errors.InvalidArgument.WithMessage("name is required")
	}
	return errors.OK
}

func (vpc *VPC) InsertIntoProject(project domain2.Project, upsert bool) errors.Error {
	for _, r := range project.Resources {
		for _, v := range r.VPCs {
			if v.Metadata.ID == vpc.Metadata.ID && upsert {
				v = *vpc
				return errors.OK
			} else if v.Metadata.ID == vpc.Metadata.ID && !upsert {
				return errors.AlreadyExists.WithMessage("vpc already exists")
			}
		}
	}
	return errors.OK
}
