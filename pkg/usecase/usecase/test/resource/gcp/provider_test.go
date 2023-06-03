package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	repository2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	secretRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	secretModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	domainTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

var (
	providerGVK = schema.GroupVersionKind{
		Kind:    "ProviderConfig",
		Group:   "gcp.upbound.io",
		Version: "v1beta1",
	}
)

type wantProvider struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.ProviderConfigSpec
}

func clearProvider(ctx context.Context) {
	clear(ctx)
	providers := &v1beta1.ProviderConfigList{}

	err := kubernetes.Client().Client.List(ctx, providers)
	if err != nil {
		return
	}
	for _, provider := range providers.Items {
		err = kubernetes.Client().Client.Delete(ctx, &provider)
		if err != nil {
			logger.Warning.Printf("Error deleting provider %s: %v", provider.Name, err)
			continue
		}
	}
}
func Test_providerUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, uc := initTest(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearProvider(ctx)
	}()

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)
	t.Run("Create a valid provider", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)
		kubeResource := &v1beta1.ProviderConfig{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: provider.IdentifierID.Provider}, kubeResource)
		assert.NoError(t, errk)

		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          provider.IdentifierID.Provider,
			"identifier.vpc":               provider.IdentifierID.VPC,
			"name.vpc":                     provider.IdentifierID.VPC,
			"name.provider":                "test",
			"name.secret":                  testResource.SecretTestName,
		}
		usedSecret := ctx.Value(test.TestSecretContextKey).(secret.Secret)
		wantSpecs := v1beta1.ProviderConfigSpec{
			ProjectID: provider.IdentifierName.VPC,
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
		want := wantProvider{
			Name:   provider.IdentifierID.Provider,
			Labels: wantLabels,
			Spec:   wantSpecs,
		}
		got := wantProvider{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertProviderEqual(t, want, got)
	})
	t.Run("Create a provider with a non-existing secret should fail", func(t *testing.T) {
		req := dto.CreateProviderRequest{
			Name:           "test",
			VPC:            testResource.SecretTestName,
			SecretAuthName: "this-secret-does-not-exist",
		}
		ctx.Set(context.RequestKey, req)
		provider := &resource.Provider{}
		err := uc.Create(ctx, provider)

		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
	t.Run("Create a provider with an already existing name should fail", func(t *testing.T) {
		_ = ProviderFixture(ctx, t, uc)
		newProvider := &resource.Provider{}
		err := uc.Create(ctx, newProvider)
		assert.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_providerUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, testRes, uc := initTest(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearProvider(ctx)
	}()

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)

	t.Run("Delete a valid provider should succeed", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)

		delReq := dto.DeleteProviderRequest{
			IdentifierID: identifier.Provider{
				Provider: provider.IdentifierID.Provider,
			},
		}
		ctx.Set(context.RequestKey, delReq)
		delProvider := &resource.Provider{}
		err := uc.Delete(ctx, delProvider)

		assert.True(t, err.IsOk())
		assert.Equal(t, errors.NoContent.Code, err.Code)
	})

	t.Run("Delete a non-existing provider should fail", func(t *testing.T) {
		provider := &resource.Provider{}
		delReq := dto.DeleteProviderRequest{
			IdentifierID: identifier.Provider{
				Provider: "this-provider-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, delReq)
		err := uc.Delete(ctx, provider)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})

	t.Run("Delete cascade a provider should succeed", func(t *testing.T) {
		gcpRepo := gcp.NewRepository()
		sshRepo := repository2.NewSSHKeyRepository()

		nuc := usecase.NewNetworkUseCase(testRes.ProjectRepo, gcpRepo, nil, nil)
		suc := usecase.NewSubnetworkUseCase(testRes.ProjectRepo, gcpRepo, nil, nil)
		vuc := usecase.NewVMUseCase(testRes.ProjectRepo, sshRepo, gcpRepo, nil, nil)

		provider := ProviderFixture(ctx, t, uc)
		NetworkFixture(ctx, t, nuc)
		SubnetworkFixture(ctx, t, suc)
		VMFixture(ctx, t, vuc)

		delReq := dto.DeleteProviderRequest{
			IdentifierID: identifier.Provider{
				Provider: provider.IdentifierID.Provider,
			},
		}
		ctx.Set(context.RequestKey, delReq)
		delProvider := &resource.Provider{}
		err := uc.Delete(ctx, delProvider)

		assert.Equal(t, errors.NoContent.Code, err.Code)
	})
}

func Test_providerUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, uc := initTest(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearProvider(ctx)
	}()

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)
	t.Run("Get a valid provider should succeed", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)

		getReq := dto.GetProviderRequest{
			IdentifierID: identifier.Provider{
				Provider: provider.IdentifierID.Provider,
			},
		}
		ctx.Set(context.RequestKey, getReq)
		getProvider := &resource.Provider{}
		err := uc.Get(ctx, getProvider)
		assert.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, provider.IdentifierName, getProvider.IdentifierName)
		assert.Equal(t, provider.IdentifierID, getProvider.IdentifierID)
		assert.Equal(t, provider.Auth, getProvider.Auth)
	})
	t.Run("Get a non-existing provider should fail", func(t *testing.T) {
		getReq := dto.GetProviderRequest{
			IdentifierID: identifier.Provider{
				Provider: "this-provider-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, getReq)
		getProvider := &resource.Provider{}
		err := uc.Get(ctx, getProvider)
		assert.Equal(t, errors.NotFound.Code, err.Code)
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
	ctx, _, uc := initTest(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearProvider(ctx)
	}()

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)

	t.Run("List providers should succeed", func(t *testing.T) {
		create := func(name string) {
			req := dto.CreateProviderRequest{
				Name:           name,
				VPC:            testResource.SecretTestName,
				SecretAuthName: testResource.SecretTestName,
			}
			ctx.Set(context.RequestKey, req)
			provider := &resource.Provider{}
			err := uc.Create(ctx, provider)
			assert.True(t, err.IsOk())
		}
		create("test1")
		create("test2")
		providerList := resource.ProviderCollection{}
		err := uc.List(ctx, &providerList)
		assert.True(t, err.IsOk())
		assert.Equal(t, 2, len(providerList))
		for _, provider := range providerList {
			assert.True(t, provider.IdentifierName.Provider == "test1" || provider.IdentifierName.Provider == "test2")
		}
	})

	t.Run("List providers in a non-existing project should fail", func(t *testing.T) {
		ctx.Set(context.ProjectIDKey, "this-project-does-not-exist")
		providerList := resource.ProviderCollection{}
		err := uc.List(ctx, &providerList)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_providerUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, uc := initTest(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearProvider(ctx)
	}()

	t.Run("Update a valid provider should succeed", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)

		// New secret
		secrRepo := secretRepo.NewSecretRepository()
		kubeSecretRepo := secretRepo.NewKubernetesRepository()
		createdSecret, err := kubeSecretRepo.Create(ctx, option.Option{
			Value: repository.CreateKubernetesSecretRequest{
				ProjectID:   ctx.Value(context.ProjectIDKey).(string),
				SecretName:  "test-2",
				SecretKey:   "key.json",
				SecretValue: []byte("{\"test\":\"test\"}"),
			},
		})
		require.True(t, err.IsOk())
		secretAuth := secretModel.NewSecret("test-2", "A new secret", *createdSecret)
		err = secrRepo.Create(ctx, secretAuth)
		require.True(t, err.IsOk())

		updateReq := dto.UpdateProviderRequest{
			IdentifierID:   provider.IdentifierID,
			Name:           "test2",
			SecretAuthName: "test-2",
		}
		ctx.Set(context.RequestKey, updateReq)
		updatedProvider := &resource.Provider{}
		err = uc.Update(ctx, updatedProvider)
		assert.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.GetProviderRequest{
			IdentifierID: provider.IdentifierID,
		})
		err = uc.Get(ctx, updatedProvider)
		kubeResource := &v1beta1.ProviderConfig{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: provider.IdentifierID.Provider}, kubeResource)
		assert.NoError(t, errk)

		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          provider.IdentifierID.Provider,
			"identifier.vpc":               provider.IdentifierID.VPC,
			"name.provider":                "test2",
			"name.vpc":                     "test",
			"name.secret":                  "test-2",
		}
		wantSpecs := v1beta1.ProviderConfigSpec{
			ProjectID: provider.IdentifierID.VPC,
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						Key: "key.json",
						SecretReference: xpv1.SecretReference{
							Namespace: secretAuth.Kubernetes.Namespace,
							Name:      secretAuth.Kubernetes.SecretName,
						},
					},
				},
			},
		}
		want := wantProvider{
			Name:   provider.IdentifierID.Provider,
			Labels: wantLabels,
			Spec:   wantSpecs,
		}
		got := wantProvider{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertProviderEqual(t, want, got)

	})
	t.Run("Update a non-existing provider should fail", func(t *testing.T) {
		updateReq := dto.UpdateProviderRequest{
			IdentifierID:   identifier.Provider{Provider: "this-provider-does-not-exist"},
			Name:           "test2",
			SecretAuthName: "test-2",
		}
		ctx.Set(context.RequestKey, updateReq)
		updatedProvider := &resource.Provider{}
		err := uc.Update(ctx, updatedProvider)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
	t.Run("Update a provider with a non-existing secret should return bad request", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)

		updateReq := dto.UpdateProviderRequest{
			IdentifierID:   provider.IdentifierID,
			Name:           "test2",
			SecretAuthName: "this-secret-does-not-exist",
		}
		ctx.Set(context.RequestKey, updateReq)
		updatedProvider := &resource.Provider{}
		err := uc.Update(ctx, updatedProvider)
		assert.Equal(t, errors.BadRequest.Code, err.Code)
	})
}

func assertProviderEqual(t *testing.T, want wantProvider, got wantProvider) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec, got.Spec)
}
