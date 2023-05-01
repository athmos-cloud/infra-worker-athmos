package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"

type UpdateResourceRequest struct {
	ProjectID        string
	ResourceID       identifier.IdPayload
	NewResourceSpecs map[string]interface{}
}
