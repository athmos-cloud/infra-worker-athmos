package aws

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/aws"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-aws/apis/networkfirewall/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

type wantFirewall struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.FirewallSpec
}

func tearDown(ctx context.Context) {
	ClearFirewallFixtures(ctx)
}

func Test_firewallUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearFixtures(ctx)
	}()

	t.Run("Create a valid firewall", func(t *testing.T) {
		defer tearDown(ctx)

		firewall := FirewallFixture(ctx, t, fuc)

		kubeResource := &v1beta1.Firewall{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall}, kubeResource)

		assert.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          firewall.IdentifierID.Provider,
			"identifier.vpc":               firewall.IdentifierID.VPC,
			"identifier.network":           firewall.IdentifierID.Network,
			"identifier.firewall":          firewall.IdentifierID.Firewall,
			"name.provider":                "fixture-provider",
			"name.vpc":                     "test",
			"name.network":                 "fixture-network",
			"name.firewall":                firewall.IdentifierName.Firewall,
		}

		wantSpec := v1beta1.FirewallSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: firewall.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.FirewallParameters{
				VPCID: &firewall.IdentifierID.Network,
			},
		}

		want := wantFirewall{
			Name:   firewall.IdentifierID.Firewall,
			Labels: wantLabels,
			Spec:   wantSpec,
		}

		got := wantFirewall{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertFirewallEqual(t, want, got)
	})

	t.Run("Create a firewall with an already existing name should fail", func(t *testing.T) {
		defer tearDown(ctx)

		firewall := FirewallFixture(ctx, t, fuc)
		err := fuc.Create(ctx, firewall)

		assert.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_firewallUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearFixtures(ctx)
	}()

	t.Run("Delete a valid firewall should succeed", func(t *testing.T) {
		defer tearDown(ctx)

		firewall := FirewallFixture(ctx, t, fuc)
		ctx.Set(context.RequestKey, dto.DeleteFirewallRequest{IdentifierID: firewall.IdentifierID})
		delFirewall := &network.Firewall{}
		err := fuc.Delete(ctx, delFirewall)
		require.Equal(t, errors.NoContent.Code, err.Code)
	})

	t.Run("Delete a non-existing firewall should fail", func(t *testing.T) {
		defer tearDown(ctx)

		id := identifier.Firewall{
			Provider: "test",
			Network:  "test",
			Firewall: "this-firewall-does-not-exist",
		}
		ctx.Set(context.RequestKey, dto.DeleteFirewallRequest{IdentifierID: id})
		delFirewall := &network.Firewall{}
		err := fuc.Delete(ctx, delFirewall)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_firewallUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearFixtures(ctx)
	}()

	t.Run("Get a valid firewall should succeed", func(t *testing.T) {
		defer tearDown(ctx)

		firewall := FirewallFixture(ctx, t, fuc)
		getReq := dto.GetFirewallRequest{IdentifierID: firewall.IdentifierID}
		ctx.Set(context.RequestKey, getReq)
		getFirewall := &network.Firewall{}
		err := fuc.Get(ctx, getFirewall)

		assert.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, firewall.IdentifierName, getFirewall.IdentifierName)
		assert.Equal(t, firewall.IdentifierID, getFirewall.IdentifierID)
		assert.Equal(t, firewall.Allow, getFirewall.Allow)
		assert.Equal(t, firewall.Deny, getFirewall.Deny)
	})

	t.Run("Get a non-existing firewall should fail", func(t *testing.T) {
		defer tearDown(ctx)

		getReq := dto.GetFirewallRequest{IdentifierID: identifier.Firewall{
			Provider: "test",
			Network:  "test",
			Firewall: "this-firewall-does-not-exist",
		}}
		ctx.Set(context.RequestKey, getReq)
		getFirewall := &network.Firewall{}
		err := fuc.Get(ctx, getFirewall)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func Test_firewallUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, resourceTest := initTest(t)
	awsRepo := aws.NewRepository()
	puc := usecase.NewProviderUseCase(
		resourceTest.ProjectRepo,
		resourceTest.SecretRepo,
		nil,
		awsRepo,
		nil)
	nuc := usecase.NewNetworkUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)
	fuc := usecase.NewFirewallUseCase(resourceTest.ProjectRepo, nil, awsRepo, nil)

	ProviderFixture(ctx, t, puc)
	NetworkFixture(ctx, t, nuc)

	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		ClearFixtures(ctx)
	}()

	t.Run("Update a valid firewall should succeed", func(t *testing.T) {
		defer tearDown(ctx)

		firewall := FirewallFixture(ctx, t, fuc)
		port := "42"
		protocol := "tcp"
		updateReq := dto.UpdateFirewallRequest{
			IdentifierID: firewall.IdentifierID,
			AllowRules: &network.FirewallRuleList{
				{
					Protocol: protocol,
					Ports:    []string{port},
				},
			},
		}
		ctx.Set(context.RequestKey, updateReq)
		updateFirewall := &network.Firewall{}
		err := fuc.Update(ctx, updateFirewall)
		if !err.IsOk() {
			fmt.Print(err.Message)
		}
		assert.Equal(t, errors.NoContent.Code, err.Code)

		kubeResource := &v1beta1.Firewall{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall}, kubeResource)
		assert.NoError(t, errk)
	})

	t.Run("Update a non-existing firewall should fail", func(t *testing.T) {
		defer tearDown(ctx)

		updateReq := dto.UpdateFirewallRequest{
			IdentifierID: identifier.Firewall{
				Provider: "test",
				Network:  "test",
				Firewall: "this-firewall-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, updateReq)
		updateFirewall := &network.Firewall{}
		err := fuc.Update(ctx, updateFirewall)
		assert.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func assertFirewallEqual(t *testing.T, want wantFirewall, got wantFirewall) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec.ResourceSpec, got.Spec.ResourceSpec)
	assert.Equal(t, *want.Spec.ForProvider.VPCID, *got.Spec.ForProvider.VPCID)
}
