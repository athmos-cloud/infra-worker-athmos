package crossplane

import "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"

const (
	managedByLabel = "app.kubernetes.io/managed-by"
	managedByValue = "athmos"
)

func GetBaseLabels(projectID string) map[string]string {
	return map[string]string{
		managedByLabel:          managedByValue,
		model.ProjectIDLabelKey: projectID,
	}
}
