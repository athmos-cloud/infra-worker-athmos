package presenter

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/config"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output"
	"github.com/gin-gonic/gin"
)

type Secret struct{}

func NewSecretPresenter() output.SecretPort {
	return &Secret{}
}

func (s *Secret) Render(ctx context.Context, secretAuth *secret.Secret) {
	resp := dto.GetSecretResponse{
		ID:          secretAuth.IDField.ID.Hex(),
		Name:        secretAuth.Name,
		Description: secretAuth.Description,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (s *Secret) RenderAll(ctx context.Context, secretAuths *[]secret.Secret) {
	resp := make([]dto.ListSecretResponseItem, len(*secretAuths))
	for i, secretAuth := range *secretAuths {
		resp[i] = dto.ListSecretResponseItem{
			ID:          secretAuth.IDField.ID.Hex(),
			Name:        secretAuth.Name,
			Description: secretAuth.Description,
		}
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (s *Secret) RenderCreate(ctx context.Context, secretAuth *secret.Secret) {
	resp := dto.CreateSecretResponse{
		RedirectionURL: fmt.Sprintf("%s/projects/%s/secrets/%s", config.Current.RedirectionURL, ctx.Value(context.ProjectIDKey), secretAuth.Name),
		Prerequisites:  secretAuth.Prerequisites,
	}
	logger.Info.Printf("Resp: %v", resp)
	ctx.JSON(201, gin.H{"payload": resp})
}

func (s *Secret) RenderUpdate(ctx context.Context, secretAuth *secret.Secret) {
	ctx.JSON(204, gin.H{"message": "Secret " + secretAuth.IDField.ID.Hex() + " updated"})
}

func (s *Secret) RenderDelete(ctx context.Context) {
	ctx.JSON(204, gin.H{"message": "Secret deleted"})
}
