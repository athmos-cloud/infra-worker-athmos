package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
)

type SqlDBChannel struct {
	Channel      chan *instance.SqlDBCollection
	ErrorChannel chan errors.Error
}

type SqlDB interface {
	FindSqlDB(context.Context, option.Option) (*instance.SqlDB, errors.Error)
	FindAllSqlDBs(context.Context, option.Option) (*instance.SqlDBCollection, errors.Error)
	FindAllRecursiveSqlDBs(context.Context, option.Option, *SqlDBChannel)
	CreateSqlDB(context.Context, *instance.SqlDB) errors.Error
	UpdateSqlDB(context.Context, *instance.SqlDB) errors.Error
	DeleteSqlDB(context.Context, *instance.SqlDB) errors.Error
	SqlDBExists(context.Context, *instance.SqlDB) (bool, errors.Error)
}
