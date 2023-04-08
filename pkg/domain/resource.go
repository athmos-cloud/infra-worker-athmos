package domain

import (
	dto "github.com/PaulBarrie/infra-worker/pkg/common/dto/resource"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
)

type IResource interface {
	GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error)
	FromMap(map[string]interface{}) errors.Error
	WithMetadata(metadata CreateMetadataRequest)
	GetMetadata() Metadata
	InsertIntoProject(project Project, upsert bool) errors.Error
}
