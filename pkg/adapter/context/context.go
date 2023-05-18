package context

import "context"

type Context interface {
	context.Context
	JSON(int, any)
	BindJSON(any) error
}
