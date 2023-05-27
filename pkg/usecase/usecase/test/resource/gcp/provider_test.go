package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	kubernetesRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

var (
	providerGVK = schema.GroupVersionKind{
		Kind:    "ProviderConfig",
		Group:   "gcp.upbound.io",
		Version: "v1beta1",
	}
)

func Test_providerUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	testRes, ctx, uc := initTest(t)
	ctx.Set(context.ResourceTypeKey, types.ProviderResource)
	t.Run("Create a valid provider", func(t *testing.T) {
		req := dto.CreateProviderRequest{
			Name:           "test",
			VPC:            testResource.SecretTestName,
			SecretAuthName: testResource.SecretTestName,
		}
		ctx.Set(context.RequestKey, req)
		provider := &resource.Provider{}
		err := uc.Create(ctx, provider)
		assert.True(t, err.IsOk())

		getRes, errK := testRes.KubernetesRepo.Get(ctx, kubernetesRepo.Resource{
			GVK:  providerGVK,
			Name: provider.IdentifierID.Provider,
		})
		assert.True(t, errK.IsOk())

		kubeResource := &v1beta1.ProviderConfig{}
		errConvert := kubernetes.Client().Client.Scheme().Convert(getRes, kubeResource, ctx)
		assert.NoError(t, errConvert)

		assert.Equal(t, provider.IdentifierID.Provider, kubeResource.Name)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"name.provider":                "test",
			"name.secret":                  testResource.SecretTestName,
		}
		assert.Equal(t, wantLabels, kubeResource.Labels)
		usedSecret := ctx.Value(test.TestSecretContextKey).(secret.Secret)
		wantSpecs := v1beta1.ProviderConfigSpec{
			ProjectID: req.VPC,
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: usedSecret.Kubernetes.SecretKey,
						SecretReference: xpv1.SecretReference{
							Namespace: usedSecret.Kubernetes.Namespace,
							Name:      usedSecret.Kubernetes.SecretName,
						},
					},
				},
			},
		}
		assert.Equal(t, wantSpecs, kubeResource.Spec)
	})
	t.Run("Create a provider with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a provider with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Delete a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a provider with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Get a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_GetRecursively(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("GetRecursively a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_List(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("List providers should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("List providers in a non-existing project should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_providerUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
	}()
	t.Run("Update a valid provider should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing provider should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}
