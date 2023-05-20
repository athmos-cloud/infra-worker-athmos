package repository

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Project interface {
	Find(context.Context, option.Option) (*model.Project, errors.Error)
	FindAll(context.Context, option.Option) (*[]model.Project, errors.Error)
	Create(context.Context, *model.Project) errors.Error
	Update(context.Context, *model.Project) errors.Error
	Delete(context.Context, *model.Project) errors.Error
}

type FindProjectByIDRequest struct {
	ID string
}

type FindAllProjectByOwnerRequest struct {
	Owner string
}
