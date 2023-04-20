package repository

import (
	"context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type IRepository interface {
	Create(context.Context, option.Option) (interface{}, errors.Error)
	Get(context.Context, option.Option) (interface{}, errors.Error)
	Watch(context.Context, option.Option) (interface{}, errors.Error)
	List(context.Context, option.Option) (interface{}, errors.Error)
	Update(context.Context, option.Option) errors.Error
	Delete(context.Context, option.Option) errors.Error
}
