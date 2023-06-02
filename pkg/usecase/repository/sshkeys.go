package repository

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
)

type SSHKeys interface {
	Create(context.Context, *model.SSHKey) errors.Error
	CreateList(context.Context, model.SSHKeyList) errors.Error
	Get(context.Context, *model.SSHKey) errors.Error
	GetList(context.Context, model.SSHKeyList) errors.Error
	Delete(context.Context, *model.SSHKey) errors.Error
}
