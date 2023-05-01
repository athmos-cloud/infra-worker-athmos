package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
)

type IResource interface {
	New(identifier.ID, types.ProviderType) IResource
	GetPluginReference() plugin.Reference
	FromMap(map[string]interface{})
	SetMetadata(metadata metadata.CreateMetadataRequest)
	GetMetadata() metadata.Metadata
	SetStatus(resourceStatus status.ResourceStatus)
	GetStatus() status.ResourceStatus
	Insert(project Project, update ...bool)
	Remove(project Project)
}
