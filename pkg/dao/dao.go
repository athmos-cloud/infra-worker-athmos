package dao

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type IDataAccessObject interface {
	Create(context.Context, option.Option) interface{}
	Get(context.Context, option.Option) interface{}
	Exists(context.Context, option.Option) bool
	GetAll(context.Context, option.Option) interface{}
	Update(context.Context, option.Option)
	Delete(context.Context, option.Option)
	Close(context context.Context)
}
