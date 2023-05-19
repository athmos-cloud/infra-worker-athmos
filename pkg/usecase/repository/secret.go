package repository

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type Secret interface {
	Find(context.Context, option.Option) (*secret.Secret, errors.Error)
	FindAll(context.Context, option.Option) (*[]secret.Secret, errors.Error)
	Create(context.Context, *secret.Secret) errors.Error
	Update(context.Context, *secret.Secret) errors.Error
	Delete(context.Context, *secret.Secret) errors.Error
}

type KubernetesSecret interface {
	Create(context.Context, option.Option) (*secret.Kubernetes, errors.Error)
	Update(context.Context, option.Option) errors.Error
	Delete(context.Context, option.Option) errors.Error
}
