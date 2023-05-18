package repository

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Project interface {
	Find(context.Context, option.Option) *model.Project
	FindAll(context.Context, option.Option) []*model.Project
	Create(context.Context, chan *model.Project, chan errors.Error)
	Update(context.Context, *model.Project) *model.Project
	Delete(context.Context, *model.Project)
}
