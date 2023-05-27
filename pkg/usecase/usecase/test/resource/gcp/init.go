package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"testing"
)

func initTest(t *testing.T) (*testResource.TestResource, context.Context, usecase.Provider) {
	ctx := test.NewContext()
	resourceTest := testResource.NewTest(ctx, t)
	gcpRepo := gcp.NewRepository()
	uc := usecase.NewProviderUseCase(resourceTest.ProjectRepo, resourceTest.SecretRepo, gcpRepo, nil, nil)
	ctx.Set(context.ProviderTypeKey, types.ProviderGCP)

	return resourceTest, ctx, uc
}
