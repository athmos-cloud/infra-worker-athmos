package context

import "context"

const (
	RequestKey      = "request"
	ResponseKey     = "response"
	ResponseCodeKey = "response_code"
	ProjectIDKey    = "project_id"
	OwnerIDKey      = "owner_id"
	ProviderTypeKey = "provider_type"
	ResourceTypeKey = "resource_type"
)

type Context interface {
	context.Context
	JSON(int, any)
	Set(string, any)
}
