package aws

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	resourceRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/resource"
)

type awsRepository struct{}

func (aws *awsRepository) FindSqlDB(context context.Context, option option.Option) (*instance.SqlDB, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) FindAllSqlDBs(context context.Context, option option.Option) (*instance.SqlDBCollection, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) FindAllRecursiveSqlDBs(context context.Context, option option.Option, channel *resourceRepo.SqlDBChannel) {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) CreateSqlDB(context context.Context, db *instance.SqlDB) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) UpdateSqlDB(context context.Context, db *instance.SqlDB) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) DeleteSqlDB(context context.Context, db *instance.SqlDB) errors.Error {
	//TODO implement me
	panic("implement me")
}

func (aws *awsRepository) SqlDBExists(context context.Context, db *instance.SqlDB) (bool, errors.Error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository() resourceRepo.Resource {
	return &awsRepository{}
}
