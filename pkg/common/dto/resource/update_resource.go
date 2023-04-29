package resource

import "github.com/athmos-cloud/infra-worker-athmos/pkg/data/resource/identifier"

type UpdateResourceRequest struct {
	ProjectID        string
	ResourceID       identifier.ID
	NewResourceSpecs map[string]interface{}
}
