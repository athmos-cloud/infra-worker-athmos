package resource

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type IResource interface {
	GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error)
	FromMap(map[string]interface{}) errors.Error
	WithMetadata(metadata metadata.CreateMetadataRequest)
	GetMetadata() metadata.Metadata
	ToDomain() (interface{}, errors.Error)
	Insert(project Project, update ...bool) errors.Error
}
