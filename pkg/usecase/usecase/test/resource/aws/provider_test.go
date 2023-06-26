package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	repository2 "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	secretRepo "github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/secret"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/instance"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
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

	defer suiteTeardown(ctx, t, mongoC)

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)
	t.Run("Create a valid provider", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

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

	defer suiteTeardown(ctx, t, mongoC)

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)

	t.Run("Delete a valid provider should succeed", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

		sshRepo := repository2.NewSSHKeyRepository()
		nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
		suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
		vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, sshRepo, nil, awsRepo, nil)

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

	defer suiteTeardown(ctx, t, mongoC)

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)
	t.Run("Get a valid provider should succeed", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

		provider := ProviderFixture(ctx, t, uc)

		getReq := dto.GetResourceRequest{
			Identifier: identifier.Provider{
				Provider: provider.IdentifierID.Provider,
			}.Provider,
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
		defer ClearProviderFixtures(ctx)

		getReq := dto.GetResourceRequest{
			Identifier: identifier.Provider{
				Provider: "this-provider-does-not-exist",
			}.Provider,
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
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
	suc := usecase.NewSubnetworkUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
	fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, aws.NewRepository(), nil)
	vuc := usecase.NewVMUseCase(resourceTest.ProjectRepo, repository2.NewSSHKeyRepository(), nil, aws.NewRepository(), nil)

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Get a valid Provider Stack should succeed", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

		region := "eu-west-1"
		provider := ProviderFixture(ctx, t, uc)
		ctx.Set(context.RequestKey, dto.CreateNetworkRequest{
			ParentIDProvider: &provider.IdentifierID,
			Name:             "test-network-1",
			Region:           region,
		})
		net1 := &network.Network{}
		err := nuc.Create(ctx, net1)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateNetworkRequest{
			ParentIDProvider: &provider.IdentifierID,
			Name:             "test-network-2",
			Region:           region,
		})
		net2 := &network.Network{}
		err = nuc.Create(ctx, net2)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateSubnetworkRequest{
			ParentID: net1.IdentifierID,
			Name:     "test-subnetwork-11",
			Region:   region,
		})

		subnet11 := &network.Subnetwork{}
		err = suc.Create(ctx, subnet11)
		require.True(t, err.IsOk())
		ctx.Set(context.RequestKey, dto.CreateSubnetworkRequest{
			ParentID: net1.IdentifierID,
			Name:     "test-subnetwork-12",
			Region:   region,
		})
		subnet12 := &network.Subnetwork{}
		err = suc.Create(ctx, subnet12)
		require.True(t, err.IsOk())

		ctx.Set(context.RequestKey, dto.CreateFirewallRequest{
			ParentID: net2.IdentifierID,
			Name:     "test-firewall-21",
			AllowRules: network.FirewallRuleList{
				{
					Protocol: "tcp",
					Ports:    []string{"80", "443"},
				},
				{
					Protocol: "udp",
					Ports:    []string{"53"},
				},
			},
			DenyRules: network.FirewallRuleList{
				{
					Protocol: "tcp",
					Ports:    []string{"65"},
				},
			},
			Managed: false,
		})
		fw21 := &network.Firewall{}
		err = fuc.Create(ctx, fw21)
		require.True(t, err.IsOk())

		ctx.Set(context.RequestKey, dto.CreateVMRequest{
			ParentID:    subnet11.IdentifierID,
			Name:        "test-vm-111",
			Zone:        region,
			MachineType: "t2.micro",
			Auths: []dto.VMAuth{
				{
					Username: "admin",
				}, {
					Username:     "test",
					RSAKeyLength: 1024,
				},
			},
			OS: instance.VMOS{
				ID:   "ami-0a5d9cd4e632d99c1",
				Name: "ami-0a5d9cd4e632d99c1",
			},
			Disks: []instance.VMDisk{
				{
					AutoDelete: true,
					Mode:       instance.DiskModeReadWrite,
					Type:       instance.DiskTypeSSD,
					SizeGib:    10,
				},
			},
			Managed: true,
		})
		vm111 := &instance.VM{}
		err = vuc.Create(ctx, vm111)
		fmt.Println(err.Message)
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

	defer suiteTeardown(ctx, t, mongoC)

	ctx.Set(context.ResourceTypeKey, domainTypes.ProviderResource)

	t.Run("List providers should succeed", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

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

	defer suiteTeardown(ctx, t, mongoC)

	t.Run("Update a valid provider should succeed", func(t *testing.T) {
		defer ClearProviderFixtures(ctx)

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
		ctx.Set(context.RequestKey, dto.GetResourceRequest{
			Identifier: provider.IdentifierID.Provider,
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
		defer ClearProviderFixtures(ctx)

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
		defer ClearProviderFixtures(ctx)

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
