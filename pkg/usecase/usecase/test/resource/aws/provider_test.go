package aws

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	repository2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	secretRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	resourceInstance "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	resourceNetwork "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	secretModel "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/secret"
	domainTypes "github.com/athmos-cloud/infra-worker-athmos/pkg/domain/types"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/option"
	secret2 "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/repository/secret"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-aws/apis/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

var (
	providerGVK = schema.GroupVersionKind{
		Kind:    "ProviderConfig",
		Group:   "aws.upbound.io",
		Version: "v1beta1",
	}
)

type wantProvider struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.ProviderConfigSpec
}

/*
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
*/

func Test_providerUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearFixtures(ctx)
	}()

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)
	t.Run("Create a valid provider", func(t *testing.T) {
		defer func() {
			ClearProviderFixtures(ctx)
			clear(ctx)
		}()

		provider := ProviderFixture(ctx, t, uc)

		ctx.Set(testResource.ProviderIDKey, provider.IdentifierID)

		kubeResource := &v1beta1.ProviderConfig{}
		errK := kubernetes.Client().Client.Get(
			ctx,
			types.NamespacedName{Name: provider.IdentifierID.Provider},
			kubeResource)
		assert.NoError(t, errK)

		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          provider.IdentifierID.Provider,
			"identifier.vpc":               provider.IdentifierID.VPC,
			"name.vpc":                     provider.IdentifierID.VPC,
			"name.provider":                "fixture-provider",
			"name.secret":                  testResource.SecretTestName,
		}
		usedSecret := ctx.Value(test.TestSecretContextKey).(secret.Secret)
		wantSpecs := v1beta1.ProviderConfigSpec{
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
		defer func() {
			ClearProviderFixtures(ctx)
			clear(ctx)
		}()

		req := dto.CreateProviderRequest{
			Name:           "test",
			VPC:            testResource.SecretTestName,
			SecretAuthName: "this-secret-does-not-exist",
		}
		ctx.Set(context.RequestKey, req)
		provider := &resource.Provider{}
		err := uc.Create(ctx, provider)

		ctx.Set(testResource.ProviderIDKey, provider.IdentifierID)

		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
	t.Run("Create a provider with an already existing name should fail", func(t *testing.T) {
		defer func() {
			ClearProviderFixtures(ctx)
			clear(ctx)
		}()

		_ = ProviderFixture(ctx, t, uc)
		newProvider := &resource.Provider{}
		err := uc.Create(ctx, newProvider)
		assert.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_providerUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearProviderFixtures(ctx)
		clear(ctx)
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

		nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, gcpRepo, nil, nil)
		suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, gcpRepo, nil, nil)
		vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, gcpRepo, nil, nil)

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
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearProviderFixtures(ctx)
		clear(ctx)
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

func Test_providerUseCase_GetProviderStack(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearProviderFixtures(ctx)
		clear(ctx)
	}()

	t.Run("Get a valid Provider Stack should succeed", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)
		nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
		suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
		fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
		vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, repository2.NewSSHKeyRepository(), nil, aws.NewRepository(), nil)

		ctx.Set(context.RequestKey, dto.CreateNetworkRequest{
			ParentIDProvider: &provider.IdentifierID,
			Name:             "test-network-1",
		})
		net1 := &resourceNetwork.Network{}
		err := nuc.Create(ctx, net1)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateNetworkRequest{
			ParentIDProvider: &provider.IdentifierID,
			Name:             "test-network-2",
		})
		net2 := &resourceNetwork.Network{}
		err = nuc.Create(ctx, net2)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateSubnetworkRequest{
			ParentID: net1.IdentifierID,
			Name:     "test-subnetwork-11",
		})
		subnet11 := &resourceNetwork.Subnetwork{}
		err = suc.Create(ctx, subnet11)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateSubnetworkRequest{
			ParentID: net1.IdentifierID,
			Name:     "test-subnetwork-12",
		})
		subnet12 := &resourceNetwork.Subnetwork{}
		err = suc.Create(ctx, subnet12)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateFirewallRequest{
			ParentID: net2.IdentifierID,
			Name:     "test-firewall-21",
		})
		fw21 := &resourceNetwork.Firewall{}
		err = fuc.Create(ctx, fw21)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateVMRequest{
			ParentID: subnet11.IdentifierID,
			Name:     "test-vm-111",
		})
		vm111 := &resourceInstance.VM{}
		err = vuc.Create(ctx, vm111)
		require.True(t, err.IsOk())

		foundProvider := &resource.Provider{}
		ctx.Set(context.RequestKey, dto.GetProviderStackRequest{
			ProviderID: provider.IdentifierID.Provider,
		})
		err = uc.GetStack(ctx, foundProvider)
		assert.True(t, err.IsOk())
		assert.Equal(t, provider.IdentifierName, foundProvider.IdentifierName)
		assert.Equal(t, 2, len(foundProvider.Networks))
		assert.Equal(t, 2, len(foundProvider.Networks["test-network-1"].Subnetworks))
		assert.Equal(t, 1, len(foundProvider.Networks["test-network-2"].Firewalls))
		assert.Equal(t, "test-vm-111", foundProvider.Networks["test-network-1"].Subnetworks["test-subnetwork-11"].VMs["test-vm-111"].IdentifierName.VM)

	})
}

func Test_providerUseCase_List(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearProviderFixtures(ctx)
		clear(ctx)
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
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	uc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearProviderFixtures(ctx)
		clear(ctx)
	}()

	t.Run("Update a valid provider should succeed", func(t *testing.T) {
		provider := ProviderFixture(ctx, t, uc)

		// New secret
		secrRepo := secretRepo.NewSecretRepository()
		kubeSecretRepo := secretRepo.NewKubernetesRepository()
		createdSecret, err := kubeSecretRepo.Create(ctx, option.Option{
			Value: secret2.CreateKubernetesSecretRequest{
				ProjectID:   ctx.Value(context.ProjectIDKey).(string),
				SecretName:  "test-2",
				SecretKey:   "key.json",
				SecretValue: []byte("{\"test\":\"test\"}"),
			},
		})
		require.True(t, err.IsOk())
		secretAuth := secretModel.NewSecret("test-2", "A new secret", *createdSecret, domainTypes.ProviderGCP)
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
