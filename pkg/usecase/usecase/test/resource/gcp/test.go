package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func initTest(t *testing.T) (context.Context, *testResource.TestResource, usecase.Provider) {
	ctx := test.NewContext()
	resourceTest := testResource.NewTest(ctx, t)
	gcpRepo := gcp.NewRepository()
	uc := usecase.NewProviderUseCase(resourceTest.ProjectRepo, resourceTest.SecretRepo, gcpRepo, nil, nil)
	ctx.Set(context.ProviderTypeKey, types.ProviderGCP)

	return ctx, resourceTest, uc
}

func clear(ctx context.Context) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ctx.Value(test.TestNamespaceContextKey).(string),
		},
	}
	_ = kubernetes.Client().Client.Delete(ctx, namespace)
}
