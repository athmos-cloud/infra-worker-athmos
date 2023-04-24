package domain

import (
	dto "github.com/athmos-cloud/infra-worker-athmos/pkg/common/dto/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/project"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type IResource interface {
	GetPluginReference(request dto.GetPluginReferenceRequest) (dto.GetPluginReferenceResponse, errors.Error)
	FromMap(map[string]interface{}) errors.Error
	WithMetadata(metadata CreateMetadataRequest)
	GetMetadata() Metadata
	InsertIntoProject(project domain.Project, upsert bool) errors.Error
}
