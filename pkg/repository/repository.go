package repository

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type IRepository interface {
	Create(context.Context, option.Option) interface{}
	Get(context.Context, option.Option) interface{}
	Watch(context.Context, option.Option) interface{}
	List(context.Context, option.Option) interface{}
	Update(context.Context, option.Option) interface{}
	Delete(context.Context, option.Option)
}
