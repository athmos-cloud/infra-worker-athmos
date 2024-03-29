package gcp

import (
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/network"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/logger"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

type wantFirewall struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.FirewallSpec
}

func initFirewall(t *testing.T) (context.Context, *testResource.TestResource, usecase.Firewall) {
	ctx, testNet, nuc := initNetwork(t)
	parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
	req := dto.CreateNetworkRequest{
		ParentIDProvider: &parentID,
		Name:             "test-net",
		Managed:          false,
	}
	ctx.Set(context.RequestKey, req)
	net := &network.Network{}
	err := nuc.Create(ctx, net)
	require.True(t, err.IsOk())
	ctx.Set(testResource.NetworkIDKey, net.IdentifierID)
	uc := usecase.NewFirewallUseCase(testNet.ProjectRepo, gcp.NewRepository(), nil, nil)
	return ctx, testNet, uc
}

func clearFirewall(ctx context.Context) {
	clearSubnetwork(ctx)
	firewalls := &v1beta1.FirewallList{}

	err := kubernetes.Client().Client.List(ctx, firewalls)
	if err != nil {
		return
	}
	for _, firewall := range firewalls.Items {
		err = kubernetes.Client().Client.Delete(ctx, &firewall)
		if err != nil {
			logger.Warning.Printf("Error deleting firewall %s: %v", firewall.Name, err)
			continue
		}
	}
}

func Test_firewallUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, fuc := initFirewall(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearFirewall(ctx)
	}()
	t.Run("_createSqlPasswordSecret a valid firewall", func(t *testing.T) {
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
			"name.provider":                "test",
			"name.vpc":                     "test",
			"name.network":                 "test-net",
			"name.firewall":                firewall.IdentifierName.Firewall,
		}
		var allow []v1beta1.AllowParameters
		for _, rule := range firewall.Allow {
			for _, port := range rule.Ports {
				p := port
				allow = append(allow, v1beta1.AllowParameters{
					Protocol: &rule.Protocol,
					Ports:    []*string{&p},
				})
			}
		}
		var deny []v1beta1.DenyParameters
		for _, rule := range firewall.Deny {
			for _, port := range rule.Ports {
				p := port
				deny = append(deny, v1beta1.DenyParameters{
					Protocol: &rule.Protocol,
					Ports:    []*string{&p},
				})
			}

		}
		wantSpec := v1beta1.FirewallSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Delete",
				ProviderConfigReference: &v1.Reference{
					Name: firewall.IdentifierID.Provider,
				},
			},
			ForProvider: v1beta1.FirewallParameters{
				Network: &firewall.IdentifierID.Network,
				Project: &firewall.IdentifierID.VPC,
				Allow:   allow,
				Deny:    deny,
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

	t.Run("_createSqlPasswordSecret a firewall with an already existing name should fail", func(t *testing.T) {
		firewall := FirewallFixture(ctx, t, fuc)
		err := fuc.Create(ctx, firewall)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_firewallUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, fuc := initFirewall(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearFirewall(ctx)
	}()
	t.Run("Delete a valid firewall should succeed", func(t *testing.T) {
		firewall := FirewallFixture(ctx, t, fuc)
		ctx.Set(context.RequestKey, dto.DeleteFirewallRequest{IdentifierID: firewall.IdentifierID})
		delFirewall := &network.Firewall{}
		err := fuc.Delete(ctx, delFirewall)
		require.Equal(t, errors.NoContent.Code, err.Code)
	})
	t.Run("Delete a non-existing firewall should fail", func(t *testing.T) {
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
	ctx, _, fuc := initFirewall(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearFirewall(ctx)
	}()
	t.Run("Get a valid firewall should succeed", func(t *testing.T) {
		firewall := FirewallFixture(ctx, t, fuc)

		getReq := dto.GetResourceRequest{Identifier: firewall.IdentifierID.Firewall}
		ctx.Set(context.RequestKey, getReq)
		getFirewall := &network.Firewall{}
		err := fuc.Get(ctx, getFirewall)
		assert.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, firewall.IdentifierName, getFirewall.IdentifierName)
		assert.Equal(t, firewall.IdentifierID, getFirewall.IdentifierID)
		assert.Equal(t, firewall.Allow, getFirewall.Allow)
		assert.Equal(t, firewall.Deny, getFirewall.Deny)
	})
	t.Run("Delete a non-existing firewall should fail", func(t *testing.T) {
		getReq := dto.GetResourceRequest{Identifier: identifier.Firewall{
			Provider: "test",
			Network:  "test",
			Firewall: "this-firewall-does-not-exist",
		}.Firewall}
		ctx.Set(context.RequestKey, getReq)
		getFirewall := &network.Firewall{}
		err := fuc.Get(ctx, getFirewall)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})

}

func Test_firewallUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, fuc := initFirewall(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clearFirewall(ctx)
	}()
	t.Run("Update a valid firewall should succeed", func(t *testing.T) {
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
		assert.Equal(t, errors.NoContent.Code, err.Code)
		kubeResource := &v1beta1.Firewall{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: firewall.IdentifierID.Firewall}, kubeResource)
		assert.NoError(t, errk)

		expect := []v1beta1.AllowParameters{
			{
				Protocol: &protocol,
				Ports:    []*string{&port},
			},
		}
		assert.Equal(t, expect, kubeResource.Spec.ForProvider.Allow)

	})
	t.Run("Update a non-existing firewall should fail", func(t *testing.T) {
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
	assert.Equal(t, want.Spec, got.Spec)
}
