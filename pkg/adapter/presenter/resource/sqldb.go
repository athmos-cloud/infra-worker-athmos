package resource

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
	"github.com/gin-gonic/gin"
)

type sqlDB struct{}

func NewSqlDBPresenter() output.SqlDBPort {
	return &sqlDB{}
}

func (s *sqlDB) Render(ctx context.Context, db *instance.SqlDB) {
	resp := &dto.GetSqlDBResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *db,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (s *sqlDB) RenderCreate(ctx context.Context, db *instance.SqlDB) {
	resp := &dto.CreateSqlDBResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *db,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (s *sqlDB) RenderUpdate(ctx context.Context, db *instance.SqlDB) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("db %s updated", db.IdentifierID.SqlDB)})

}

func (s *sqlDB) RenderDelete(ctx context.Context, db *instance.SqlDB) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("db %s deleted", db.IdentifierID.SqlDB)})

}
