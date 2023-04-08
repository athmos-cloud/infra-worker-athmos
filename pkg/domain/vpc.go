package domain

import (
	"fmt"
	"github.com/PaulBarrie/infra-worker/pkg/common"
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/config"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/utils"
)

type VPC struct {
	Metadata Metadata  `bson:"metadata"`
	Networks []Network `bson:"networks"`
}

func (vpc *VPC) WithMetadata(request CreateMetadataRequest) {
	vpc.Metadata = New(request)
}

func (vpc *VPC) GetMetadata() Metadata {
	return vpc.Metadata
}

func (vpc *VPC) GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error) {
	switch request.ProviderType {
	case common.GCP:
		return dto.GetPluginReferenceResponse{
			ChartName:    config.Current.Plugins.Crossplane.GCP.VPC.ChartName,
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

func (vpc *VPC) InsertIntoProject(project Project, upsert bool) errors.Error {
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
