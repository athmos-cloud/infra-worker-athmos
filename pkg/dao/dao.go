package dao

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
)

type ClientType string

const (
	Minio ClientType = "minio"
)

type IDataAccessObject interface {
	Create(context.Context, option.Option) (interface{}, errors.Error)
	Get(context.Context, option.Option) (interface{}, errors.Error)
	Exists(context.Context, option.Option) (bool, errors.Error)
	GetAll(context.Context, option.Option) (interface{}, errors.Error)
	Update(context.Context, option.Option) errors.Error
	Delete(context.Context, option.Option) errors.Error
	Close(context context.Context) errors.Error
}
