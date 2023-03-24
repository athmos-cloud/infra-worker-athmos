package repository

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/errors"
	"github.com/PaulBarrie/infra-worker/pkg/kernel/option"
)

type ClientType string

const (
	Minio ClientType = "minio"
)

type IRepository interface {
	Create(context.Context, option.Option) (interface{}, errors.Error)
	Get(context.Context, option.Option) (interface{}, errors.Error)
	GetAll(context.Context, option.Option) (interface{}, errors.Error)
	Update(context.Context, option.Option) errors.Error
	Delete(context.Context, option.Option) errors.Error
	Close(context context.Context) errors.Error
}
