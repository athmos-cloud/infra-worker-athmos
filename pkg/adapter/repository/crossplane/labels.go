package crossplane

const (
	managedByLabel = "app.kubernetes.io/managed-by"
	managedByValue = "athmos"
	projectIDLabel = "athmos.cloud/project-id"
)

func GetBaseLabels(projectID string) map[string]string {
	return map[string]string{
		managedByLabel: managedByValue,
		projectIDLabel: projectID,
	}
}
