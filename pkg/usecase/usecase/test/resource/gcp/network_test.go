package gcp

import (
	"fmt"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/controller/context"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/dto"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/adapter/repository/gcp"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/domain/model/resource/identifier"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/infrastructure/kubernetes"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/errors"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/kernel/utils"
	usecase "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/resource"
	"github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"

	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/types"
	"testing"
)

type wantNetwork struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.NetworkSpec
}

func initNetwork(t *testing.T) (context.Context, *testResource.TestResource, usecase.Network) {
	ctx, testNet, puc := initTest(t)

	req := dto.CreateProviderRequest{
		Name:           "test",
		VPC:            testResource.SecretTestName,
		SecretAuthName: testResource.SecretTestName,
	}
	ctx.Set(context.RequestKey, req)
	provider := &resource.Provider{}
	err := puc.Create(ctx, provider)
	require.True(t, err.IsOk())
	ctx.Set(testResource.ProviderIDKey, provider.IdentifierID)

	nuc := usecase.NewNetworkUseCase(testNet.ProjectRepo, gcp.NewRepository(), nil, nil)
	return ctx, testNet, nuc
}

func Test_networkUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, nuc := initNetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()

	t.Run("Create a valid network", func(t *testing.T) {
		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
			Managed:          false,
		}
		ctx.Set(context.RequestKey, req)
		net := &resource.Network{}
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Created.Code, err.Code)

		kubeResource := &v1beta1.Network{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: net.IdentifierID.Network}, kubeResource)
		assert.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          net.IdentifierID.Provider,
			"identifier.vpc":               net.IdentifierID.VPC,
			"identifier.network":           net.IdentifierID.Network,
			"name.provider":                "test",
			"name.vpc":                     "test",
			"name.network":                 netName,
		}
		autoCreateSubnet := false
		wantSpec := v1beta1.NetworkSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Orphan",
				ProviderConfigReference: &v1.Reference{
					Name: net.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.NetworkParameters{
				Project:               &net.IdentifierID.VPC,
				AutoCreateSubnetworks: &autoCreateSubnet,
			},
		}
		wantNet := wantNetwork{
			Name:   net.IdentifierID.Network,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantNetwork{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertNetworkEqual(t, wantNet, gotNet)
	})
	t.Run("Create a network with an already existing name should raise conflict error", func(t *testing.T) {
		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
		}
		ctx.Set(context.RequestKey, req)
		net := &resource.Network{}
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Created.Code, err.Code)
		err = nuc.Create(ctx, net)
		require.Equal(t, errors.Conflict.Code, err.Code)
	})
}

func Test_networkUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, nuc := initNetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()

	t.Run("Delete a valid network should succeed", func(t *testing.T) {
		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
		}
		ctx.Set(context.RequestKey, req)
		net := &resource.Network{}
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Created.Code, err.Code)
		delReq := dto.DeleteNetworkRequest{
			IdentifierID: net.IdentifierID,
		}
		ctx.Set(context.RequestKey, delReq)
		err = nuc.Delete(ctx, net)
		require.Equal(t, errors.NoContent.Code, err.Code)
	})
	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		delReq := dto.DeleteNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, delReq)
		net := &resource.Network{}
		err := nuc.Delete(ctx, net)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
	t.Run("Delete a network with children should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete cascade a network should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_networkUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, nuc := initNetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()

	t.Run("Get a valid network should succeed", func(t *testing.T) {
		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
		}
		ctx.Set(context.RequestKey, req)
		net := &resource.Network{}
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Created.Code, err.Code)
		getReq := dto.GetNetworkRequest{
			IdentifierID: net.IdentifierID,
		}
		ctx.Set(context.RequestKey, getReq)
		gotNet := &resource.Network{}
		err = nuc.Get(ctx, gotNet)
		require.Equal(t, errors.OK.Code, err.Code)
		assert.Equal(t, net, gotNet)

	})
	t.Run("Delete a non-existing network should fail", func(t *testing.T) {
		getReq := dto.GetNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, getReq)
		gotNet := &resource.Network{}
		err := nuc.Get(ctx, gotNet)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})

}

func Test_networkUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, nuc := initNetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Update a valid network should succeed", func(t *testing.T) {
		parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		req := dto.CreateNetworkRequest{
			ParentIDProvider: &parentID,
			Name:             netName,
		}
		ctx.Set(context.RequestKey, req)
		net := &resource.Network{}
		err := nuc.Create(ctx, net)
		require.Equal(t, errors.Created.Code, err.Code)
		managed := false
		toUp := dto.UpdateNetworkRequest{
			IdentifierID: net.IdentifierID,
			Managed:      &managed,
		}
		ctx.Set(context.RequestKey, toUp)
		err = nuc.Update(ctx, net)
		require.Equal(t, errors.NoContent.Code, err.Code)
		kubeResource := &v1beta1.Network{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: net.IdentifierID.Network}, kubeResource)
		assert.NoError(t, errk)
		assert.Equal(t, kubeResource.Spec.DeletionPolicy, v1.DeletionOrphan)
	})
	t.Run("Update a non-existing network should fail", func(t *testing.T) {
		toUp := dto.UpdateNetworkRequest{
			IdentifierID: identifier.Network{
				Provider: "test",
				VPC:      "test",
				Network:  "this-network-does-not-exist",
			},
		}
		ctx.Set(context.RequestKey, toUp)
		net := &resource.Network{}
		err := nuc.Update(ctx, net)
		require.Equal(t, errors.NotFound.Code, err.Code)
	})
}

func assertNetworkEqual(t *testing.T, want wantNetwork, got wantNetwork) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec, got.Spec)
}
