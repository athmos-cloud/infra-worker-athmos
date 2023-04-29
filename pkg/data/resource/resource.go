package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/common"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type IResource interface {
	New(identifier.ID, common.ProviderType) (IResource, errors.Error)
	GetPluginReference() (plugin.Reference, errors.Error)
	FromMap(map[string]interface{}) errors.Error
	SetMetadata(metadata metadata.CreateMetadataRequest)
	GetMetadata() metadata.Metadata
	SetStatus(resourceStatus status.ResourceStatus)
	GetStatus() status.ResourceStatus
	Insert(project Project, update ...bool) errors.Error
	Remove(project Project) errors.Error
}
