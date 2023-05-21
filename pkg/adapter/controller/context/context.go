package context

import "context"

const (
	RequestKey      = "request"
	ErrorKey        = "error"
	ProjectIDKey    = "project_id"
	OwnerIDKey      = "owner_id"
	ProviderTypeKey = "provider_type"
	ResourceTypeKey = "resource_type"
)

type Context interface {
	context.Context
	JSON(int, any)
	BindJSON(any) error
}
