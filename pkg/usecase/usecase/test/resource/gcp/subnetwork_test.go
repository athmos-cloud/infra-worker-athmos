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
	testResource "github.com/athmos-cloud/infra-worker-athmos/pkg/usecase/usecase/test/resource"
	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/orlangure/gnomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upbound/provider-gcp/apis/compute/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

type wantSubnetwork struct {
	Name   string
	Labels map[string]string
	Spec   v1beta1.SubnetworkSpec
}

func initSubnetwork(t *testing.T) (context.Context, *testResource.TestResource, usecase.Subnetwork) {
	ctx, testNet, nuc := initNetwork(t)
	parentID := ctx.Value(testResource.ProviderIDKey).(identifier.Provider)
	req := dto.CreateNetworkRequest{
		ParentIDProvider: &parentID,
		Name:             "test-net",
		Managed:          false,
	}
	ctx.Set(context.RequestKey, req)
	net := &resource.Network{}
	err := nuc.Create(ctx, net)
	require.True(t, err.IsOk())
	ctx.Set(testResource.NetworkIDKey, net.IdentifierID)
	uc := usecase.NewSubnetworkUseCase(testNet.ProjectRepo, gcp.NewRepository(), nil, nil)
	return ctx, testNet, uc
}

func Test_subnetworkUseCase_Create(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, suc := initSubnetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()

	t.Run("Create a valid subnetwork", func(t *testing.T) {
		parentID := ctx.Value(testResource.NetworkIDKey).(identifier.Network)
		netName := fmt.Sprintf("%s-%s", "test", utils.RandomString(5))
		region := "europe-west1"
		ipCIDR := "10.0.0.1/26"
		req := dto.CreateSubnetworkRequest{
			ParentID:    parentID,
			Name:        netName,
			Region:      region,
			IPCIDRRange: ipCIDR,
			Managed:     false,
		}
		ctx.Set(context.RequestKey, req)
		subnet := &resource.Subnetwork{}
		err := suc.Create(ctx, subnet)
		require.Equal(t, errors.Created.Code, err.Code)

		kubeResource := &v1beta1.Subnetwork{}
		errk := kubernetes.Client().Client.Get(ctx, types.NamespacedName{Name: subnet.IdentifierID.Network}, kubeResource)
		assert.NoError(t, errk)
		wantLabels := map[string]string{
			"app.kubernetes.io/managed-by": "athmos",
			"athmos.cloud/project-id":      ctx.Value(context.ProjectIDKey).(string),
			"identifier.provider":          subnet.IdentifierID.Provider,
			"identifier.vpc":               subnet.IdentifierID.VPC,
			"identifier.network":           subnet.IdentifierID.Network,
			"identifier.subnetwork":        subnet.IdentifierID.Subnetwork,
			"name.provider":                "test",
			"name.vpc":                     "test",
			"name.network":                 "test-net",
			"name.subnetwork":              netName,
		}
		wantSpec := v1beta1.SubnetworkSpec{
			ResourceSpec: v1.ResourceSpec{
				DeletionPolicy: "Orphan",
				ProviderConfigReference: &v1.Reference{
					Name: subnet.IdentifierID.Provider,
				},
				ManagementPolicy: "FullControl",
			},
			ForProvider: v1beta1.SubnetworkParameters_2{
				Network:     &subnet.IdentifierID.Network,
				Project:     &subnet.IdentifierID.VPC,
				Region:      &region,
				IPCidrRange: &ipCIDR,
			},
		}
		wantNet := wantSubnetwork{
			Name:   subnet.IdentifierID.Network,
			Labels: wantLabels,
			Spec:   wantSpec,
		}
		gotNet := wantSubnetwork{
			Name:   kubeResource.Name,
			Labels: kubeResource.Labels,
			Spec:   kubeResource.Spec,
		}
		assertSubnetworkEqual(t, wantNet, gotNet)
	})
	t.Run("Create a subnetwork with a non-existing secret should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Create a subnetwork with an already existing name should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_subnetworkUseCase_Delete(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, _ := initSubnetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Delete a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a subnetwork with children should fail", func(t *testing.T) {

	})
	t.Run("Delete cascade a subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func Test_subnetworkUseCase_Get(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, _ := initSubnetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Get a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Delete a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})

}

func Test_subnetworkUseCase_Update(t *testing.T) {
	mongoC := test.Init(t)
	ctx, _, _ := initSubnetwork(t)
	defer func() {
		require.NoError(t, gnomock.Stop(mongoC))
		clear(ctx)
	}()
	t.Run("Update a valid subnetwork should succeed", func(t *testing.T) {
		t.Skip("TODO")
	})
	t.Run("Update a non-existing subnetwork should fail", func(t *testing.T) {
		t.Skip("TODO")
	})
}

func assertSubnetworkEqual(t *testing.T, want wantSubnetwork, got wantSubnetwork) {
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.Labels, got.Labels)
	assert.Equal(t, want.Spec, got.Spec)
}
