package resource

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	errorCtrl "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/error"
	resourceValidator "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/validator/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	resource2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
)

type SqlDB interface {
	GetSqlDB(context.Context)
	CreateSqlDB(context.Context)
	UpdateSqlDB(context.Context)
	DeleteSqlDB(context.Context)
}

type sqlDBController struct {
	dbUseCase resource2.SqlDB
	dbOutput  output.SqlDBPort
}

func NewSqlDBController(dbUseCase resource2.SqlDB, dbOutput output.SqlDBPort) SqlDB {
	return &sqlDBController{dbUseCase: dbUseCase, dbOutput: dbOutput}
}

func (dc *sqlDBController) GetSqlDB(ctx context.Context) {
	if err := resourceValidator.GetSqlDB(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	net := &instance.SqlDB{}
	if err := dc.dbUseCase.Get(ctx, net); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		dc.dbOutput.Render(ctx, net)
	}
}

func (dc *sqlDBController) CreateSqlDB(ctx context.Context) {
	if err := resourceValidator.CreateSqlDB(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	net := &instance.SqlDB{}
	if err := dc.dbUseCase.Create(ctx, net); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		dc.dbOutput.RenderCreate(ctx, net)
	}
}

func (dc *sqlDBController) UpdateSqlDB(ctx context.Context) {
	if err := resourceValidator.UpdateSqlDB(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	}
	network := &instance.SqlDB{}
	if err := dc.dbUseCase.Update(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	} else {
		dc.dbOutput.RenderUpdate(ctx, network)
	}
}

func (dc *sqlDBController) DeleteSqlDB(ctx context.Context) {
	if err := resourceValidator.DeleteSqlDB(ctx); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
		return
	}
	network := &instance.SqlDB{}
	if err := dc.dbUseCase.Delete(ctx, network); !err.IsOk() {
		errorCtrl.RaiseError(ctx, err)
	} else {
		dc.dbOutput.Render(ctx, network)
	}
}
