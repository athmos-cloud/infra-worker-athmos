package resource

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type IResource interface {
	New(identifier.ID) (IResource, errors.Error)
	GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error)
	FromMap(map[string]interface{}) errors.Error
	WithMetadata(metadata metadata.CreateMetadataRequest)
	GetMetadata() metadata.Metadata
	Insert(project Project, update ...bool) errors.Error
	Remove(project Project) errors.Error
}
