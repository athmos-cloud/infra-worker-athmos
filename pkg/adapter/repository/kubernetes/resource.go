package kubernetes

import (
	"github.com/fatih/structs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type Resource struct {
	Resource  interface{}
	GVK       schema.GroupVersionKind
	Name      string
	Namespace *string
	Labels    *map[string]string
}

func toCamelCase(m interface{}) interface{} {
	switch item := m.(type) {

	case map[string]interface{}:
		newMap := make(map[string]interface{}, len(item))
		for k, v := range item {
			lowerKey := string([]rune(k)[:0]) + strings.ToLower(k[0:1]) + string([]rune(k)[1:])
			newMap[lowerKey] = toCamelCase(v)
		}
		return newMap

	case []interface{}:
		for i, v := range item {
			item[i] = toCamelCase(v)
		}
	}

	return m
}

func toUnstructured(crossplaneResource Resource) *unstructured.Unstructured {
	m := structs.Map(crossplaneResource.Resource)
	u := &unstructured.Unstructured{
		Object: toCamelCase(m).(map[string]interface{}),
	}
	u.SetGroupVersionKind(crossplaneResource.GVK)
	u.SetLabels(*crossplaneResource.Labels)
	u.SetName(crossplaneResource.Name)
	if crossplaneResource.Namespace != nil {
		u.SetNamespace(*crossplaneResource.Namespace)
	}
	return u
}

func gVKtoGVR(gvk schema.GroupVersionKind) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind) + "s",
	}
}
