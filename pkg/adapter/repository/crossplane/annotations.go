package crossplane

import "fmt"

const (
	ExternalNameAnnotationKey = "crossplane.io/external-name"
	pauseAnnotationKey        = "crossplane.io/paused"
)

func GetAnnotations(managed bool, name string) map[string]string {
	return map[string]string{
		ExternalNameAnnotationKey: name,
		pauseAnnotationKey:        fmt.Sprintf("%t", !managed),
	}
}
