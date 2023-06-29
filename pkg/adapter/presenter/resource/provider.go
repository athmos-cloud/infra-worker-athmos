package resource

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	output "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/output/resource"
)

type provider struct{}

func NewProviderPresenter() output.ProviderPort {
	return &provider{}
}

func (p *provider) Render(ctx context.Context, provider *resource.Provider) {
	resp := dto.GetProviderResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *provider,
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (p *provider) RenderCreate(ctx context.Context, provider *resource.Provider) {
	resp := dto.CreateProviderResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *provider,
	}
	ctx.JSON(201, gin.H{"payload": resp})
}

func (p *provider) RenderUpdate(ctx context.Context, provider *resource.Provider) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("provider %s updated", provider.IdentifierID.Provider)})
}

func (p *provider) RenderAll(ctx context.Context, providers *resource.ProviderCollection) {
	resp := dto.ListProvidersResponse{
		ProjectID: ctx.Value(context.ProjectIDKey).(string),
		Payload:   *providers,
	}
	if len(*providers) == 0 {
		ctx.JSON(204, []string{})
		return
	}
	ctx.JSON(200, gin.H{"payload": resp})
}

func (p *provider) RenderDelete(ctx context.Context, provider *resource.Provider) {
	ctx.JSON(204, gin.H{"message": fmt.Sprintf("provider %s deleted", provider.IdentifierID.Provider)})
}
