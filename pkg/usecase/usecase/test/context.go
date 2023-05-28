package test

import (
	goContext "context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
)

const (
	TestNamespaceContextKey = "testNamespace"
	TestSecretContextKey    = "testSecret"
)

type Context struct {
	goContext.Context
}

func NewContext() context.Context {
	return &Context{goContext.Background()}
}

func (c *Context) Set(key string, val any) {
	c.Context = goContext.WithValue(c.Context, key, val)
}

func (c *Context) JSON(i int, a any) {
	//TODO implement me
	panic("implement me")
}

func ClearContext(ctx *Context) {
	ctx.Set(context.RequestKey, nil)
}
