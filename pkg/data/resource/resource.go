package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/plugin"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/metadata"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/status"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

const (
	resourceIDSuffixLength = 10
)

type IResource interface {
	New(NewResourcePayload) IResource
	GetPluginReference() plugin.Reference
	FromMap(map[string]interface{})
	GetIdentifier() identifier.ID
	SetMetadata(metadata metadata.CreateMetadataRequest)
	GetMetadata() metadata.Metadata
	SetStatus(resourceStatus status.ResourceStatus)
	GetStatus() status.ResourceStatus
	Insert(resource IResource, update ...bool)
	Remove(resource IResource)
}

type NewResourcePayload struct {
	Name             string
	ParentIdentifier identifier.ID
	Provider         types.ProviderType
	Monitored        bool
	Tags             map[string]string
}

func (payload NewResourcePayload) Validate() {
	if payload.Name == "" {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("invalid name: %s", payload.Name)))
	}
	if payload.Provider == "" {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("invalid provider: %s", payload.Provider)))
	}
}
