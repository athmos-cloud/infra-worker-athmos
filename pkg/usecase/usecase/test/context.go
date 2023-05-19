package test

import "context"

type Context struct {
	context.Context
}

func NewContext() Context {
	return Context{context.Background()}
}

func (c Context) WithValue(key, val any) Context {
	c.Context = context.WithValue(c.Context, key, val)
	return c
}

func (c Context) JSON(i int, a any) {
	//TODO implement me
	panic("implement me")
}

func (c Context) BindJSON(a any) error {
	//TODO implement me
	panic("implement me")
}
