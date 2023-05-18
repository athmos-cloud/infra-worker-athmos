package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

const (
	resourceIDSuffixLength = 10
)

type NewResourcePayload struct {
	Name             string
	ParentIdentifier identifier.ID
	//Provider         types.ProviderType
	Managed bool
	Tags    map[string]string
}

func (payload NewResourcePayload) Validate() {
	if payload.Name == "" {
		panic(errors.InternalError.WithMessage(fmt.Sprintf("invalid name: %s", payload.Name)))
	}
}
