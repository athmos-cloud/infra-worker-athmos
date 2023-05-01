package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type DeleteResourceRequest struct {
	ProjectID  string               `json:"projectID"`
	ResourceID identifier.IdPayload `json:"resourceID"`
}

type DeleteResourceResponse struct {
	Message errors.Error `json:"message"`
}
