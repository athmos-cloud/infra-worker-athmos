package crossplane

import "fmt"

const (
	ExternalNameAnnotation = "crossplane.io/external-name"
	pauseAnnotation        = "crossplane.io/paused"
)

func GetAnnotations(managed bool, name string) map[string]string {
	return map[string]string{
		ExternalNameAnnotation: name,
		pauseAnnotation:        fmt.Sprintf("%t", !managed),
	}
}
