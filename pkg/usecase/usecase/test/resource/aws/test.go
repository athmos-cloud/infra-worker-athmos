package aws

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func initTest(t *testing.T) (context.Context, *testResource.TestResource) {
	ctx := test.NewContext()
	resourceTest := testResource.NewTest(ctx, t)
	ctx.Set(context.ProviderTypeKey, types.ProviderAWS)
	ctx.Set(context.CurrentNamespace, test.TestNamespaceContextKey)

	return ctx, resourceTest
}

func clear(ctx context.Context) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ctx.Value(test.TestNamespaceContextKey).(string),
		},
	}
	_ = kubernetes.Client().Client.Delete(ctx, namespace)
}

func suiteTeardown(ctx context.Context, t *testing.T, container *gnomock.Container) {
	require.NoError(t, gnomock.Stop(container))
	ClearFixtures(ctx)
	clear(ctx)
}
