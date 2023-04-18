package kubernetes

import "k8s.io/apimachinery/pkg/runtime/schema"

type GetResourceRequest struct {
	ResourceID schema.GroupVersionResource
	Namespace  string
	Name       string
}

type GetListResourceRequest struct {
	ResourceID schema.GroupVersionResource
	Namespace  string
	Labels     map[string]string
}