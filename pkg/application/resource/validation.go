package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"reflect"
)

func validateCreate(request CreateResourceRequest) {
	if request.ProjectID == "" {
		panic(errors.BadRequest.WithMessage("projectID is required"))
	}
	if request.ProviderType == "" {
		panic(errors.BadRequest.WithMessage("providerType is required"))
	}
	if request.ResourceType == "" {
		panic(errors.BadRequest.WithMessage("resourceType is required"))
	}
	if request.ParentIdentifier == nil && request.ResourceType != types.Provider {
		panic(errors.BadRequest.WithMessage("parentIdentifier is required"))
	} else if !identifier.IDParentMatchesWithResource(request.ParentIdentifier, request.ResourceType) {
		panic(errors.BadRequest.WithMessage(
			fmt.Sprintf("Wrong parent id of type %s to create %s resource", reflect.TypeOf(request.ParentIdentifier), request.ResourceType),
		))
	}
	if request.ResourceSpecs == nil {
		panic(errors.BadRequest.WithMessage("resourceSpecs is required"))
	}
}
